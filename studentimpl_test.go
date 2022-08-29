package main

import (
	"encoding/json"
	"math"
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

	event := eventRequest{
		Positive:    false,
		Description: "Pay Taxes",
		Title:       "Taxes",
	}

	marshal, _ := json.Marshal(event)

	err := createJobOrEvent(db, marshal, KeyNEvents, "Teacher")
	require.Nil(t, err)

	event = eventRequest{
		Positive:    true,
		Description: "Pay Taxes",
		Title:       "Taxes",
	}

	marshal, _ = json.Marshal(event)

	err = createJobOrEvent(db, marshal, KeyPEvents, "Teacher")
	require.Nil(t, err)

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

	job := Job{
		Pay:         53000,
		Description: "Teach Stuff",
		College:     true,
	}

	marshal, err := json.Marshal(job)

	err = createJobOrEvent(db, marshal, KeyCollegeJobs, "Teacher")
	require.Nil(t, err)

	r := CollegeIfNeeded(db, &clock, student)
	require.False(t, r)

	body := openapi.UsersUserBody{
		College: true,
	}

	err = userEdit(db, &clock, student, body)
	require.Nil(t, err)

	r = CollegeIfNeeded(db, &clock, student)
	require.False(t, r)

	clock.TickOne(time.Hour * 720)

	r = CollegeIfNeeded(db, &clock, student)
	require.True(t, r)

	r = CollegeIfNeeded(db, &clock, student)
	require.False(t, r)

	student, err = getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	require.NotEqual(t, student.Name, "")
}

func TestCareer(t *testing.T) {

	lgr.Printf("INFO TestEvents")
	t.Log("INFO TestEvents")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("-event")
	defer dbTearDown()
	_, _, _, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 10)

	student, _ := getUserInLocalStore(db, students[0])

	job := Job{
		Pay:         53000,
		Description: "Teach Stuff",
		College:     true,
	}

	marshal, err := json.Marshal(job)

	err = createJobOrEvent(db, marshal, KeyCollegeJobs, "Teacher")
	require.Nil(t, err)

	job2 := Job{
		Pay:         53000,
		Description: "Teach Stuff",
		College:     false,
	}

	marshal, err = json.Marshal(job2)

	err = createJobOrEvent(db, marshal, KeyJobs, "Teacher")
	require.Nil(t, err)

	r := CareerIfNeeded(db, &clock, student)
	require.False(t, r)

	body := openapi.UsersUserBody{
		CareerTransition: true,
	}

	err = userEdit(db, &clock, student, body)
	require.Nil(t, err)

	r = CareerIfNeeded(db, &clock, student)
	require.False(t, r)

	clock.TickOne(time.Hour * 720)

	r = CareerIfNeeded(db, &clock, student)
	require.True(t, r)

	r = CareerIfNeeded(db, &clock, student)
	require.False(t, r)

	student, err = getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	require.NotEqual(t, student.Name, "")
}

func TestDailyPayment(t *testing.T) {

	lgr.Printf("INFO TestDailyPayment")
	t.Log("INFO TestDailyPayment")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("-pay")
	defer dbTearDown()
	_, _, _, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 1)

	student, _ := getUserInLocalStore(db, students[0])

	r := DailyPayIfNeeded(db, &clock, student)
	require.True(t, r)

	require.Equal(t, float64(student.Income), StudentNetWorth(db, student.Name).InexactFloat64())

	r = DailyPayIfNeeded(db, &clock, student)
	require.False(t, r)

	clock.Tick()
	r = DailyPayIfNeeded(db, &clock, student)
	require.False(t, r)

	clock.TickOne(24 * time.Hour)
	r = DailyPayIfNeeded(db, &clock, student)
	require.True(t, r)

	netWorth := decimal.Zero
	_ = db.View(func(tx *bolt.Tx) error {
		netWorth = StudentNetWorthTx(tx, students[0])
		return nil
	})

	require.True(t, netWorth.GreaterThan(decimal.NewFromInt(200)))

	err := chargeStudentUbuck(db, &clock, student, decimal.NewFromInt(2000), "debt", false)
	require.Nil(t, err)

	clock.TickOne(24 * time.Hour)
	r = DailyPayIfNeeded(db, &clock, student)
	require.True(t, r)

	var account openapi.ResponseCurrencyExchange
	err = db.View(func(tx *bolt.Tx) error {
		studentBucket, err := getStudentBucketRx(tx, student.Name)
		if err != nil {
			return err
		}
		accounts := studentBucket.Bucket([]byte(KeyAccounts))
		debt := accounts.Bucket([]byte(KeyDebt))
		account, err = getStudentAccountRx(tx, debt, student.Name)
		if err != nil {
			return err
		}
		return nil
	})

	require.Nil(t, err)

	require.Equal(t, math.Round(float64(student.Income*.3)*1000)/1000, math.Round(float64((2000-account.Balance)*1000))/1000)

}

func TestDailyPaymentGarnish(t *testing.T) {

	lgr.Printf("INFO TestDailyPaymentGarnish")
	t.Log("INFO TestDailyPaymentGarnish")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("payGarnish")
	defer dbTearDown()
	_, _, _, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 1)

	student, _ := getUserInLocalStore(db, students[0])

	r := DailyPayIfNeeded(db, &clock, student)
	require.True(t, r)

	require.Equal(t, float64(student.Income), StudentNetWorth(db, student.Name).InexactFloat64())

	chargeStudentUbuck(db, &clock, student, decimal.NewFromFloat32(student.Income+1), "to debt", false)
	chargeStudent(db, &clock, student, decimal.NewFromFloat32(student.Income), KeyDebt, "minimize debt", false)
	netWorth := StudentNetWorth(db, student.Name)
	require.Equal(t, float64(student.Income-1), netWorth.InexactFloat64())

	clock.TickOne(24 * time.Hour)
	r = DailyPayIfNeeded(db, &clock, student)
	require.True(t, r)

	var account openapi.ResponseCurrencyExchange
	err := db.View(func(tx *bolt.Tx) error {
		studentBucket, err := getStudentBucketRx(tx, student.Name)
		if err != nil {
			return err
		}
		accounts := studentBucket.Bucket([]byte(KeyAccounts))
		debt := accounts.Bucket([]byte(KeyDebt))
		account, err = getStudentAccountRx(tx, debt, student.Name)
		if err != nil {
			return err
		}
		return nil
	})

	require.Nil(t, err)

	require.Equal(t, float64(0), math.Round(float64((account.Balance)*1000))/1000)

}

func TestDebtInterest(t *testing.T) {

	lgr.Printf("INFO TestDebtInterest")
	t.Log("INFO TestDebtInterest")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("debtInterest")
	defer dbTearDown()
	_, _, _, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 1)

	student, _ := getUserInLocalStore(db, students[0])

	r := DebtIfNeeded(db, &clock, student)
	require.False(t, r)

	r = DailyPayIfNeeded(db, &clock, student)
	require.True(t, r)

	require.Equal(t, float64(student.Income), StudentNetWorth(db, student.Name).InexactFloat64())

	err := chargeStudentUbuck(db, &clock, student, decimal.NewFromInt(2000), "debt", false)
	require.Nil(t, err)

	r = DebtIfNeeded(db, &clock, student)
	require.False(t, r)

	clock.Tick()
	r = DebtIfNeeded(db, &clock, student)
	require.False(t, r)

	clock.TickOne(24 * time.Hour)
	r = DebtIfNeeded(db, &clock, student)
	require.True(t, r)

	var account openapi.ResponseCurrencyExchange
	err = db.View(func(tx *bolt.Tx) error {
		studentBucket, err := getStudentBucketRx(tx, student.Name)
		if err != nil {
			return err
		}
		accounts := studentBucket.Bucket([]byte(KeyAccounts))
		debt := accounts.Bucket([]byte(KeyDebt))
		account, err = getStudentAccountRx(tx, debt, student.Name)
		if err != nil {
			return err
		}
		return nil
	})

	require.Nil(t, err)

	require.Greater(t, account.Balance, float32(2000))
}

func TestGetCryptoForStudentRequest(t *testing.T) {

	lgr.Printf("INFO TestGetCryptoForStudentRequest")
	t.Log("INFO TestGetCryptoForStudentRequest")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("getCryptoForStudentRequest")
	defer dbTearDown()
	coinGecko(db)
	_, _, _, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 1)

	student, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	err = pay2Student(db, &clock, student, decimal.NewFromFloat(10000), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	resp, err := getCryptoForStudentRequest(db, student, "bitCoin")
	require.Nil(t, err)

	resp2, err := getCryptoForStudentRequest(db, student, "bitCoin")
	require.Nil(t, err)
	require.Equal(t, resp.Usd, resp2.Usd)

	clock.TickOne(time.Minute * 2)

	resp, err = getCryptoForStudentRequest(db, student, "bitCoin")

	var bitcoin openapi.CryptoCb
	err = db.View(func(tx *bolt.Tx) error {
		cryptos := tx.Bucket([]byte(KeyCryptos))
		bitcoinData := cryptos.Get([]byte("bitcoin"))
		err = json.Unmarshal(bitcoinData, &bitcoin)

		return err
	})

	require.Nil(t, err)

	require.Less(t, time.Now().Truncate(time.Second).Sub(bitcoin.UpdatedAt), time.Second*5)
}

func TestCryptoTransaction(t *testing.T) {

	lgr.Printf("INFO TestCryptoTransaction")
	t.Log("INFO TestCryptoTransaction")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("cryptoTransaction")
	defer dbTearDown()
	coinGecko(db)
	_, _, _, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 1)

	student, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	err = pay2Student(db, &clock, student, decimal.NewFromFloat(10000), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	body := openapi.RequestCryptoConvert{
		Name: "cardano",
		Buy:  10,
		Sell: 0,
	}

	resp, err := getCryptoForStudentRequest(db, student, "Cardano")
	require.Nil(t, err)
	require.Equal(t, float32(0), resp.Owned)

	err = cryptoTransaction(db, &clock, student, body)

	resp, err = getCryptoForStudentRequest(db, student, "Cardano")
	require.Nil(t, err)
	require.Equal(t, float32(10), resp.Owned)

	resp2, _, err := getStudentCrypto(db, student, "cardano")
	require.Nil(t, err)
	require.NotZero(t, resp2.Basis)

	clock.TickOne(time.Minute * 2)

	body.Buy = 0
	body.Sell = 11

	err = cryptoTransaction(db, &clock, student, body)
	require.NotNil(t, err)

	body.Sell = 5

	err = cryptoTransaction(db, &clock, student, body)
	require.Nil(t, err)
	resp, err = getCryptoForStudentRequest(db, student, "Cardano")
	require.Nil(t, err)
	require.Equal(t, float32(5), resp.Owned)

	err = cryptoTransaction(db, &clock, student, body)
	require.Nil(t, err)
	resp, err = getCryptoForStudentRequest(db, student, "Cardano")
	require.Nil(t, err)
	require.Equal(t, float32(0), resp.Owned)

}
