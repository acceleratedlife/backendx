package main

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/go-pkgz/lgr"
	"github.com/shopspring/decimal"
	bolt "go.etcd.io/bbolt"
)

// calculate rate of currency 'from' comparative to 'base'
// how much 'base' to buy 1 'from'
// from, base - empty value refers to uBuck
// NB: if uBuck involved iterates through all currencies in CB
func xRateToBaseInstantRx(tx *bolt.Tx, schoolId, from, base string) (rate decimal.Decimal, err error) {

	if from == base {
		return decimal.NewFromInt32(1), nil
	}
	fromValue := decimal.Zero
	baseValue := decimal.Zero

	if from == "" || from == CurrencyUBuck {
		fromValue, err = getUbuckValueRx(tx, schoolId)
		if err != nil {
			return decimal.Zero, err
		}
	} else if from == KeyDebt {
		fromValue, err = getUbuckValueRx(tx, schoolId)
		if err != nil {
			return decimal.Zero, err
		}
		fromValue = fromValue.Neg()
	} else {
		fromValue, err = getCurrencyMMARx(tx, schoolId, from)
		if err != nil {
			return decimal.Zero, err
		}
	}

	if base == "" || base == CurrencyUBuck {
		baseValue, err = getUbuckValueRx(tx, schoolId)
		if err != nil {
			return decimal.Zero, err
		}
	} else if base == KeyDebt {
		baseValue, err = getUbuckValueRx(tx, schoolId)
		if err != nil {
			return decimal.Zero, err
		}
		baseValue = baseValue.Neg()
	} else {
		baseValue, err = getCurrencyMMARx(tx, schoolId, base)
		if err != nil {
			return decimal.Zero, err
		}
	}

	return baseValue.DivRound(fromValue, 6), nil
}

// calculate rate of currency 'from' comparative to 'base'
// how much 'base' to buy 1 'from'
// from, base - empty value refers to uBuck
// NB: uses saved xRates for currencies, only takes saved rates
func xRateToBaseHistoricalRx(tx *bolt.Tx, schoolId, from, base string) (rate decimal.Decimal, err error) {
	if from == base {
		return decimal.NewFromInt32(1), nil
	}
	fromValue, err := getSavedXRateRx(tx, schoolId, from)
	if err != nil {
		return fromValue, err
	}

	baseValue, err := getSavedXRateRx(tx, schoolId, base)
	if err != nil {
		return baseValue, err
	}

	return fromValue.Div(baseValue), nil
}

// calculate rate of currency 'from' comparative to 'base'
// how much 'base' to buy 1 'from'
// from, base - empty value refers to uBuck
func xRateToBaseRx(tx *bolt.Tx, schoolId, from, base string) (rate decimal.Decimal, err error) {
	if (from == KeyDebt && base == CurrencyUBuck) || (from == KeyDebt && base == "") || (from == CurrencyUBuck && base == KeyDebt) || (from == "" && base == KeyDebt) {
		return decimal.NewFromInt32(-1), nil
	}

	if from == base || (from == CurrencyUBuck && base == "") || (from == "" && base == CurrencyUBuck) {
		return decimal.NewFromInt32(1), nil
	}

	if from == "" || from == CurrencyUBuck || base == "" || base == CurrencyUBuck || from == KeyDebt || base == KeyDebt {
		rate, err = xRateToBaseHistoricalRx(tx, schoolId, from, base)
		if err == nil {
			return
		}
	}

	rate, err = xRateToBaseInstantRx(tx, schoolId, from, base)

	return
}

// adds step for pay frequency, modifies MMA based on mean pay frequency
func modifyMmaTx(tx *bolt.Tx, schoolId, currencyId string, currentTrans time.Time, mma decimal.Decimal, clock Clock) (decimal.Decimal, error) {
	account, err := getAccountBucketTx(tx, schoolId, currencyId)
	if err != nil {
		return decimal.Zero, err
	}

	lastTrans, err := getLastTeacherPaymentRX(account, clock)
	if err != nil {
		return decimal.Zero, err
	}

	currentDiff := currentTrans.Sub(lastTrans.Ts)

	var t time.Duration
	w := account.Get([]byte(KeyPayFrequency))
	if w != nil {
		var last time.Duration
		err = json.Unmarshal(w, &last)
		if err != nil {
			return decimal.Zero, err
		}
		t = time.Duration(decimal.NewFromInt(last.Nanoseconds()).Mul(decimal.NewFromFloat(99)).Add(decimal.NewFromInt(currentDiff.Nanoseconds())).Div(decimal.NewFromFloat(100)).IntPart())
	} else {
		t = currentDiff
	}

	marshal, err := json.Marshal(t)
	if err != nil {
		return decimal.Zero, err
	}

	//saving average frequency of payouts
	err = account.Put([]byte(KeyPayFrequency), marshal)
	if err != nil {
		return decimal.Zero, err
	}

	zScore, err := getPayZScoreRx(tx, schoolId, t)
	if err != nil {
		return decimal.Zero, err
	}

	//I think this should be Add but when I put it I get the answer in reverse
	modMMA := decimal.NewFromFloat(1).Sub((zScore.Div(decimal.NewFromFloat(10)))).Mul(mma)

	//saving modMMA
	step1, err := modMMA.MarshalText()
	if err != nil {
		return decimal.Zero, err
	}

	//saving average payouts
	err = account.Put([]byte(KeyModMMA), step1)
	if err != nil {
		return decimal.Zero, err
	}

	return modMMA, nil

}

// updates MMA
// returns MMA
func addStepTx(tx *bolt.Tx, schoolId string, currencyId string, amount float32) (decimal.Decimal, error) {
	account, err := getAccountBucketTx(tx, schoolId, currencyId)
	if err != nil {
		return decimal.Zero, err
	}

	var d decimal.Decimal
	v := account.Get([]byte(KeyMMA))
	if v != nil {
		var last decimal.Decimal
		err = last.UnmarshalText(v)
		if err != nil {
			return decimal.Zero, err
		}
		d = (last.Mul(decimal.NewFromFloat(99)).Add(decimal.NewFromFloat32(amount))).Div(decimal.NewFromFloat(100))
	} else {
		d = decimal.NewFromFloat32(amount)
	}
	step1, err := d.MarshalText()
	if err != nil {
		return decimal.Zero, err
	}

	//saving average payouts
	err = account.Put([]byte(KeyMMA), step1)
	if err != nil {
		return decimal.Zero, err
	}

	return d, nil

}

func getPayZScoreRx(tx *bolt.Tx, schoolId string, payFreq time.Duration) (zScore decimal.Decimal, err error) {
	cb, err := getCbRx(tx, schoolId)
	if err != nil {
		return
	}

	accounts := cb.Bucket([]byte(KeyAccounts))
	if accounts == nil {
		return zScore, fmt.Errorf("ERROR cannot get accounts bucket")
	}

	c := accounts.Cursor()
	var counter = int64(0)
	var sum = time.Duration(0)
	durations := make([]time.Duration, 0)
	for k, _ := c.First(); k != nil; k, _ = c.Next() {
		account := accounts.Bucket(k)
		if account == nil {
			return zScore, fmt.Errorf("ERROR cannot get account bucket")
		}

		payFreqData := account.Get([]byte(KeyPayFrequency))
		if payFreqData == nil {
			continue
		}

		payFreq2, err := time.ParseDuration((string(payFreqData) + "ns"))
		if err != nil {
			return zScore, fmt.Errorf("ERROR cannot parse time")
		}

		durations = append(durations, payFreq2)
		counter++
		sum = sum + payFreq2

	}

	if counter == 0 {
		return decimal.Zero, nil
	}

	meanPayFreq := time.Duration(decimal.NewFromInt(sum.Nanoseconds()).Div(decimal.NewFromInt(counter)).IntPart())

	squaredSums := decimal.Zero
	for _, duration := range durations {
		squaredSums = squaredSums.Add(decimal.NewFromInt((duration - meanPayFreq).Nanoseconds()).Pow(decimal.NewFromFloat(2)))
	}

	variance := squaredSums.Div(decimal.NewFromInt(counter))
	if variance.IsZero() {
		return decimal.Zero, nil
	}
	standardDev := math.Sqrt(variance.InexactFloat64())
	standardDev2 := decimal.NewFromFloat(standardDev)

	zScore = decimal.NewFromInt((payFreq - meanPayFreq).Nanoseconds()).Div(standardDev2)

	return
}

func getLastTeacherPaymentRX(account *bolt.Bucket, clock Clock) (lastTrans Transaction, err error) {
	previousTransBucket := account.Bucket([]byte(KeyTransactions))
	if previousTransBucket == nil {
		lastTrans = Transaction{Ts: clock.Now().AddDate(0, 0, -2)}
		return lastTrans, nil
	}

	c := previousTransBucket.Cursor()
	for k, v := c.Last(); k != nil; k, v = c.Prev() {
		if v == nil {
			continue
		}

		err = json.Unmarshal(v, &lastTrans)
		if err != nil {
			return lastTrans, fmt.Errorf("ERROR cannot unmarshal transaction details")
		}

		if lastTrans.Source != "" {
			continue
		}

		return
	}

	return lastTrans, fmt.Errorf("ERROR cannot find teachers last payout")
}

// calculates avg MMA of all currencies
func getUbuckValueRx(tx *bolt.Tx, schoolId string) (decimal.Decimal, error) {
	cb, err := getCbRx(tx, schoolId)
	if err != nil {
		return decimal.Zero, err
	}

	accounts := cb.Bucket([]byte(KeyAccounts))
	if accounts == nil {
		return decimal.Zero, nil
	}

	return getUbuckValue1Rx(accounts)
}
func getUbuckValue1Rx(accounts *bolt.Bucket) (decimal.Decimal, error) {
	sum := decimal.Zero
	validCurrencies := int32(0)
	_ = accounts.ForEach(func(k, v []byte) error {
		if v != nil {
			return nil
		}
		if !isTeacherAccount(string(k)) {
			return nil
		}
		account := accounts.Bucket(k)

		d, err := getCurrencyAccMMARx(account)
		if err != nil {
			lgr.Printf("ERROR exclude currency %s from calculating avg - %v ", string(k), err)
			return nil
		}

		sum = sum.Add(d)
		validCurrencies += 1
		return nil
	})

	if validCurrencies == 0 {
		return decimal.NewFromInt(1), nil
	}
	return sum.Div(decimal.NewFromInt32(validCurrencies)), nil
}

// saves xrates to uBuck for every currency
func updateXRatesTx(accounts *bolt.Bucket, clock Clock) error {
	uBuckValue, err := getUbuckValue1Rx(accounts)
	if err != nil {
		return err
	}
	now := clock.Now()
	return accounts.ForEach(func(k, v []byte) error {
		if v != nil {
			return nil
		}
		account := accounts.Bucket(k)
		if account == nil {
			return fmt.Errorf("account does not exist")
		}
		d, err := getCurrencyAccMMARx(account)
		if err != nil {
			lgr.Printf("ERROR exclude currency %s from calculating avg - %v ", string(k), err)
			return nil
		}

		xRate := uBuckValue.Div(d)

		value, err := account.CreateBucketIfNotExists([]byte(KeyValue))
		if err != nil {
			return fmt.Errorf("cannot create bucket: %v ", err)
		}

		valueX, err := xRate.MarshalText()
		if err != nil {
			return err
		}
		err = value.Put([]byte(now.String()), valueX)
		if err != nil {
			return fmt.Errorf("cannot save new xrate: %v", err)
		}
		return nil
	})
}

func getSavedXRateRx(tx *bolt.Tx, schoolId, currency string) (rate decimal.Decimal, err error) {
	if currency == "" || currency == CurrencyUBuck {
		return decimal.NewFromInt(1), nil
	}

	if currency == KeyDebt {
		return decimal.NewFromInt(-1), nil
	}

	bucketRx := getAccountBucketRx(tx, schoolId, currency)
	if bucketRx == nil {
		return decimal.Zero, fmt.Errorf("no account for %s", currency)
	}
	history := bucketRx.Bucket([]byte(KeyValue))
	if history == nil {
		return decimal.Zero, fmt.Errorf("no saved history for %s", currency)
	}
	cursor := history.Cursor()
	_, v := cursor.Last()
	if v == nil {
		return decimal.Zero, fmt.Errorf("no saved rates for %s", currency)
	}
	err = rate.UnmarshalText(v)
	return

}

func getAccountBucketRx(tx *bolt.Tx, schoolId, currencyId string) *bolt.Bucket {
	cb, err := getCbRx(tx, schoolId)
	if err != nil {
		return nil
	}

	accounts := cb.Bucket([]byte(KeyAccounts))
	if accounts == nil {
		return nil
	}
	return accounts.Bucket([]byte(currencyId))
}

func getAccountBucketTx(tx *bolt.Tx, schoolId, currencyId string) (*bolt.Bucket, error) {
	cb, err := getCbTx(tx, schoolId)
	if err != nil {
		return nil, err
	}

	accounts, err := cb.CreateBucketIfNotExists([]byte(KeyAccounts))
	if accounts == nil {
		return nil, err
	}
	return accounts.CreateBucketIfNotExists([]byte(currencyId))
}

// returns current MMA
func getCurrencyMMARx(tx *bolt.Tx, schoolId, currencyId string) (decimal.Decimal, error) {
	account := getAccountBucketRx(tx, schoolId, currencyId)
	if account == nil {
		return decimal.Zero, fmt.Errorf("account does not exist, you may need to pay someone first")
	}
	return getCurrencyAccMMARx(account)
}

// returns modMMA
func getCurrencyAccMMARx(account *bolt.Bucket) (decimal.Decimal, error) {

	v := account.Get([]byte(KeyModMMA))
	if v == nil {
		return decimal.Zero, fmt.Errorf("MMA not defined")
	}
	var d decimal.Decimal
	err := d.UnmarshalText(v)
	if err != nil {
		return decimal.Zero, err
	}
	return d, nil
}

func convertRx(tx *bolt.Tx, schoolId, from, to string, amount float64) (converted decimal.Decimal, xRate decimal.Decimal, err error) {
	xRate, err = xRateToBaseRx(tx, schoolId, from, to)
	if err != nil {
		return converted, xRate, err
	}

	converted = xRate.Mul(decimal.NewFromFloat(amount))

	return converted, xRate, nil

}

func isTeacherAccount(id string) bool {
	return strings.ContainsAny(id, "@")
}
