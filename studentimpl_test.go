package main

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func Test_ubuckFlow(t *testing.T) {

	clock := AppClock{}
	db, dbTearDown := OpenTestDB("")
	defer dbTearDown()

	_, _, teachers, _, students, err := CreateTestAccounts(db, 2, 2, 2, 3)

	userInfo, _ := getUserInLocalStore(db, students[0])
	err = addUbuck2Student(db, &clock, userInfo, decimal.NewFromFloat(1.01), "daily payment")
	require.Nil(t, err)

	balance := StudentNetWorth(db, students[0])
	require.Equal(t, 1.01, balance.InexactFloat64())

	err = chargeStudentUbuck(db, &clock, userInfo, decimal.NewFromFloat(0.51), "some reason", false)

	balance = StudentNetWorth(db, students[0])
	require.Equal(t, 0.5, balance.InexactFloat64())

	err = chargeStudentUbuck(db, &clock, userInfo, decimal.NewFromFloat(0.51), "some reason", false)
	require.Nil(t, err)
	balance = StudentNetWorth(db, students[0])
	require.Equal(t, -0.01, balance.InexactFloat64())

	err = pay2Student(db, &clock, userInfo, decimal.NewFromFloat(0.48), teachers[0], "reward")
	require.Nil(t, err)
	balance = StudentNetWorth(db, students[0])
	require.Equal(t, .47, balance.InexactFloat64())

	err = pay2Student(db, &clock, userInfo, decimal.NewFromFloat(0.01), teachers[1], "some reason")
	require.Nil(t, err)

	err = chargeStudent(db, &clock, userInfo, decimal.NewFromFloat(1), teachers[1], "some reason", false)
	require.Nil(t, err)
	balance = StudentNetWorth(db, students[0])
	require.Equal(t, -24.01999984, balance.InexactFloat64())

}
