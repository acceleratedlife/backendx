package main

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"
	"testing"
)

func Test_addUbuck2Student(t *testing.T) {
	clock := AppClock{}
	db, dbTearDown := OpenTestDB("")
	defer dbTearDown()

	_, _, _, _, students, err := CreateTestAccounts(db, 2, 2, 2, 3)

	require.Nil(t, err)
	require.Equal(t, 24, len(students))

	for _, s := range students {
		err := addUbuck2Student(db, clock, s, decimal.NewFromInt(1), "daily payment")
		require.Nil(t, err)
	}

	var balance decimal.Decimal

	_ = db.View(func(tx *bolt.Tx) error {
		cb := tx.Bucket([]byte(KeyCB))
		ub := cb.Get([]byte(CurrencyUBuck))
		return balance.UnmarshalText(ub)
	})

	assert.Equal(t, "24", balance.String())
}

//func TestStudentNetWorthTx(t *testing.T) {
//	type args struct {
//		tx       *bolt.Tx
//		userName string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		wantRes decimal.Decimal
//	}{
//		{name: , args: , wantRes: },
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			assert.Equalf(t, tt.wantRes, StudentNetWorthTx(tt.args.tx, tt.args.userName), "StudentNetWorthTx(%v, %v)", tt.args.tx, tt.args.userName)
//		})
//	}
//}
