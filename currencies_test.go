package main

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"
	"testing"
)

func TestSchool_xRateTx(t *testing.T) {
	type fields struct {
		db       *bolt.DB
		schoolId string
	}
	type args struct {
		tx       *bolt.Tx
		schoolId string
		from     string
		to       string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := xRateFromToBaseTx(tt.args.tx, tt.args.schoolId, tt.args.from, tt.args.to)
			if !tt.wantErr(t, err, fmt.Sprintf("xRateTx(%v, %v, %v, %v)", tt.args.tx, tt.args.schoolId, tt.args.from, tt.args.to)) {
				return
			}
			assert.Equalf(t, tt.want, got, "xRateTx(%v, %v, %v, %v)", tt.args.tx, tt.args.schoolId, tt.args.from, tt.args.to)
		})
	}
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
					clock:      &clock,
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
					clock:      &clock,
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
				got, err := addStepTx(tt.args.tx, tt.args.schoolId,, tt.args.currencyId, tt.args.amount)
				if !tt.wantErr(t, err, fmt.Sprintf("addStepTx(%v, %v, %v, %v, %v)", tt.args.tx, tt.args.schoolId, tt.args.clock, tt.args.currencyId, tt.args.amount)) {
					return
				}
				assert.True(t, tt.want.Sub(got).LessThan(decimal.NewFromFloat(0.001)), "addStepTx(%v, %v, %v, %v, %v - %v %v)", tt.args.tx, tt.args.schoolId, tt.args.clock, tt.args.currencyId, tt.args.amount, tt.want.InexactFloat64(), got.InexactFloat64())
			})
			clock.Tick()
		}

		mmaTx, err2 := getCurrencyMMATx(tx, schools[0], teachers[0])
		require.Nil(t, err2)
		assert.True(t, mmaTx.Sub(decimal.NewFromFloat32(10.90)).LessThan(decimal.NewFromFloat(0.001)))
		return nil
	})
}

func Test_convert(t *testing.T) {
	type args struct {
		schoolId string
		from     string
		to       string
		amount   float64
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convert(tt.args.schoolId, tt.args.from, tt.args.to, tt.args.amount)
			if !tt.wantErr(t, err, fmt.Sprintf("convert(%v, %v, %v, %v)", tt.args.schoolId, tt.args.from, tt.args.to, tt.args.amount)) {
				return
			}
			assert.Equalf(t, tt.want, got, "convert(%v, %v, %v, %v)", tt.args.schoolId, tt.args.from, tt.args.to, tt.args.amount)
		})
	}
}
