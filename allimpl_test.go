package main

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func Test_student2student(t *testing.T) {

	clock := TestClock{}
	db, dbTearDown := OpenTestDB("student2student")
	defer dbTearDown()

	_, _, _, _, students, _ := CreateTestAccounts(db, 1, 2, 2, 3)

	student0, _ := getUserInLocalStore(db, students[0])
	err := addUbuck2Student(db, &clock, student0, decimal.NewFromFloat(200), "daily payment")
	require.Nil(t, err)

	err = executeStudentTransaction(db, &clock, 75, students[1], student0, "Pencil")
	require.Nil(t, err)

	balance := StudentNetWorth(db, students[0])
	require.Equal(t, float64(124.25), balance.InexactFloat64())

	balance = StudentNetWorth(db, students[1])
	require.Equal(t, float64(75), balance.InexactFloat64())
}

func TestClearMessagesImpl(t *testing.T) {
	db, closeDB := OpenTestDB("")
	defer closeDB()

	_, _, _, _, students, err := CreateTestAccounts(db, 1, 1, 1, 1)
	require.NoError(t, err)

	theMessage := "Hello World"
	err = message(db, &theMessage, &students[0], nil, false, false)
	require.NoError(t, err)

	user, err := getUserInLocalStore(db, students[0])
	require.NoError(t, err)
	require.Len(t, user.Messages, 1)
	require.Equal(t, theMessage, user.Messages[0])

	err = clearMessages(db, user)
	require.NoError(t, err)

	user, err = getUserInLocalStore(db, students[0])
	require.NoError(t, err)
	require.Len(t, user.Messages, 0)

}
