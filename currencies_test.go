package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"
)

func TestSchool_xRateTx(t *testing.T) {
	// if only currency defined - it is equivalent to uBuck
	db, teardown := OpenTestDB("currency")
	defer teardown()
	clock := TestClock{}

	_, schools, teachers, _, _, err := CreateTestAccounts(db, 2, 10, 1, 2)
	require.Nil(t, err)

	_ = db.Update(func(tx *bolt.Tx) error {
		_, err2 := xRateToBaseInstantRx(tx, schools[0], teachers[0], "")

		require.NotNil(t, err2)

		// add first currency in a school
		mma, _ := addStepTx(tx, schools[0], teachers[0], 10)
		_, _ = addStepHelperTx(tx, schools[0], teachers[0], clock.Now(), mma, &clock)
		r, _ := xRateToBaseInstantRx(tx, schools[0], teachers[0], "")
		require.Equal(t, 1.0, r.InexactFloat64())

		// payment by 2nd teacher
		mma, _ = addStepTx(tx, schools[0], teachers[1], 20)
		_, _ = addStepHelperTx(tx, schools[0], teachers[1], clock.Now(), mma, &clock)

		r, _ = xRateToBaseInstantRx(tx, schools[0], teachers[0], "")
		require.Equal(t, 1.5, r.InexactFloat64())
		r, _ = xRateToBaseInstantRx(tx, schools[0], teachers[1], "")
		require.Equal(t, 0.75, r.InexactFloat64())

		r, _ = xRateToBaseInstantRx(tx, schools[0], teachers[0], teachers[1])
		require.Equal(t, 2.0, r.InexactFloat64())
		r, _ = xRateToBaseInstantRx(tx, schools[0], teachers[1], teachers[0])
		require.Equal(t, 0.5, r.InexactFloat64())

		return nil
	})
}

func Test_addStepTx(t *testing.T) {
	db, teardown := OpenTestDB("currency")
	defer teardown()
	clock := TestClock{}

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
				mma, err := addStepTx(tt.args.tx, tt.args.schoolId, tt.args.currencyId, tt.args.amount)
				got, _ := addStepHelperTx(tt.args.tx, tt.args.schoolId, tt.args.currencyId, clock.Now(), mma, &clock)
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

		mma, err := addStepTx(tx, schools[0], teachers[0], 10)
		stepTx, _ := addStepHelperTx(tx, schools[0], teachers[0], clock.Now(), mma, &clock)
		require.Nil(t, err)
		require.Equal(t, 10.0, stepTx.InexactFloat64())

		err = updateXRatesTx(accounts, &clock)
		require.Nil(t, err)

		mma, err = addStepTx(tx, schools[0], teachers[0], 10)
		stepTx, _ = addStepHelperTx(tx, schools[0], teachers[0], clock.Now(), mma, &clock)
		require.Nil(t, err)
		require.Equal(t, 10.0, stepTx.InexactFloat64())
		err = updateXRatesTx(accounts, &clock)
		require.Nil(t, err)

		rx, err := xRateToBaseHistoricalRx(tx, schools[0], teachers[0], "")
		require.Nil(t, err)
		require.Equal(t, 1.0, rx.InexactFloat64())

		mma, err = addStepTx(tx, schools[0], teachers[1], 100)
		stepTx, _ = addStepHelperTx(tx, schools[0], teachers[1], clock.Now(), mma, &clock)
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

	_, _, teachers, _, students, _ := CreateTestAccounts(db, 2, 2, 2, 3)

	userInfo, _ := getUserInLocalStore(db, students[0])
	err := addUbuck2Student(db, &clock, userInfo, decimal.NewFromFloat(1.0), "daily payment")
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

func Test_ModMMA(t *testing.T) {
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("modMMA")
	defer dbTearDown()

	_, _, teachers, _, students, _ := CreateTestAccounts(db, 1, 2, 2, 3)

	require.True(t, StudentNetWorth(db, students[0]).Equal(StudentNetWorth(db, students[1])))

	student0, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	err = addBuck2Student(db, &clock, student0, decimal.NewFromFloat(100), teachers[0], "pre load")
	require.Nil(t, err)
	clock.TickOne(time.Millisecond * 1)
	err = addBuck2Student(db, &clock, student0, decimal.NewFromFloat(100), teachers[0], "pre load")
	require.Nil(t, err)
	clock.TickOne(time.Millisecond * 1)
	err = addBuck2Student(db, &clock, student0, decimal.NewFromFloat(100), teachers[0], "pre load")
	require.Nil(t, err)

	student1, err := getUserInLocalStore(db, students[1])
	require.Nil(t, err)
	err = addBuck2Student(db, &clock, student1, decimal.NewFromFloat(100), teachers[1], "pre load")
	require.Nil(t, err)
	clock.TickOne(time.Hour * 48)
	err = addBuck2Student(db, &clock, student1, decimal.NewFromFloat(100), teachers[1], "pre load")
	require.Nil(t, err)
	clock.TickOne(time.Hour * 48)
	err = addBuck2Student(db, &clock, student1, decimal.NewFromFloat(100), teachers[1], "pre load")
	require.Nil(t, err)

	err = addBuck2Student(db, &clock, student0, decimal.NewFromFloat(100), teachers[0], "pre load")
	require.Nil(t, err)
	clock.TickOne(time.Hour * 48)
	err = addBuck2Student(db, &clock, student1, decimal.NewFromFloat(100), teachers[1], "pre load")
	require.Nil(t, err)

	temp0 := StudentNetWorth(db, students[0])
	temp1 := StudentNetWorth(db, students[1])

	require.True(t, temp0.LessThan(temp1))
}
