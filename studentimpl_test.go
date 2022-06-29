package main

import (
	"testing"
	"time"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/go-pkgz/lgr"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"
)

func Test_ubuckFlow(t *testing.T) {

	clock := TestClock{}
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

func TestEvents(t *testing.T) {

	lgr.Printf("INFO TestEvents")
	t.Log("INFO TestEvents")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("-event")
	defer dbTearDown()
	_, _, _, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 10)

	student, _ := getUserInLocalStore(db, students[0])

	for _, student := range students {
		studentDetails, _ := getUserInLocalStore(db, student)
		err := addUbuck2Student(db, &clock, studentDetails, decimal.NewFromFloat(100), "pre load")
		require.Nil(t, err)

	}

	r := EventIfNeeded(db, &clock, student)
	require.False(t, r)

	r = EventIfNeeded(db, &clock, student)
	require.False(t, r)

	clock.TickOne(time.Hour * 216)

	r = EventIfNeeded(db, &clock, student)
	require.True(t, r)

	r = EventIfNeeded(db, &clock, student)
	require.False(t, r)

	netWorth := decimal.Zero
	_ = db.View(func(tx *bolt.Tx) error {
		netWorth = StudentNetWorthTx(tx, students[0])
		return nil
	})

	require.False(t, netWorth.Equal(decimal.NewFromFloat(100)))
}

func TestCollege(t *testing.T) {

	lgr.Printf("INFO TestCollege")
	t.Log("INFO TestCollege")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("-college")
	defer dbTearDown()
	_, _, _, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 1)

	student, _ := getUserInLocalStore(db, students[0])

	r := CollegeIfNeeded(db, &clock, student)
	require.False(t, r)

	body := openapi.UsersUserBody{
		College: true,
	}

	err := userEdit(db, &clock, student, body)
	require.Nil(t, err)

	r = CollegeIfNeeded(db, &clock, student)
	require.False(t, r)

	clock.TickOne(time.Hour * 720)

	r = CollegeIfNeeded(db, &clock, student)
	require.True(t, r)

	r = CollegeIfNeeded(db, &clock, student)
	require.False(t, r)
}

func TestCareer(t *testing.T) {

	lgr.Printf("INFO TestEvents")
	t.Log("INFO TestEvents")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("-event")
	defer dbTearDown()
	_, _, _, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 10)

	student, _ := getUserInLocalStore(db, students[0])

	r := CareerIfNeeded(db, &clock, student)
	require.False(t, r)

	body := openapi.UsersUserBody{
		CareerTransition: true,
	}

	err := userEdit(db, &clock, student, body)
	require.Nil(t, err)

	r = CareerIfNeeded(db, &clock, student)
	require.False(t, r)

	clock.TickOne(time.Hour * 720)

	r = CareerIfNeeded(db, &clock, student)
	require.True(t, r)

	r = CareerIfNeeded(db, &clock, student)
	require.False(t, r)
}
