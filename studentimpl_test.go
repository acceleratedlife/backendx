package main

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
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

	_, _, teachers, _, students, _ := CreateTestAccounts(db, 2, 2, 2, 3)

	userInfo, _ := getUserInLocalStore(db, students[0])
	err := addUbuck2Student(db, &clock, userInfo, decimal.NewFromFloat(1.01), "daily payment")
	require.Nil(t, err)

	balance := StudentNetWorth(db, students[0])
	require.Equal(t, 1.01, balance.InexactFloat64())

	_ = chargeStudentUbuck(db, &clock, userInfo, decimal.NewFromFloat(0.51), "some reason", false)

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

	event := EventRequest{
		Positive:    false,
		Description: "Pay Taxes",
		Title:       "Taxes",
	}

	marshal, _ := json.Marshal(event)

	err := createJobOrEvent(db, marshal, KeyNEvents, "Negative")
	require.Nil(t, err)

	event = EventRequest{
		Positive:    true,
		Description: "Tax Refund",
		Title:       "Taxes",
	}

	marshal, _ = json.Marshal(event)

	err = createJobOrEvent(db, marshal, KeyPEvents, "Positive")
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

func TestEventsLowUbuck(t *testing.T) {

	lgr.Printf("INFO TestEventsLowUbuck")
	t.Log("INFO TestEventsLowUbuck")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("-eventLowUbuck")
	defer dbTearDown()
	_, _, _, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 10)

	student, _ := getUserInLocalStore(db, students[0])
	student2, _ := getUserInLocalStore(db, students[2])

	event := EventRequest{
		Positive:    false,
		Description: "Pay Taxes",
		Title:       "Taxes",
	}

	marshal, _ := json.Marshal(event)

	err := createJobOrEvent(db, marshal, KeyNEvents, "Negative")
	require.Nil(t, err)

	event = EventRequest{
		Positive:    true,
		Description: "Tax Refund",
		Title:       "Taxes",
	}

	marshal, _ = json.Marshal(event)

	err = createJobOrEvent(db, marshal, KeyPEvents, "Positive")
	require.Nil(t, err)

	for _, student := range students {
		studentDetails, _ := getUserInLocalStore(db, student)
		err := addUbuck2Student(db, &clock, studentDetails, decimal.NewFromFloat(1), "pre load")
		require.Nil(t, err)

	}

	err = addUbuck2Student(db, &clock, student, decimal.NewFromFloat(1000), "overload load")
	require.Nil(t, err)

	r := EventIfNeeded(db, &clock, student2)
	require.False(t, r)

	clock.TickOne(time.Hour * 216)

	r = EventIfNeeded(db, &clock, student2)
	require.True(t, r)

	keys := 0
	err = db.View(func(tx *bolt.Tx) error {
		cb, err := getCbRx(tx, student2.SchoolId)
		if err != nil {
			return err
		}

		accountsBucket := cb.Bucket([]byte(KeyAccounts))
		if accountsBucket == nil {
			return fmt.Errorf("cannot find cb accounts")
		}

		debtBucket := accountsBucket.Bucket([]byte(KeyDebt))
		if debtBucket == nil {
			return nil
		}

		trans := debtBucket.Bucket([]byte(KeyTransactions))
		if trans == nil {
			return fmt.Errorf("cannot find cb account debt trans")
		}

		keys = trans.Stats().KeyN //KeyN is failing elsewhere so I need to check that it is working true
		return nil
	})

	require.Nil(t, err)

	require.LessOrEqual(t, keys, 1)
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

	marshal, _ := json.Marshal(job)

	err := createJobOrEvent(db, marshal, KeyCollegeJobs, "Teacher")
	require.Nil(t, err)

	r := CollegeIfNeeded(db, &clock, student)
	require.False(t, r)

	body := openapi.RequestUserEdit{
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

	marshal, _ := json.Marshal(job)

	err := createJobOrEvent(db, marshal, KeyCollegeJobs, "Teacher")
	require.Nil(t, err)

	job2 := Job{
		Pay:         53000,
		Description: "Teach Stuff",
		College:     false,
	}

	marshal, _ = json.Marshal(job2)

	err = createJobOrEvent(db, marshal, KeyJobs, "Teacher")
	require.Nil(t, err)

	r := CareerIfNeeded(db, &clock, student)
	require.False(t, r)

	body := openapi.RequestUserEdit{
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

	clock.TickOne(5 * time.Hour)
	r = DailyPayIfNeeded(db, &clock, student)
	require.False(t, r)

	clock.TickOne(3 * time.Hour)
	r = DailyPayIfNeeded(db, &clock, student)
	require.False(t, r)

	clock.TickOne(1 * time.Hour)
	r = DailyPayIfNeeded(db, &clock, student)
	require.True(t, r)

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

	require.Equal(t, math.Round(float64(student.Income*.75)*1000)/1000, math.Round(float64((2000-account.Balance)*1000))/1000)

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

	_ = cryptoTransaction(db, &clock, student, body)

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

func TestTrueAuctionFalse(t *testing.T) {

	lgr.Printf("INFO TestTrueAuctionFalse")
	t.Log("INFO TestTrueAuctionFalse")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("TrueAuctionFalse")
	defer dbTearDown()
	_, _, teachers, classes, students, err := CreateTestAccounts(db, 1, 1, 1, 2)
	require.Nil(t, err)

	SetTestLoginUser(teachers[0])

	student0, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	student1, err := getUserInLocalStore(db, students[1])
	require.Nil(t, err)
	teacher, err := getUserInLocalStore(db, teachers[0])
	require.Nil(t, err)

	err = addUbuck2Student(db, &clock, student0, decimal.NewFromFloat(100), "starter")
	require.Nil(t, err)
	err = addUbuck2Student(db, &clock, student1, decimal.NewFromFloat(100), "starter")
	require.Nil(t, err)

	err = MakeAuctionImpl(db, teacher, openapi.RequestMakeAuction{
		Bid:         0,
		MaxBid:      0,
		Description: "test auc",
		EndDate:     time.Now().Add(time.Minute),
		StartDate:   time.Now(),
		OwnerId:     teacher.Name,
		Visibility:  classes,
		TrueAuction: false,
	}, true)
	require.Nil(t, err)

	auctions, err := getTeacherAuctions(db, teacher)
	require.Nil(t, err)

	//overdrawn student 0 max 0 bid 0
	timeId := auctions[0].Id.Format(time.RFC3339Nano)

	clock.TickOne(time.Second * 30)

	_, err = placeBid(db, &clock, student1, timeId, 1)
	require.Nil(t, err)

	auctions, err = getTeacherAuctions(db, teacher)
	require.Nil(t, err)

	require.Equal(t, timeId, auctions[0].EndDate.Format((time.RFC3339Nano)))

}

func TestTrueAuctionTrue(t *testing.T) {

	lgr.Printf("INFO TestTrueAuctionFalse")
	t.Log("INFO TestTrueAuctionFalse")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("TrueAuctionFalse")
	defer dbTearDown()
	_, _, teachers, classes, students, err := CreateTestAccounts(db, 1, 1, 1, 2)
	require.Nil(t, err)

	SetTestLoginUser(teachers[0])

	student0, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	student1, err := getUserInLocalStore(db, students[1])
	require.Nil(t, err)
	teacher, err := getUserInLocalStore(db, teachers[0])
	require.Nil(t, err)

	err = addUbuck2Student(db, &clock, student0, decimal.NewFromFloat(100), "starter")
	require.Nil(t, err)
	err = addUbuck2Student(db, &clock, student1, decimal.NewFromFloat(100), "starter")
	require.Nil(t, err)

	err = MakeAuctionImpl(db, teacher, openapi.RequestMakeAuction{
		Bid:         0,
		MaxBid:      0,
		Description: "test auc",
		EndDate:     time.Now().Add(time.Minute),
		StartDate:   time.Now(),
		OwnerId:     teacher.Name,
		Visibility:  classes,
		TrueAuction: true,
	}, true)
	require.Nil(t, err)

	auctions, err := getTeacherAuctions(db, teacher)
	require.Nil(t, err)

	//overdrawn student 0 max 0 bid 0
	timeId := auctions[0].Id.Format(time.RFC3339Nano)

	clock.TickOne(time.Second * 30)

	_, err = placeBid(db, &clock, student1, timeId, 1)
	require.Nil(t, err)

	auctions, err = getTeacherAuctions(db, teacher)
	require.Nil(t, err)

	require.NotEqual(t, timeId, auctions[0].EndDate.Format((time.RFC3339Nano)))

}

func TestPurchaseLottoOff(t *testing.T) {

	lgr.Printf("INFO TestPurchaseLottoOff")
	t.Log("INFO TestPurchaseLottoOff")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("purchaseLottoOff")
	defer dbTearDown()
	_, _, _, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 1)

	student, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	err = pay2Student(db, &clock, student, decimal.NewFromFloat(10000), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	_, err = purchaseLotto(db, &clock, student, 5)
	require.Equal(t, "the lotto has not been initialized", err.Error())

	var settings = openapi.Settings{
		Lottery: true,
		Odds:    20,
	}

	err = setSettings(db, &clock, student, settings)
	require.Nil(t, err)

	var settings2 = openapi.Settings{
		Lottery: false,
		Odds:    20,
	}

	err = setSettings(db, &clock, student, settings2)
	require.Nil(t, err)

	winner, err := purchaseLotto(db, &clock, student, 60)
	require.Nil(t, err)
	require.True(t, winner)

	_, err = purchaseLotto(db, &clock, student, 60)
	require.NotNil(t, err)

	err = setSettings(db, &clock, student, settings)
	require.Nil(t, err)

	winner, err = purchaseLotto(db, &clock, student, 60)
	require.Nil(t, err)
	require.True(t, winner)
}

func TestPurchaseLottoSingle(t *testing.T) {

	lgr.Printf("INFO TestPurchaseLottoSingle")
	t.Log("INFO TestPurchaseLottoSingle")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("purchaseLottoSingle")
	defer dbTearDown()
	_, _, _, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 1)

	student, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	err = pay2Student(db, &clock, student, decimal.NewFromFloat(10000), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	var settings = openapi.Settings{
		Lottery: true,
		Odds:    30,
	}

	err = setSettings(db, &clock, student, settings)
	require.Nil(t, err)

	winner, err := purchaseLotto(db, &clock, student, 1)
	require.Nil(t, err)
	count := 0

	for !winner {
		winner, err = purchaseLotto(db, &clock, student, 1)
		require.Nil(t, err)
		count++
		lgr.Printf(strconv.Itoa(count))
	}

	require.True(t, winner)
}

func TestBuyCDFunction(t *testing.T) {

	lgr.Printf("INFO TestBuyCDFunction")
	t.Log("INFO TestBuyCDFunction")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("buyCDFunction")
	defer dbTearDown()
	_, _, _, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 1)

	student, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	err = pay2Student(db, &clock, student, decimal.NewFromFloat(199), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	body := openapi.RequestBuyCd{
		PrinInv: 100,
		Time:    60,
	}
	err = buyCD(db, &clock, student, body)
	require.Nil(t, err)
	err = buyCD(db, &clock, student, body)
	require.NotNil(t, err)
}

func TestGetCDSFunction(t *testing.T) {

	lgr.Printf("INFO TestGetCDSFunction")
	t.Log("INFO TestGetCDSFunction")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("getCDSFunction")
	defer dbTearDown()
	_, _, _, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 1)

	student, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	err = pay2Student(db, &clock, student, decimal.NewFromFloat(1000), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	body := openapi.RequestBuyCd{
		PrinInv: 100,
		Time:    60,
	}

	for i := 0; i < 5; i++ {
		err = buyCD(db, &clock, student, body)
		require.Nil(t, err)
	}

	resp, err := getCDS(db, student)
	require.Nil(t, err)

	require.Equal(t, 5, len(resp))

}

func TestGetCDTransaction(t *testing.T) {

	lgr.Printf("INFO TestGetCDTransaction")
	t.Log("INFO TestGetCDTransaction")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("getCDTransaction")
	defer dbTearDown()
	_, _, _, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 1)

	student, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	err = pay2Student(db, &clock, student, decimal.NewFromFloat(10000), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	body := openapi.RequestBuyCd{
		PrinInv: 100,
		Time:    60,
	}

	for i := 0; i < 5; i++ {
		err = buyCD(db, &clock, student, body)
		require.Nil(t, err)
	}

	resp, err := getCDTransactions(db, student)
	require.Nil(t, err)

	require.Equal(t, 5, len(resp))

	for i := 0; i < 30; i++ {
		err = buyCD(db, &clock, student, body)
		require.Nil(t, err)
	}

	resp, err = getCDTransactions(db, student)
	require.Nil(t, err)

	require.Equal(t, 25, len(resp))

}

func TestRefundCDFunction(t *testing.T) {

	lgr.Printf("INFO TestRefundCDFunction")
	t.Log("INFO TestRefundCDFunction")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("refundCDFunction")
	defer dbTearDown()
	_, _, _, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 1)

	student, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	err = pay2Student(db, &clock, student, decimal.NewFromFloat(500), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	body := openapi.RequestBuyCd{
		PrinInv: 100,
		Time:    60,
	}

	for i := 0; i < 5; i++ {
		err = buyCD(db, &clock, student, body)
		require.Nil(t, err)
	}

	resp, err := getCDS(db, student)
	require.Nil(t, err)

	require.Equal(t, 5, len(resp))

	trans, err := getCDTransactions(db, student)
	require.Nil(t, err)

	require.Equal(t, 5, len(trans))

	err = refundCD(db, &clock, student, resp[0].Ts.Format(time.RFC3339Nano))
	require.Nil(t, err)

	resp, err = getCDS(db, student)
	require.Nil(t, err)

	require.Equal(t, 4, len(resp))

	trans, err = getCDTransactions(db, student)
	require.Nil(t, err)

	require.Equal(t, 6, len(trans))

	netWorth := StudentNetWorth(db, student.Email)

	require.True(t, netWorth.Equal(decimal.NewFromInt32(450)))

	//I think this test is working correctly but I need another one that uses multiple different lengths and I need to make the time change with tick
}
