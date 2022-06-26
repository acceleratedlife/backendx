package main

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"
)

func TestSchool_xRateTx(t *testing.T) {
	// if only currency defined - it is equivalent to uBuck
	db, teardown := OpenTestDB("currency")
	defer teardown()

	_, schools, teachers, _, _, err := CreateTestAccounts(db, 2, 10, 1, 2)
	require.Nil(t, err)

	_ = db.Update(func(tx *bolt.Tx) error {
		_, err2 := xRateToBaseInstantRx(tx, schools[0], teachers[0], "")

		require.NotNil(t, err2)

		// add first currency in a school
		_, _ = addStepTx(tx, schools[0], teachers[0], 10)
		r, err2 := xRateToBaseInstantRx(tx, schools[0], teachers[0], "")
		require.Equal(t, 1.0, r.InexactFloat64())

		// payment by 2nd teacher
		_, _ = addStepTx(tx, schools[0], teachers[1], 20)

		r, err2 = xRateToBaseInstantRx(tx, schools[0], teachers[0], "")
		require.Equal(t, 1.5, r.InexactFloat64())
		r, err2 = xRateToBaseInstantRx(tx, schools[0], teachers[1], "")
		require.Equal(t, 0.75, r.InexactFloat64())

		r, err2 = xRateToBaseInstantRx(tx, schools[0], teachers[0], teachers[1])
		require.Equal(t, 2.0, r.InexactFloat64())
		r, err2 = xRateToBaseInstantRx(tx, schools[0], teachers[1], teachers[0])
		require.Equal(t, 0.5, r.InexactFloat64())

		return nil
	})
}

func Test_addStepTx(t *testing.T) {
	db, teardown := OpenTestDB("currency")
	defer teardown()

	_, schools, teachers, _, _, err := CreateTestAccounts(db, 2, 10, 1, 2)
	require.Nil(t, err)

	type args struct {
		tx         *bolt.Tx
		schoolId   string
		clock      Clock
		currencyId string
		amount     float32
	}

	_ = db.Update(func(tx *bolt.Tx) error {
		tests := []struct {
			name    string
			args    args
			want    decimal.Decimal
			wantErr assert.ErrorAssertionFunc
		}{
			{
				name: "",
				args: args{
					tx:         tx,
					schoolId:   schools[0],
					currencyId: teachers[0],
					amount:     10,
				},
				want: decimal.NewFromFloat(10),
				wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
					return err == nil
				},
			},
			{
				name: "",
				args: args{
					tx:         tx,
					schoolId:   schools[0],
					currencyId: teachers[0],
					amount:     100,
				},
				want: decimal.NewFromFloat(10.90), // why not 10.91?
				wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
					return err == nil
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := addStepTx(tt.args.tx, tt.args.schoolId, tt.args.currencyId, tt.args.amount)
				if !tt.wantErr(t, err, fmt.Sprintf("addStepTx(%v, %v, %v, %v, %v)", tt.args.tx, tt.args.schoolId, tt.args.clock, tt.args.currencyId, tt.args.amount)) {
					return
				}
				assert.True(t, tt.want.Sub(got).LessThan(decimal.NewFromFloat(0.001)),
					"addStepTx(%v, %v, %v - %v %v)",
					tt.args.schoolId, tt.args.currencyId, tt.args.amount, tt.want.InexactFloat64(), got.InexactFloat64())
			})
		}

		mmaTx, err2 := getCurrencyMMARx(tx, schools[0], teachers[0])
		require.Nil(t, err2)
		assert.True(t, mmaTx.Sub(decimal.NewFromFloat32(10.90)).LessThan(decimal.NewFromFloat(0.001)))
		return nil
	})
}

func Test_HistoricalRates(t *testing.T) {
	db, teardown := OpenTestDB("historical change")
	defer teardown()

	clock := TestClock{}

	_, schools, teachers, _, _, err := CreateTestAccounts(db, 2, 10, 1, 2)
	require.Nil(t, err)

	_ = db.Update(func(tx *bolt.Tx) error {
		cb, err := getCbRx(tx, schools[0])
		require.Nil(t, err)
		accounts := cb.Bucket([]byte(KeyAccounts))
		require.NotNil(t, accounts)

		stepTx, err := addStepTx(tx, schools[0], teachers[0], 10)
		require.Nil(t, err)
		require.Equal(t, 10.0, stepTx.InexactFloat64())

		err = updateXRatesTx(accounts, &clock)
		require.Nil(t, err)

		stepTx, err = addStepTx(tx, schools[0], teachers[0], 10)
		require.Nil(t, err)
		require.Equal(t, 10.0, stepTx.InexactFloat64())
		err = updateXRatesTx(accounts, &clock)
		require.Nil(t, err)

		rx, err := xRateToBaseHistoricalRx(tx, schools[0], teachers[0], "")
		require.Nil(t, err)
		require.Equal(t, 1.0, rx.InexactFloat64())

		stepTx, err = addStepTx(tx, schools[0], teachers[1], 100)
		require.Nil(t, err)
		require.Equal(t, 100.0, stepTx.InexactFloat64())
		err = updateXRatesTx(accounts, &clock)
		require.Nil(t, err)

		/// historical
		rx, err = xRateToBaseHistoricalRx(tx, schools[0], teachers[0], "")
		require.Nil(t, err)
		require.Equal(t, 5.5, rx.InexactFloat64())

		rx, err = xRateToBaseHistoricalRx(tx, schools[0], teachers[1], "")
		require.Nil(t, err)
		require.Equal(t, 0.55, rx.InexactFloat64())

		rx, err = xRateToBaseHistoricalRx(tx, schools[0], teachers[0], teachers[1])
		require.Nil(t, err)
		require.Equal(t, 10.0, rx.InexactFloat64())

		// instant
		rx, err = xRateToBaseInstantRx(tx, schools[0], teachers[0], "")
		require.Nil(t, err)
		require.Equal(t, 5.5, rx.InexactFloat64())

		rx, err = xRateToBaseInstantRx(tx, schools[0], teachers[1], "")
		require.Nil(t, err)
		require.Equal(t, 0.55, rx.InexactFloat64())

		rx, err = xRateToBaseInstantRx(tx, schools[0], teachers[0], teachers[1])
		require.Nil(t, err)
		require.Equal(t, 10.0, rx.InexactFloat64())

		rx, err = xRateToBaseRx(tx, schools[0], teachers[1], teachers[0])
		require.Nil(t, err)
		require.Equal(t, 0.1, rx.InexactFloat64())

		return nil
	})
}

func Test_ubuck2ubuck(t *testing.T) {
	rx, err := xRateToBaseRx(nil, "werwer", "", "")
	require.Nil(t, err)
	require.Equal(t, 1.0, rx.InexactFloat64())

}

func Test_convertRx(t *testing.T) {
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("convertRx")
	defer dbTearDown()

	_, _, teachers, _, students, err := CreateTestAccounts(db, 2, 2, 2, 3)

	userInfo, _ := getUserInLocalStore(db, students[0])
	err = addUbuck2Student(db, &clock, userInfo, decimal.NewFromFloat(1.0), "daily payment")
	require.Nil(t, err)

	balance := StudentNetWorth(db, students[0])
	require.Equal(t, 1.0, balance.InexactFloat64())

	err = pay2Student(db, &clock, userInfo, decimal.NewFromFloat(0.5), teachers[0], "reward")
	require.Nil(t, err)

	balance = StudentNetWorth(db, students[0])
	require.Equal(t, 1.5, balance.InexactFloat64())

	err = pay2Student(db, &clock, userInfo, decimal.NewFromFloat(0.5), teachers[1], "some reason")
	require.Nil(t, err)

	balance = StudentNetWorth(db, students[0])
	require.Equal(t, 2.0, balance.InexactFloat64())
}
