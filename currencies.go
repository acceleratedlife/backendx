package main

import (
	"fmt"
	"github.com/go-pkgz/lgr"
	"github.com/shopspring/decimal"
	bolt "go.etcd.io/bbolt"
)

// rate of currency 'from' comparative to 'base'
// how much 'base' to buy 1 'from'
// from, base - empty value refers to uBuck
func xRateFromToBaseTx(tx *bolt.Tx, schoolId, from, base string) (rate decimal.Decimal, err error) {

	if from == base {
		return decimal.NewFromInt32(1), nil
	}
	fromValue := decimal.Zero
	baseValue := decimal.Zero

	if from == "" {
		fromValue, err = getUbuckValueTx(tx, schoolId)
		if err != nil {
			return decimal.Zero, err
		}
	} else {
		fromValue, err = getCurrencyMMATx(tx, schoolId, from)
		if err != nil {
			return decimal.Zero, err
		}
	}

	if base == "" {
		baseValue, err = getUbuckValueTx(tx, schoolId)
		if err != nil {
			return decimal.Zero, err
		}
	} else {
		baseValue, err = getCurrencyMMATx(tx, schoolId, from)
		if err != nil {
			return decimal.Zero, err
		}
	}

	return baseValue.DivRound(fromValue, 6), nil

}

// updates MMA
// returns MMA
func addStepTx(tx *bolt.Tx, schoolId string, currencyId string, amount float32) (decimal.Decimal, error) {
	cb, err := getCbTx(tx, schoolId)
	if err != nil {
		return decimal.Zero, err
	}

	accounts, err := cb.CreateBucketIfNotExists([]byte(KeyAccounts))
	if err != nil {
		return decimal.Zero, err
	}

	account, err := accounts.CreateBucketIfNotExists([]byte(currencyId))
	if err != nil {
		return decimal.Zero, err
	}

	v := account.Get([]byte(KeyMMA))
	var d decimal.Decimal
	if v != nil {
		var last decimal.Decimal
		err = last.UnmarshalText(v)
		if err != nil {
			return decimal.Zero, err
		}
		d = (last.Mul(decimal.NewFromFloat(99)).Add(d)).Div(decimal.NewFromFloat(100))
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

// calculates avg of all currencies
func getUbuckValueTx(tx *bolt.Tx, schoolId string) (decimal.Decimal, error) {
	cb, err := getCbTx(tx, schoolId)
	if err != nil {
		return decimal.Zero, err
	}

	accounts, err := cb.CreateBucketIfNotExists([]byte(KeyAccounts))
	if err != nil {
		return decimal.Zero, err
	}

	return getUbuckValue1Tx(accounts)
}
func getUbuckValue1Tx(accounts *bolt.Bucket) (decimal.Decimal, error) {
	sum := decimal.Zero
	validCurrencies := int32(0)
	_ = accounts.ForEach(func(k, v []byte) error {
		if v != nil {
			return nil
		}
		account := accounts.Bucket(k)

		d, err := getCurrencyAccMMATx(account)
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

func updateXRates(accounts bolt.Bucket, uBuckValue decimal.Decimal, clock Clock) error {
	return accounts.ForEach(func(k, v []byte) error {
		if v != nil {
			return nil
		}
		account := accounts.Bucket(k)

		d, err := getCurrencyAccMMATx(account)
		if err != nil {
			lgr.Printf("ERROR exclude currency %s from calculating avg - %v ", string(k), err)
			return nil
		}

		xRate := uBuckValue.Div(d)

		value, err := account.CreateBucketIfNotExists([]byte(KeyValue))
		if err != nil {
			return fmt.Errorf("cannot create bucket: %v ", err)
		}
		now := clock.Now()

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

// returns current MMA
func getCurrencyMMATx(tx *bolt.Tx, schoolId, currencyId string) (decimal.Decimal, error) {
	cb, err := getCbTx(tx, schoolId)
	if err != nil {
		return decimal.Zero, err
	}

	accounts := cb.Bucket([]byte(KeyAccounts))
	if accounts == nil {
		return decimal.Zero, fmt.Errorf("no accounts bucket")
	}
	account := accounts.Bucket([]byte(currencyId))
	return getCurrencyAccMMATx(account)
}
func getCurrencyAccMMATx(account *bolt.Bucket) (decimal.Decimal, error) {
	if account == nil {
		return decimal.Zero, fmt.Errorf("account does not exist")
	}
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
