package main

import (
	"encoding/json"
	"testing"
	"time"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/go-pkgz/lgr"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"
)

func TestAuctionsAll(t *testing.T) {

	lgr.Printf("INFO TestAuctionsAll")
	t.Log("INFO TestAuctionsAll")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("AuctionsAll")
	defer dbTearDown()
	_, _, teachers, classes, students, _ := CreateTestAccounts(db, 1, 1, 1, 1)

	student, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	teacher, err := getUserInLocalStore(db, teachers[0])
	require.Nil(t, err)

	body := openapi.RequestMakeAuction{
		Bid:         4,
		MaxBid:      4,
		Description: "Test Auction",
		EndDate:     clock.Now().Add(time.Minute * 100),
		StartDate:   clock.Now().Add(time.Minute * -10),
		OwnerId:     teachers[0],
		Visibility:  classes,
	}

	err = MakeAuctionImpl(db, teacher, body, true)
	require.Nil(t, err)

	body = openapi.RequestMakeAuction{
		Bid:         4,
		MaxBid:      4,
		Description: "Test Auction",
		EndDate:     clock.Now().Add(time.Minute * 10),
		StartDate:   clock.Now().Add(time.Minute * -10),
		OwnerId:     students[0],
		Visibility:  classes,
	}

	err = MakeAuctionImpl(db, student, body, false)
	require.Nil(t, err)

	auctions, err := getAllAuctions(db, &clock, teacher)
	require.Nil(t, err)

	require.Equal(t, 2, len(auctions))

	clock.TickOne(time.Minute * 12)

	auctions, err = getAllAuctions(db, &clock, teacher)
	require.Nil(t, err)

	require.Equal(t, 1, len(auctions))
}

func TestApproveAuction(t *testing.T) {

	lgr.Printf("INFO TestApproveAuction")
	t.Log("INFO TestApproveAuction")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("approveAuction")
	defer dbTearDown()
	_, _, teachers, classes, students, _ := CreateTestAccounts(db, 1, 1, 1, 1)

	student, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	teacher, err := getUserInLocalStore(db, teachers[0])
	require.Nil(t, err)

	body := openapi.RequestMakeAuction{
		Bid:         4,
		MaxBid:      4,
		Description: "Test Auction",
		EndDate:     clock.Now().Add(time.Minute * 100),
		StartDate:   clock.Now().Add(time.Minute * -10),
		OwnerId:     students[0],
		Visibility:  classes,
	}

	err = MakeAuctionImpl(db, student, body, false)
	require.Nil(t, err)

	auctions, err := getAllAuctions(db, &clock, teacher)
	require.Nil(t, err)

	require.False(t, auctions[0].Approved)

	actionBody := openapi.RequestAuctionAction{
		AuctionId: auctions[0].Id.Format(time.RFC3339Nano),
	}

	err = approveAuction(db, teacher, actionBody)
	require.Nil(t, err)

	auctions, err = getAllAuctions(db, &clock, teacher)
	require.Nil(t, err)

	require.True(t, auctions[0].Approved)

}

func TestRejectAuction(t *testing.T) {

	lgr.Printf("INFO TestRejectAuction")
	t.Log("INFO TestRejectAuction")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("rejectAuction")
	defer dbTearDown()
	_, _, teachers, classes, students, _ := CreateTestAccounts(db, 1, 1, 1, 1)

	student, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	teacher, err := getUserInLocalStore(db, teachers[0])
	require.Nil(t, err)

	body := openapi.RequestMakeAuction{
		Bid:         4,
		MaxBid:      4,
		Description: "Test Auction",
		EndDate:     clock.Now().Add(time.Minute * 100),
		StartDate:   clock.Now().Add(time.Minute * -10),
		OwnerId:     students[0],
		Visibility:  classes,
	}

	err = MakeAuctionImpl(db, student, body, false)
	require.Nil(t, err)

	body = openapi.RequestMakeAuction{
		Bid:         4,
		MaxBid:      4,
		Description: "Test Auction",
		EndDate:     clock.Now().Add(time.Minute * 10),
		StartDate:   clock.Now().Add(time.Minute * -10),
		OwnerId:     students[0],
		Visibility:  classes,
	}

	err = MakeAuctionImpl(db, student, body, false)
	require.Nil(t, err)

	auctions, err := getAllAuctions(db, &clock, teacher)
	require.Nil(t, err)

	auctionId := auctions[0].Id.Format(time.RFC3339Nano)

	err = rejectAuction(db, teacher, auctionId)
	require.Nil(t, err)

	auctions, err = getAllAuctions(db, &clock, teacher)
	require.Nil(t, err)

	require.Equal(t, 1, len(auctions))

}

func TestMakeMarketItemImpl(t *testing.T) {

	lgr.Printf("INFO TestMakeMarketItemImpl")
	t.Log("INFO TestMakeMarketItemImpl")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("MakeMarketItemImpl")
	defer dbTearDown()
	_, _, teachers, _, _, _ := CreateTestAccounts(db, 1, 1, 1, 1)

	teacher, err := getUserInLocalStore(db, teachers[0])
	require.Nil(t, err)

	body := openapi.RequestMakeMarketItem{
		Title: "Candy",
		Count: 54,
		Cost:  23,
	}

	id, err := makeMarketItem(db, &clock, teacher, body)
	require.Nil(t, err)

	items, err := getMarketItems(db, teacher)
	require.Nil(t, err)

	require.Equal(t, 1, len(items))
	require.Equal(t, "Candy", items[0].Title)

	_ = db.View(func(tx *bolt.Tx) error {
		_, itemBucket, err := getMarketItemRx(tx, teacher, id)
		require.Nil(t, err)

		var details MarketItem

		itemData := itemBucket.Get([]byte(KeyMarketData))
		err = json.Unmarshal(itemData, &details)
		require.Nil(t, err)

		item, err := packageMarketItemRx(tx, details, teacher, itemBucket, id)
		require.Nil(t, err)

		require.Equal(t, "Candy", item.Title)

		return nil
	})

}

func TestMakeMarketItemImplBuyers(t *testing.T) {

	lgr.Printf("INFO TestMakeMarketItemImplBuyers")
	t.Log("INFO TestMakeMarketItemImplBuyers")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("MakeMarketItemImplBuyers")
	defer dbTearDown()
	_, _, teachers, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 4)

	teacher, err := getUserInLocalStore(db, teachers[0])
	require.Nil(t, err)

	student0, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	err = pay2Student(db, &clock, student0, decimal.NewFromFloat(100), teachers[0], "pre load")
	require.Nil(t, err)

	student1, err := getUserInLocalStore(db, students[1])
	require.Nil(t, err)

	err = pay2Student(db, &clock, student1, decimal.NewFromFloat(100), teachers[0], "pre load")
	require.Nil(t, err)

	student2, err := getUserInLocalStore(db, students[2])
	require.Nil(t, err)

	err = pay2Student(db, &clock, student2, decimal.NewFromFloat(100), teachers[0], "pre load")
	require.Nil(t, err)

	student3, err := getUserInLocalStore(db, students[3])
	require.Nil(t, err)

	err = pay2Student(db, &clock, student3, decimal.NewFromFloat(100), teachers[0], "pre load")
	require.Nil(t, err)

	body := openapi.RequestMakeMarketItem{
		Title: "Candy",
		Count: 3,
		Cost:  4,
	}

	id, err := makeMarketItem(db, &clock, teacher, body)
	require.Nil(t, err)

	_, err = buyMarketItem(db, &clock, student0, teacher, id)
	require.Nil(t, err)

	_, err = buyMarketItem(db, &clock, student1, teacher, id)
	require.Nil(t, err)

	_, err = buyMarketItem(db, &clock, student2, teacher, id)
	require.Nil(t, err)

	_, err = buyMarketItem(db, &clock, student2, teacher, id)
	require.NotNil(t, err)
	require.Equal(t, "ERROR there is nothing left to buy", err.Error())

	_ = db.View(func(tx *bolt.Tx) error {
		_, itemBucket, err := getMarketItemRx(tx, teacher, id)
		require.Nil(t, err)

		var details MarketItem

		itemData := itemBucket.Get([]byte(KeyMarketData))
		err = json.Unmarshal(itemData, &details)
		require.Nil(t, err)

		item, err := packageMarketItemRx(tx, details, teacher, itemBucket, id)
		require.Nil(t, err)

		require.Equal(t, 3, len(item.Buyers))

		return nil
	})

}

func TestMakeMarketItemImplBuyers1multi(t *testing.T) {

	lgr.Printf("INFO TestMakeMarketItemImplBuyers1Multi")
	t.Log("INFO TestMakeMarketItemImplBuyers1Multi")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("MakeMarketItemImplBuyers1Multi")
	defer dbTearDown()
	_, _, teachers, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 4)

	teacher, err := getUserInLocalStore(db, teachers[0])
	require.Nil(t, err)

	student0, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	err = pay2Student(db, &clock, student0, decimal.NewFromFloat(8), teachers[0], "pre load")
	require.Nil(t, err)

	body := openapi.RequestMakeMarketItem{
		Title: "Candy",
		Count: 3,
		Cost:  4,
	}

	id, err := makeMarketItem(db, &clock, teacher, body)
	require.Nil(t, err)

	_, err = buyMarketItem(db, &clock, student0, teacher, id)
	require.Nil(t, err)

	_, err = buyMarketItem(db, &clock, student0, teacher, id)
	require.Nil(t, err)

	_, err = buyMarketItem(db, &clock, student0, teacher, id)
	require.NotNil(t, err)
	require.Equal(t, "insufficient funds", err.Error())

	_ = db.View(func(tx *bolt.Tx) error {
		_, itemBucket, err := getMarketItemRx(tx, teacher, id)
		require.Nil(t, err)

		var details MarketItem

		itemData := itemBucket.Get([]byte(KeyMarketData))
		err = json.Unmarshal(itemData, &details)
		require.Nil(t, err)

		item, err := packageMarketItemRx(tx, details, teacher, itemBucket, id)
		require.Nil(t, err)

		require.Equal(t, 2, len(item.Buyers))

		return nil
	})

}

func TestMarketResolveTx(t *testing.T) {

	lgr.Printf("INFO TestMarketResolveTx")
	t.Log("INFO TestMarketResolveTx")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("marketResolveTx")
	defer dbTearDown()
	_, _, teachers, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 4)

	teacher, err := getUserInLocalStore(db, teachers[0])
	require.Nil(t, err)

	student0, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	err = pay2Student(db, &clock, student0, decimal.NewFromFloat(8), teachers[0], "pre load")
	require.Nil(t, err)

	body := openapi.RequestMakeMarketItem{
		Title: "Candy",
		Count: 1,
		Cost:  4,
	}

	id, err := makeMarketItem(db, &clock, teacher, body)
	require.Nil(t, err)

	purchaseId, err := buyMarketItem(db, &clock, student0, teacher, id)
	require.Nil(t, err)

	_ = db.Update(func(tx *bolt.Tx) error {
		marketBucket, itemBucket, err := getMarketItemRx(tx, teacher, id)
		if err != nil {
			return err
		}
		err = marketItemResolveTx(marketBucket, itemBucket, purchaseId)
		require.Nil(t, err)

		_, itemBucket, err = getMarketItemRx(tx, teacher, id)
		require.Nil(t, err)

		var details MarketItem

		itemData := itemBucket.Get([]byte(KeyMarketData))
		err = json.Unmarshal(itemData, &details)
		require.Nil(t, err)

		item, err := packageMarketItemRx(tx, details, teacher, itemBucket, id)
		require.Nil(t, err)

		require.False(t, details.Active)

		require.Equal(t, 0, len(item.Buyers))

		return nil
	})

}

func TestMarketItemRefundTx(t *testing.T) {

	lgr.Printf("INFO TestMarketItemRefundTx")
	t.Log("INFO TestMarketItemRefundTx")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("marketItemRefundTx")
	defer dbTearDown()
	_, _, teachers, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 4)

	teacher, err := getUserInLocalStore(db, teachers[0])
	require.Nil(t, err)

	student0, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	err = pay2Student(db, &clock, student0, decimal.NewFromFloat(8), teachers[0], "pre load")
	require.Nil(t, err)

	body := openapi.RequestMakeMarketItem{
		Title: "Candy",
		Count: 3,
		Cost:  4,
	}

	id, err := makeMarketItem(db, &clock, teacher, body)
	require.Nil(t, err)

	purchaseId, err := buyMarketItem(db, &clock, student0, teacher, id)
	require.Nil(t, err)

	_ = db.Update(func(tx *bolt.Tx) error {
		resp, err := getStudentBuckRx(tx, student0, teacher.Email)
		require.Nil(t, err)
		require.Equal(t, float32(4), resp.Value)

		_, itemBucket, err := getMarketItemRx(tx, teacher, id)
		require.Nil(t, err)

		err = marketItemRefundTx(tx, &clock, itemBucket, purchaseId, teacher.Email)
		require.Nil(t, err)

		resp, err = getStudentBuckRx(tx, student0, teacher.Email)
		require.Nil(t, err)
		require.Equal(t, float32(8), resp.Value)

		_, itemBucket, err = getMarketItemRx(tx, teacher, id)
		require.Nil(t, err)

		var details MarketItem

		itemData := itemBucket.Get([]byte(KeyMarketData))
		err = json.Unmarshal(itemData, &details)
		require.Nil(t, err)

		item, err := packageMarketItemRx(tx, details, teacher, itemBucket, id)
		require.Nil(t, err)

		require.Equal(t, 0, len(item.Buyers))
		require.Equal(t, int32(3), item.Count)

		return nil
	})

}

func TestMarketItemDeleteTx(t *testing.T) {

	lgr.Printf("INFO TestMarketItemDeleteTx")
	t.Log("INFO TestMarketItemDeleteTx")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("marketItemDeleteTx")
	defer dbTearDown()
	_, _, teachers, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 4)

	teacher, err := getUserInLocalStore(db, teachers[0])
	require.Nil(t, err)

	student0, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	student1, err := getUserInLocalStore(db, students[1])
	require.Nil(t, err)

	student2, err := getUserInLocalStore(db, students[2])
	require.Nil(t, err)

	student3, err := getUserInLocalStore(db, students[3])
	require.Nil(t, err)

	err = pay2Student(db, &clock, student0, decimal.NewFromFloat(8), teachers[0], "pre load")
	require.Nil(t, err)

	err = pay2Student(db, &clock, student1, decimal.NewFromFloat(8), teachers[0], "pre load")
	require.Nil(t, err)

	err = pay2Student(db, &clock, student2, decimal.NewFromFloat(8), teachers[0], "pre load")
	require.Nil(t, err)

	err = pay2Student(db, &clock, student3, decimal.NewFromFloat(8), teachers[0], "pre load")
	require.Nil(t, err)

	body := openapi.RequestMakeMarketItem{
		Title: "Candy",
		Count: 8,
		Cost:  1,
	}

	id, err := makeMarketItem(db, &clock, teacher, body)
	require.Nil(t, err)

	_, err = buyMarketItem(db, &clock, student0, teacher, id)
	require.Nil(t, err)

	_, err = buyMarketItem(db, &clock, student1, teacher, id)
	require.Nil(t, err)

	_, err = buyMarketItem(db, &clock, student2, teacher, id)
	require.Nil(t, err)

	_, err = buyMarketItem(db, &clock, student3, teacher, id)
	require.Nil(t, err)

	purchaseIdResolve, err := buyMarketItem(db, &clock, student0, teacher, id)
	require.Nil(t, err)

	purchaseIdRefund, err := buyMarketItem(db, &clock, student0, teacher, id)
	require.Nil(t, err)

	_, err = buyMarketItem(db, &clock, student0, teacher, id)
	require.Nil(t, err)

	_, err = buyMarketItem(db, &clock, student0, teacher, id)
	require.Nil(t, err)

	_ = db.Update(func(tx *bolt.Tx) error {
		resp, err := getStudentBuckRx(tx, student0, teacher.Email)
		require.Nil(t, err)
		require.Equal(t, float32(3), resp.Value)

		marketBucket, itemBucket, err := getMarketItemRx(tx, teacher, id)
		require.Nil(t, err)

		err = marketItemResolveTx(marketBucket, itemBucket, purchaseIdResolve)
		require.Nil(t, err)

		err = marketItemRefundTx(tx, &clock, itemBucket, purchaseIdRefund, teacher.Email)
		require.Nil(t, err)

		err = marketItemDeleteTx(tx, &clock, marketBucket, itemBucket, id, teacher.Email)
		require.Nil(t, err)

		resp, err = getStudentBuckRx(tx, student0, teacher.Email)
		require.Nil(t, err)
		require.Equal(t, float32(7), resp.Value)

		return nil
	})

	_, _, err = getMarketItem(db, teacher, id)
	require.Nil(t, err)

}

func TestInitializeLottery(t *testing.T) {

	lgr.Printf("INFO TestInitializeLottery")
	t.Log("INFO TestInitializeLottery")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("InitializeLottery")
	defer dbTearDown()
	admins, _, _, _, students, err := CreateTestAccounts(db, 1, 1, 1, 1)
	require.Nil(t, err)

	adminDetails, err := getUserInLocalStore(db, admins[0])
	require.Nil(t, err)

	settings := openapi.Settings{
		Lottery: true,
		Odds:    10,
	}

	err = setSettings(db, &clock, adminDetails, settings)
	require.Nil(t, err)

	lottery, err := getLottoLatest(db, adminDetails)
	require.Nil(t, err)

	require.Equal(t, settings.Odds, lottery.Odds)
	require.Equal(t, "", lottery.Winner)

	settings2 := openapi.Settings{
		Lottery: false,
		Odds:    10,
	}

	err = setSettings(db, &clock, adminDetails, settings2)
	require.Nil(t, err)

	lottery, err = getLottoLatest(db, adminDetails)
	require.Nil(t, err)

	require.Equal(t, settings.Odds, lottery.Odds)
	require.Equal(t, "", lottery.Winner)

	settings3 := openapi.Settings{
		Lottery: true,
		Odds:    20,
	}

	err = setSettings(db, &clock, adminDetails, settings3)
	require.Nil(t, err)

	//newest settings have odds of 20 but that game is not over so you should still see an odds of 10 on the current game
	lottery, err = getLottoLatest(db, adminDetails)
	require.Nil(t, err)

	require.Equal(t, settings.Odds, lottery.Odds)
	require.Equal(t, "", lottery.Winner)

	student, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	err = pay2Student(db, &clock, student, decimal.NewFromFloat(800), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	//should have winner chosen and new lottery should be created
	winner, err := purchaseLotto(db, &clock, student, 300)
	require.Nil(t, err)
	require.True(t, winner)

	lottery, err = getLottoLatest(db, adminDetails)
	require.Nil(t, err)

	require.Equal(t, settings3.Odds, lottery.Odds)
	require.Equal(t, "", lottery.Winner)

}

func TestLotteryProgression(t *testing.T) {

	lgr.Printf("INFO TestLotteryProgression")
	t.Log("INFO TestLotteryProgression")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("LotteryProgression")
	defer dbTearDown()
	admins, _, _, _, students, err := CreateTestAccounts(db, 1, 1, 1, 1)
	require.Nil(t, err)

	adminDetails, err := getUserInLocalStore(db, admins[0])
	require.Nil(t, err)

	settings := openapi.Settings{
		Lottery: true,
		Odds:    10000,
	}

	err = setSettings(db, &clock, adminDetails, settings)
	require.Nil(t, err)

	student, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	err = pay2Student(db, &clock, student, decimal.NewFromFloat(1000), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	//should have winner chosen and new lottery should be created
	winner, err := purchaseLotto(db, &clock, student, 10)
	require.Nil(t, err)
	require.False(t, winner)

	lottery, err := getLottoLatest(db, adminDetails)
	require.Nil(t, err)

	mean, err := getMeanNetworth(db, student)
	require.Nil(t, err)
	jp := mean.IntPart() + 10

	require.Equal(t, int32(jp), lottery.Jackpot)

}

func TestLotteryLastWinner(t *testing.T) {

	lgr.Printf("INFO TestLotteryLastWinner")
	t.Log("INFO TestLotteryLastWinner")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("LotteryLastWinner")
	defer dbTearDown()
	admins, _, _, _, students, err := CreateTestAccounts(db, 1, 1, 1, 1)
	require.Nil(t, err)

	adminDetails, err := getUserInLocalStore(db, admins[0])
	require.Nil(t, err)

	settings := openapi.Settings{
		Lottery: true,
		Odds:    10,
	}

	err = setSettings(db, &clock, adminDetails, settings)
	require.Nil(t, err)

	err = initializeLottery(db, adminDetails, settings, &clock)
	require.Nil(t, err)

	student, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	prevLotto, err := getLottoPrevious(db, student)
	require.Nil(t, err)
	require.Equal(t, "No Previous Lotto", prevLotto.Winner)

	err = pay2Student(db, &clock, student, decimal.NewFromFloat(1000), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	_, _, err = getSchoolStudents(db, student)
	require.Nil(t, err)

	winner, err := purchaseLotto(db, &clock, student, 100)
	require.Nil(t, err)
	require.True(t, winner)

	_, _, err = getSchoolStudents(db, student)
	require.Nil(t, err)

	prevLotto, err = getLottoPrevious(db, student)
	require.Nil(t, err)
	require.Equal(t, student.Email, prevLotto.Winner)

}
