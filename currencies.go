package main

import (
	"github.com/shopspring/decimal"
	bolt "go.etcd.io/bbolt"
)

type School struct {
	db       *bolt.DB
	schoolId string
}

func (s School) xRateTx(tx *bolt.Tx, schoolId, from, to string) (float64, error) {
	if from == to {
		return 1, nil
	}

	return 1.0, nil

}

func addStepTx(tx *bolt.Tx, schoolId string, clock Clock, currencyId string, amount float32) error {
	cb, err := getCbTx(tx, schoolId)
	if err != nil {
		return err
	}

	accounts, err := cb.CreateBucketIfNotExists([]byte(KeyAccounts))
	if err != nil {
		return err
	}

	account, err := accounts.CreateBucketIfNotExists([]byte(currencyId))
	if err != nil {
		return err
	}
	value, err := account.CreateBucketIfNotExists([]byte(KeyValue))
	if err != nil {
		return err
	}

	cursor := value.Cursor()
	k, v := cursor.Last()
	if k == nil {
		now := clock.Now()
		d := decimal.NewFromFloat32(amount)
		step, err := d.MarshalText()
		if err != nil {
			return err
		}
		value.Put([]byte(now.String()), step)
	} else {
		var last decimal.Decimal
		err = last.UnmarshalText(v)
		if err != nil {
			return err
		}

		d := decimal.NewFromFloat32(amount)

		d = (last.Mul(decimal.NewFromInt(99)).Add(d)).Div(decimal.NewFromInt(100))

		step, err := d.MarshalText()
		if err != nil {
			return err
		}
		err = value.Put([]byte(k), step)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}

func getCurrencyValue(schoolId, currencyId string) (float64, error) {
	return 1.0, nil
}

func convert(schoolId, from, to string, amount float64) (float64, error) {
	return amount, nil

}
