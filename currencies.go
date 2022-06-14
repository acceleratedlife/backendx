package main

import (
	"fmt"
	"strings"

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

	if from == "" {
		fromValue, err = getUbuckValueRx(tx, schoolId)
		if err != nil {
			return decimal.Zero, err
		}
	} else {
		fromValue, err = getCurrencyMMARx(tx, schoolId, from)
		if err != nil {
			return decimal.Zero, err
		}
	}

	if base == "" {
		baseValue, err = getUbuckValueRx(tx, schoolId)
		if err != nil {
			return decimal.Zero, err
		}
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
	if from == base {
		return decimal.NewFromInt32(1), nil
	}

	if from == "" || base == "" {
		rate, err = xRateToBaseHistoricalRx(tx, schoolId, from, base)
		if err == nil {
			return
		} else {
			lgr.Printf("WARN historical xrate calculation failed: %v", err)
		}
	}

	return xRateToBaseInstantRx(tx, schoolId, from, base)
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
	step, err := d.MarshalText()
	if err != nil {
		return decimal.Zero, err
	}
	err = account.Put([]byte(KeyMMA), step)
	if err != nil {
		return decimal.Zero, err
	}
	return d, nil

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
	if currency == "" {
		return decimal.NewFromInt(1), nil
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
		return decimal.Zero, fmt.Errorf("account does not exist")
	}
	return getCurrencyAccMMARx(account)
}
func getCurrencyAccMMARx(account *bolt.Bucket) (decimal.Decimal, error) {

	v := account.Get([]byte(KeyMMA))
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

func convert(schoolId, from, to string, amount float64) (float64, error) {
	return amount, nil

}

func isTeacherAccount(id string) bool {
	return strings.IndexAny(id, "@") != -1
}