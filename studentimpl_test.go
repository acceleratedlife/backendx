package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"testing"
	"time"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/go-pkgz/lgr"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"
)

// NoopSSEService implements SSEServiceInterface but does nothing
type NoopSSEService struct{}

func (n *NoopSSEService) BroadcastAuctionEvent(auctionID string, eventType string, data interface{}) {
	lgr.Printf("[SSE-NOOP] Would broadcast event - Type: %s, Auction ID: %s", eventType, auctionID)
}

func (n *NoopSSEService) HandleAuctionEventsSSE(w http.ResponseWriter, r *http.Request) {
	lgr.Printf("[SSE-NOOP] Would handle SSE connection from %s", r.RemoteAddr)
}

func (n *NoopSSEService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lgr.Printf("[SSE-NOOP] Would serve HTTP for %s", r.RemoteAddr)
}

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

	schoolsNetworth(db)

	r := EventIfNeeded(db, &clock, student)
	require.False(t, r)

	r = EventIfNeeded(db, &clock, student)
	require.False(t, r)

	clock.TickOne(time.Hour * 9 * 24)

	r = EventIfNeeded(db, &clock, student)
	require.True(t, r)

	r = EventIfNeeded(db, &clock, student)
	require.False(t, r)

	clock.TickOne(time.Hour * 27 * 24)

	r = EventIfNeeded(db, &clock, student)
	require.True(t, r)

	netWorth := StudentNetWorth(db, students[0])

	require.False(t, netWorth.Equal(decimal.NewFromFloat(100)))

	studentDetails, _ := getUserInLocalStore(db, students[0])
	totalEvents, err := getEventsTeacher(db, &clock, studentDetails)
	require.Nil(t, err)
	require.Greater(t, len(totalEvents), 3)

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

	schoolsNetworth(db)

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

func TestEventsGarnish(t *testing.T) {

	lgr.Printf("INFO TestEventsGarnish")
	t.Log("INFO TestEventsGarnish")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("-eventGarnish")
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

	r := DailyPayIfNeeded(db, &clock, student)
	require.True(t, r)

	err = chargeStudent(db, &clock, student, decimal.NewFromFloat(1000), CurrencyUBuck, "debt load", false)
	require.Nil(t, err)

	schoolsNetworth(db)

	r = EventIfNeeded(db, &clock, student)
	require.False(t, r)

	clock.TickOne(time.Hour * 24 * 10)

	r = EventIfNeeded(db, &clock, student)
	require.True(t, r)

	err = db.View(func(tx *bolt.Tx) error {
		studentBucket, err := getStudentBucketRx(tx, student.Name)
		require.Nil(t, err)

		_, _, balance, err := IsDebtNeededRx(studentBucket, &clock)
		require.Nil(t, err)

		require.Less(t, balance.InexactFloat64(), float64(1000))

		return err
	})

	require.Nil(t, err)

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

//****** good test but needs to be ran on its own due to coingecko chanages
// func TestGetCryptoForStudentRequest(t *testing.T) {

// 	lgr.Printf("INFO TestGetCryptoForStudentRequest")
// 	t.Log("INFO TestGetCryptoForStudentRequest")
// 	clock := TestClock{}
// 	db, dbTearDown := OpenTestDB("getCryptoForStudentRequest")
// 	defer dbTearDown()
// 	coinGecko(db)
// 	_, _, _, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 1)

// 	student, err := getUserInLocalStore(db, students[0])
// 	require.Nil(t, err)
// 	err = pay2Student(db, &clock, student, decimal.NewFromFloat(10000), CurrencyUBuck, "pre load")
// 	require.Nil(t, err)

// 	resp, err := getCryptoForStudentRequest(db, student, "bitCoin")
// 	require.Nil(t, err)

// 	resp2, err := getCryptoForStudentRequest(db, student, "bitCoin")
// 	require.Nil(t, err)
// 	require.Equal(t, resp.Usd, resp2.Usd)

// 	clock.TickOne(time.Minute * 2)

// 	resp, err = getCryptoForStudentRequest(db, student, "bitCoin")

// 	var bitcoin openapi.CryptoCb
// 	err = db.View(func(tx *bolt.Tx) error {
// 		cryptos := tx.Bucket([]byte(KeyCryptos))
// 		bitcoinData := cryptos.Get([]byte("bitcoin"))
// 		err = json.Unmarshal(bitcoinData, &bitcoin)

// 		return err
// 	})

// 	require.Nil(t, err)

// 	require.Less(t, clock.Now().Truncate(time.Second).Sub(bitcoin.UpdatedAt), time.Second*5)
// }

//****** good test but needs to be ran on its own due to coingecko chanages

// func TestCryptoTransaction(t *testing.T) {

// 	lgr.Printf("INFO TestCryptoTransaction")
// 	t.Log("INFO TestCryptoTransaction")
// 	clock := TestClock{}
// 	db, dbTearDown := OpenTestDB("cryptoTransaction")
// 	defer dbTearDown()
// 	coinGecko(db)
// 	_, _, _, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 1)

// 	student, err := getUserInLocalStore(db, students[0])
// 	require.Nil(t, err)
// 	err = pay2Student(db, &clock, student, decimal.NewFromFloat(10000), CurrencyUBuck, "pre load")
// 	require.Nil(t, err)

// 	body := openapi.RequestCryptoConvert{
// 		Name: "cardano",
// 		Buy:  10,
// 		Sell: 0,
// 	}

// 	resp, err := getCryptoForStudentRequest(db, student, "Cardano")
// 	require.Nil(t, err)
// 	require.Equal(t, float32(0), resp.Owned)

// 	_ = cryptoTransaction(db, &clock, student, body)

// 	resp, err = getCryptoForStudentRequest(db, student, "Cardano")
// 	require.Nil(t, err)
// 	require.Equal(t, float32(10), resp.Owned)

// 	resp2, _, err := getStudentCrypto(db, student, "cardano")
// 	require.Nil(t, err)
// 	require.NotZero(t, resp2.Basis)

// 	clock.TickOne(time.Minute * 2)

// 	body.Buy = 0
// 	body.Sell = 11

// 	err = cryptoTransaction(db, &clock, student, body)
// 	require.NotNil(t, err)

// 	body.Sell = 5

// 	err = cryptoTransaction(db, &clock, student, body)
// 	require.Nil(t, err)
// 	resp, err = getCryptoForStudentRequest(db, student, "Cardano")
// 	require.Nil(t, err)
// 	require.Equal(t, float32(5), resp.Owned)

// 	err = cryptoTransaction(db, &clock, student, body)
// 	require.Nil(t, err)
// 	resp, err = getCryptoForStudentRequest(db, student, "Cardano")
// 	require.Nil(t, err)
// 	require.Equal(t, float32(0), resp.Owned)

// }

//****** good test but needs to be ran on its own due to coingecko chanages

// func TestCryptoTransactionGarnish(t *testing.T) {

// 	lgr.Printf("INFO TestCryptoTransactionGarnish")
// 	t.Log("INFO TestCryptoTransactionGarnish")
// 	clock := TestClock{}
// 	db, dbTearDown := OpenTestDB("cryptoTransactionGarnish")
// 	defer dbTearDown()
// 	coinGecko(db)
// 	_, _, _, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 1)

// 	student, err := getUserInLocalStore(db, students[0])
// 	require.Nil(t, err)

// 	r := DailyPayIfNeeded(db, &clock, student)
// 	require.True(t, r)

// 	err = pay2Student(db, &clock, student, decimal.NewFromFloat(1000), CurrencyUBuck, "pre load")
// 	require.Nil(t, err)

// 	body := openapi.RequestCryptoConvert{
// 		Name: "cardano",
// 		Buy:  10,
// 		Sell: 0,
// 	}

// 	_ = cryptoTransaction(db, &clock, student, body)

// 	clock.TickOne(time.Minute * 2)

// 	err = chargeStudent(db, &clock, student, decimal.NewFromInt(2000), CurrencyUBuck, "", false)
// 	require.Nil(t, err)

// 	body.Buy = 0
// 	body.Sell = 5

// 	err = cryptoTransaction(db, &clock, student, body)
// 	require.Nil(t, err)

// 	err = db.View(func(tx *bolt.Tx) error {
// 		studentBucket, err := getStudentBucketRx(tx, student.Name)
// 		require.Nil(t, err)

// 		_, _, balance, err := IsDebtNeededRx(studentBucket, &clock)
// 		require.Nil(t, err)
// 		require.Less(t, balance.InexactFloat64(), float64(2000))
// 		return err
// 	})

// 	require.Nil(t, err)

// }

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
		EndDate:     clock.Now().Add(time.Minute),
		StartDate:   clock.Now(),
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

	_, err = placeBid(db, &clock, student1, timeId, 1, &NoopSSEService{})
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
		EndDate:     clock.Now().Add(time.Minute),
		StartDate:   clock.Now(),
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

	_, err = placeBid(db, &clock, student1, timeId, 1, &NoopSSEService{})
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

	_, err = purchaseLotto(db, &clock, student, 5)
	require.Equal(t, "the lotto has not been initialized", err.Error())

	var settings = openapi.Settings{
		Lottery: true,
		Odds:    250,
	}

	err = setSettings(db, &clock, student, settings)
	require.Nil(t, err)

	err = pay2Student(db, &clock, student, decimal.NewFromFloat(10000), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	var settings2 = openapi.Settings{
		Lottery: false,
		Odds:    250,
	}

	err = setSettings(db, &clock, student, settings2)
	require.Nil(t, err)

	winner, err := purchaseLotto(db, &clock, student, 700)
	require.Nil(t, err)
	require.True(t, winner)

	_, err = purchaseLotto(db, &clock, student, 700)
	require.NotNil(t, err)

	err = setSettings(db, &clock, student, settings)
	require.Nil(t, err)

	err = pay2Student(db, &clock, student, decimal.NewFromFloat(200000), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	winner, err = purchaseLotto(db, &clock, student, 30000)
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

	//the lotto should still have an odds of 250 as that is the lowest possible
	var settings = openapi.Settings{
		Lottery: true,
		Odds:    30,
	}

	err = setSettings(db, &clock, student, settings)
	require.Nil(t, err)

	err = pay2Student(db, &clock, student, decimal.NewFromFloat(100000), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	winner, err := purchaseLotto(db, &clock, student, 1)
	require.Nil(t, err)
	count := 0

	for !winner {
		winner, err = purchaseLotto(db, &clock, student, 1)
		require.Nil(t, err)
		count++
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
		Time:    50,
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
		Time:    50,
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
		Time:    50,
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
		Time:    50,
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

	require.True(t, netWorth.Equal(decimal.NewFromInt32(500*.75))) //.75 is early refund rate

}

func TestRefundCDWithDebtLessThan(t *testing.T) {

	lgr.Printf("INFO TestRefundCDWithDebtLessThan")
	t.Log("INFO TestRefundCDWithDebtLessThan")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("RefundCDWithDebtLessThan")
	defer dbTearDown()
	_, _, _, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 1)

	student, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	ReqUserEdit := openapi.RequestUserEdit{
		College: true,
	}

	clock.TickOne(time.Hour * 24 * 5)
	paid := DailyPayIfNeeded(db, &clock, student)
	require.True(t, paid)

	err = userEdit(db, &clock, student, ReqUserEdit)
	require.Nil(t, err)

	err = pay2Student(db, &clock, student, decimal.NewFromFloat(100000), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	body := openapi.RequestBuyCd{
		PrinInv: 100000,
		Time:    50,
	}

	err = buyCD(db, &clock, student, body)
	require.Nil(t, err)

	resp, err := getCDS(db, student)
	require.Nil(t, err)

	err = refundCD(db, &clock, student, resp[0].Ts.Format(time.RFC3339Nano))
	require.Nil(t, err)

	ubucks, err := getStudentUbuck(db, student)
	require.Nil(t, err)

	require.True(t, ubucks.Value < 90000)

}

func TestRefundCDWithDebtGreaterThan(t *testing.T) {

	lgr.Printf("INFO TestRefundCDWithDebtGreaterThan")
	t.Log("INFO TestRefundCDWithDebtGreaterThan")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("RefundCDWithDebtGreaterThan")
	defer dbTearDown()
	_, _, _, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 1)

	student, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	ReqUserEdit := openapi.RequestUserEdit{
		College: true,
	}

	clock.TickOne(time.Hour * 24 * 1)
	paid := DailyPayIfNeeded(db, &clock, student)
	require.True(t, paid)

	err = userEdit(db, &clock, student, ReqUserEdit)
	require.Nil(t, err)

	err = pay2Student(db, &clock, student, decimal.NewFromFloat(100), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	body := openapi.RequestBuyCd{
		PrinInv: 100,
		Time:    50,
	}

	err = buyCD(db, &clock, student, body)
	require.Nil(t, err)

	resp, err := getCDS(db, student)
	require.Nil(t, err)

	ubucksPre, err := getStudentUbuck(db, student)
	require.Nil(t, err)

	err = refundCD(db, &clock, student, resp[0].Ts.Format(time.RFC3339Nano))
	require.Nil(t, err)

	ubucksPost, err := getStudentUbuck(db, student)
	require.Nil(t, err)
	//this line will have to change for next term when I make it garnish 100% of a CD
	require.Equal(t, float32(body.PrinInv)*.75*(1-KeyGarnish), ubucksPost.Value-ubucksPre.Value) //.75 is early refund rate

}

func TestCDWithTimeChanges(t *testing.T) {

	lgr.Printf("INFO TestCDWithTimeChanges")
	t.Log("INFO TestCDWithTimeChanges")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("CDWithTimeChanges")
	defer dbTearDown()
	_, _, _, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 1)

	student, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	err = pay2Student(db, &clock, student, decimal.NewFromFloat(500), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	body14 := openapi.RequestBuyCd{
		PrinInv: 100,
		Time:    7,
	}
	err = buyCD(db, &clock, student, body14)
	require.Nil(t, err)

	body30 := openapi.RequestBuyCd{
		PrinInv: 100,
		Time:    14,
	}
	err = buyCD(db, &clock, student, body30)
	require.Nil(t, err)

	body50 := openapi.RequestBuyCd{
		PrinInv: 100,
		Time:    30,
	}
	err = buyCD(db, &clock, student, body50)
	require.Nil(t, err)

	body70 := openapi.RequestBuyCd{
		PrinInv: 100,
		Time:    60,
	}
	err = buyCD(db, &clock, student, body70)
	require.Nil(t, err)

	body90 := openapi.RequestBuyCd{
		PrinInv: 100,
		Time:    90,
	}
	err = buyCD(db, &clock, student, body90)
	require.Nil(t, err)

	clock.TickOne(time.Hour * 24 * 6)

	CertificateOfDepositIfNeeded(db, &clock, student)

	netWorth := StudentNetWorth(db, student.Email).InexactFloat64()

	require.True(t, netWorth > 386 && netWorth < 387)

	clock.TickOne(time.Hour * 24 * 7)

	CertificateOfDepositIfNeeded(db, &clock, student)

	netWorth = StudentNetWorth(db, student.Email).InexactFloat64()

	require.True(t, netWorth > 444 && netWorth < 445)

	clock.TickOne(time.Hour * 24 * 16)

	CertificateOfDepositIfNeeded(db, &clock, student)

	netWorth = StudentNetWorth(db, student.Email).InexactFloat64()

	require.True(t, netWorth > 618 && netWorth < 619)

	clock.TickOne(time.Hour * 24 * 30)

	CertificateOfDepositIfNeeded(db, &clock, student)

	netWorth = StudentNetWorth(db, student.Email).InexactFloat64()

	require.True(t, netWorth > 1257 && netWorth < 1258)

	clock.TickOne(time.Hour * 24 * 30)

	CertificateOfDepositIfNeeded(db, &clock, student)

	netWorth = StudentNetWorth(db, student.Email).InexactFloat64()

	require.True(t, netWorth > 3449 && netWorth < 3450)

	clock.TickOne(time.Hour * 24 * 30)

	needed := CertificateOfDepositIfNeeded(db, &clock, student)
	require.True(t, needed)
	netWorth = StudentNetWorth(db, student.Email).InexactFloat64()

	require.True(t, netWorth > 4400 && netWorth < 4401)

	needed = DailyPayIfNeeded(db, &clock, student)
	require.True(t, needed)

	needed = CertificateOfDepositIfNeeded(db, &clock, student)
	require.False(t, needed)

}

func TestCDComments(t *testing.T) {

	lgr.Printf("INFO TestCDComments")
	t.Log("INFO TestCDComments")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("CDComments")
	defer dbTearDown()
	_, _, _, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 1)

	student, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	err = pay2Student(db, &clock, student, decimal.NewFromFloat(500), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	body14 := openapi.RequestBuyCd{
		PrinInv: 100,
		Time:    14,
	}
	err = buyCD(db, &clock, student, body14)
	require.Nil(t, err)

	body30 := openapi.RequestBuyCd{
		PrinInv: 100,
		Time:    30,
	}
	err = buyCD(db, &clock, student, body30)
	require.Nil(t, err)

	clock.TickOne(time.Hour * 24 * 13)

	CertificateOfDepositIfNeeded(db, &clock, student)

	resp, err := getCDS(db, student)
	require.Nil(t, err)

	err = refundCD(db, &clock, student, resp[0].Ts.Format(time.RFC3339Nano))
	require.Nil(t, err)

	trans, err := getCDTransactions(db, student)
	require.Nil(t, err)

	require.Contains(t, trans[0].Description, "Early Refund")

	clock.TickOne(time.Hour * 24 * 18)

	CertificateOfDepositIfNeeded(db, &clock, student)

	err = refundCD(db, &clock, student, resp[1].Ts.Format(time.RFC3339Nano))
	require.Nil(t, err)

	trans, err = getCDTransactions(db, student)
	require.Nil(t, err)

	require.Contains(t, trans[0].Description, "Fully Matured")

}

func TestLotteryIfNeeded(t *testing.T) {

	lgr.Printf("INFO LotteryIfNeeded")
	t.Log("INFO LotteryIfNeeded")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("LotteryIfNeeded")
	defer dbTearDown()
	_, _, _, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 1)

	student, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	lottery1, err := getLottoLatest(db, student)
	require.Nil(t, err)

	LotteryIfNeeded(db, &clock, student)

	lottery2, err := getLottoLatest(db, student)
	require.Nil(t, err)

	require.Equal(t, lottery1.Jackpot, lottery2.Jackpot)

	clock.TickOne(time.Hour * 24)

	lottery3, err := getLottoLatest(db, student)
	require.Nil(t, err)

	growth := int32(decimal.NewFromInt32(lottery2.Jackpot).Mul(decimal.NewFromFloat32(KeyLottoGrowth)).IntPart())

	require.Equal(t, growth, lottery3.Jackpot)

	clock.TickOne(time.Hour * 10)

	lottery4, err := getLottoLatest(db, student)
	require.Nil(t, err)

	require.Equal(t, lottery3.Jackpot, lottery4.Jackpot)

}
