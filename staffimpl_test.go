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

func TestGetMarketPurchases(t *testing.T) {

	lgr.Printf("INFO TestGetMarketPurchases")
	t.Log("INFO TestGetMarketPurchases")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("GetMarketPurchases")
	defer dbTearDown()
	_, _, teachers, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 1)

	SetTestLoginUser(teachers[0])

	student, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	err = pay2Student(db, &clock, student, decimal.NewFromFloat(1000), teachers[0], "pre load")
	require.Nil(t, err)

	teacher, err := getUserInLocalStore(db, teachers[0])
	require.Nil(t, err)

	marketItem0 := openapi.RequestMakeMarketItem{
		Title: "item0",
		Cost:  1,
		Count: 10,
	}

	marketItem1 := openapi.RequestMakeMarketItem{
		Title: "item1",
		Cost:  1,
		Count: 1,
	}

	item0, err := makeMarketItem(db, &clock, teacher, marketItem0)
	require.Nil(t, err)

	item1, err := makeMarketItem(db, &clock, teacher, marketItem1)
	require.Nil(t, err)

	_, err = buyMarketItem(db, &clock, student, teacher, item0)
	require.Nil(t, err)
	_, err = buyMarketItem(db, &clock, student, teacher, item0)
	require.Nil(t, err)
	purchaseId, err := buyMarketItem(db, &clock, student, teacher, item1)
	require.Nil(t, err)

	resp, err := getMarketPurchases(db, teacher)
	require.Nil(t, err)

	require.Equal(t, int32(3), resp.Count)

	err = db.Update(func(tx *bolt.Tx) error {

		marketBucket, itemBucket, err := getMarketItemRx(tx, teacher, item1)
		require.Nil(t, err)
		err = marketItemResolveTx(marketBucket, itemBucket, purchaseId)
		require.Nil(t, err)
		return err

	})

	require.Nil(t, err)

	resp, err = getMarketPurchases(db, teacher)
	require.Nil(t, err)

	require.Equal(t, int32(2), resp.Count)

}

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

	require.Equal(t, 1, len(auctions))

	clock.TickOne(time.Minute * 12)

	auctions, err = getAllAuctions(db, &clock, teacher)
	require.Nil(t, err)

	require.Equal(t, 1, len(auctions))

	clock.TickOne(time.Hour * 24 * 7)

	auctions, err = getAllAuctions(db, &clock, teacher)
	require.Nil(t, err)

	require.Equal(t, 0, len(auctions))
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
		Odds:    300,
	}

	err = setSettings(db, &clock, adminDetails, settings)
	require.Nil(t, err)

	lottery, err := getLottoLatest(db, adminDetails)
	require.Nil(t, err)

	require.Equal(t, settings.Odds, lottery.Odds)
	require.Equal(t, "", lottery.Winner)

	settings2 := openapi.Settings{
		Lottery: false,
		Odds:    400,
	}

	err = setSettings(db, &clock, adminDetails, settings2)
	require.Nil(t, err)

	lottery, err = getLottoLatest(db, adminDetails)
	require.Nil(t, err)

	require.Equal(t, settings.Odds, lottery.Odds)
	require.Equal(t, "", lottery.Winner)

	settings3 := openapi.Settings{
		Lottery: true,
		Odds:    8000,
	}

	err = setSettings(db, &clock, adminDetails, settings3)
	require.Nil(t, err)

	//newest settings have odds of 8000 but that game is not over so you should still see an odds of 300 on the current game
	lottery, err = getLottoLatest(db, adminDetails)
	require.Nil(t, err)

	require.Equal(t, settings.Odds, lottery.Odds)
	require.Equal(t, "", lottery.Winner)

	student, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	err = pay2Student(db, &clock, student, decimal.NewFromFloat(10000), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	//should have winner chosen and new lottery should be created
	winner, err := purchaseLotto(db, &clock, student, 2000)
	require.Nil(t, err)
	require.True(t, winner)

	lottery, err = getLottoLatest(db, adminDetails)
	require.Nil(t, err)

	require.Equal(t, settings3.Odds, lottery.Odds)
	require.Equal(t, "", lottery.Winner)

}

func TestInitializeLotteryTooLowOdds(t *testing.T) {

	lgr.Printf("INFO TestInitializeLotteryTooLowOdds")
	t.Log("INFO TestInitializeLotteryTooLowOdds")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("InitializeLotteryTooLowOdds")
	defer dbTearDown()
	admins, _, _, _, students, err := CreateTestAccounts(db, 1, 1, 1, 1)
	require.Nil(t, err)

	adminDetails, err := getUserInLocalStore(db, admins[0])
	require.Nil(t, err)

	student, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	err = pay2Student(db, &clock, student, decimal.NewFromFloat(10000), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	settings := openapi.Settings{
		Lottery: true,
		Odds:    10,
	}

	err = setSettings(db, &clock, adminDetails, settings)
	require.Nil(t, err)

	lottery, err := getLottoLatest(db, adminDetails)
	require.Nil(t, err)

	require.Equal(t, int32(10000), lottery.Odds)
	require.Equal(t, "", lottery.Winner)

}

func TestInitializeLotteryMinimalOdds(t *testing.T) {

	lgr.Printf("INFO TestInitializeLotteryMinimalOdds")
	t.Log("INFO TestInitializeLotteryMinimalOdds")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("InitializeLotteryMinimalOdds")
	defer dbTearDown()
	admins, _, _, _, students, err := CreateTestAccounts(db, 1, 1, 1, 1)
	require.Nil(t, err)

	adminDetails, err := getUserInLocalStore(db, admins[0])
	require.Nil(t, err)

	student, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	err = pay2Student(db, &clock, student, decimal.NewFromFloat(10000), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	settings2 := openapi.Settings{
		Lottery: true,
		Odds:    10000 * .6,
	}

	err = setSettings(db, &clock, adminDetails, settings2)
	require.Nil(t, err)

	lottery, err := getLottoLatest(db, adminDetails)
	require.Nil(t, err)

	require.Equal(t, settings2.Odds, lottery.Odds)
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
		Odds:    300,
	}

	err = setSettings(db, &clock, adminDetails, settings)
	require.Nil(t, err)

	err = initializeLottery(db, adminDetails, settings, &clock)
	require.Nil(t, err)

	student, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	prevLotto, err := getLottoPrevious(db, student)
	require.Nil(t, err)
	require.Equal(t, "No Previous Raffle", prevLotto.Winner)

	err = pay2Student(db, &clock, student, decimal.NewFromFloat(10000), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	winner, err := purchaseLotto(db, &clock, student, 2000)
	require.Nil(t, err)
	require.True(t, winner)

	prevLotto, err = getLottoPrevious(db, student)
	require.Nil(t, err)
	require.Equal(t, student.Email, prevLotto.Winner)

}

func TestLotteryGarnish(t *testing.T) {

	lgr.Printf("INFO TestLotteryGarnish")
	t.Log("INFO TestLotteryGarnish")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("LotteryGarnish")
	defer dbTearDown()
	admins, _, _, _, students, err := CreateTestAccounts(db, 1, 1, 1, 1)
	require.Nil(t, err)

	adminDetails, err := getUserInLocalStore(db, admins[0])
	require.Nil(t, err)

	settings := openapi.Settings{
		Lottery: true,
		Odds:    3,
	}

	err = setSettings(db, &clock, adminDetails, settings)
	require.Nil(t, err)

	err = initializeLottery(db, adminDetails, settings, &clock)
	require.Nil(t, err)

	student, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	r := DailyPayIfNeeded(db, &clock, student)
	require.True(t, r)

	err = chargeStudent(db, &clock, student, decimal.NewFromInt(40000), CurrencyUBuck, "", false)
	require.Nil(t, err)

	err = pay2Student(db, &clock, student, decimal.NewFromFloat(30000), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	winner, err := purchaseLotto(db, &clock, student, 1000)
	require.Nil(t, err)
	require.True(t, winner)

	err = db.View(func(tx *bolt.Tx) error {
		studentBucket, err := getStudentBucketRx(tx, student.Name)
		require.Nil(t, err)

		_, _, balance, err := IsDebtNeededRx(studentBucket, &clock)
		require.Nil(t, err)
		require.Less(t, balance.InexactFloat64(), float64(40000))
		return err
	})

	require.Nil(t, err)

}

func TestDeleteStudentWithAuctionPreBid(t *testing.T) {

	lgr.Printf("INFO TestDeleteStudentWithAuctionPreBid")
	t.Log("INFO TestDeleteStudentWithAuctionPreBid")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("DeleteStudentWithAuctionPreBid")
	defer dbTearDown()
	_, _, _, classes, students, err := CreateTestAccounts(db, 1, 1, 1, 3)
	require.Nil(t, err)

	buyer, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	seller, err := getUserInLocalStore(db, students[1])
	require.Nil(t, err)

	err = pay2Student(db, &clock, buyer, decimal.NewFromFloat(100), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	request := openapi.RequestMakeAuction{
		Bid:         0,
		MaxBid:      0,
		Description: "test auc",
		EndDate:     clock.Now().Add(time.Minute),
		StartDate:   clock.Now(),
		OwnerId:     seller.Name,
		Visibility:  classes,
		TrueAuction: false,
	}

	err = MakeAuctionImpl(db, seller, request, false)
	require.Nil(t, err)

	err = deleteStudent(db, &clock, seller.Name)
	require.Nil(t, err)

	auctions, err := getStudentAuctions(db, &clock, seller)
	require.Nil(t, err)
	require.Equal(t, 0, len(auctions))

}

func TestDeleteStudentWithAuctionPostBid(t *testing.T) {
	//the student holds the auction and another student has bid on it
	//make sure that it repays the current winner and deletes the auction
	lgr.Printf("INFO TestDeleteStudentWithAuctionPostBid")
	t.Log("INFO TestDeleteStudentWithAuctionPostBid")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("DeleteStudentWithAuctionPostBid")
	defer dbTearDown()
	_, _, _, classes, students, err := CreateTestAccounts(db, 1, 1, 1, 3)
	require.Nil(t, err)

	buyer, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	seller, err := getUserInLocalStore(db, students[1])
	require.Nil(t, err)
	otherStudent, err := getUserInLocalStore(db, students[2])
	require.Nil(t, err)

	err = pay2Student(db, &clock, buyer, decimal.NewFromFloat(100), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	request := openapi.RequestMakeAuction{
		Bid:         0,
		MaxBid:      0,
		Description: "test auc",
		EndDate:     clock.Now().Add(time.Minute),
		StartDate:   clock.Now(),
		OwnerId:     otherStudent.Name,
		Visibility:  classes,
		TrueAuction: false,
	}

	for i := 0; i < 5; i++ {
		err = MakeAuctionImpl(db, otherStudent, request, false)
		require.Nil(t, err)
		clock.TickOne(time.Second)
		request.EndDate = clock.Now().Add(time.Minute)
	}

	clock.TickOne(time.Second)
	request.OwnerId = seller.Name
	request.EndDate = clock.Now().Add(time.Minute)

	err = MakeAuctionImpl(db, seller, request, false)
	require.Nil(t, err)

	auctions, err := getStudentAuctions(db, &clock, seller)
	require.Nil(t, err)
	require.Equal(t, 6, len(auctions))

	_, err = placeBid(db, &clock, buyer, auctions[0].Id.Format(time.RFC3339Nano), 10, &NoopSSEService{})
	require.Nil(t, err)

	ubucks, err := getStudentUbuck(db, buyer)
	require.Nil(t, err)
	require.Greater(t, float32(100), ubucks.Value)

	err = deleteStudent(db, &clock, seller.Name)
	require.Nil(t, err)

	auctions, err = getStudentAuctions(db, &clock, buyer)
	require.Nil(t, err)
	require.Equal(t, 5, len(auctions))

	ubucks, err = getStudentUbuck(db, buyer)
	require.Nil(t, err)
	require.Equal(t, float32(100), ubucks.Value)

}

func TestDeleteStudentWithAuctionPostBidderDelete(t *testing.T) {
	//someone has created an auction and student A has bid on it
	//student A is deleted and the auction is deleted
	lgr.Printf("INFO TestDeleteStudentWithAuctionPostBidderDelete")
	t.Log("INFO TestDeleteStudentWithAuctionPostBidderDelete")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("DeleteStudentWithAuctionPostBidderDelete")
	defer dbTearDown()
	_, _, _, classes, students, err := CreateTestAccounts(db, 1, 1, 1, 3)
	require.Nil(t, err)

	buyer, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	seller, err := getUserInLocalStore(db, students[1])
	require.Nil(t, err)
	otherStudent, err := getUserInLocalStore(db, students[2])
	require.Nil(t, err)

	err = pay2Student(db, &clock, buyer, decimal.NewFromFloat(100), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	request := openapi.RequestMakeAuction{
		Bid:         0,
		MaxBid:      0,
		Description: "test auc",
		EndDate:     clock.Now().Add(time.Minute),
		StartDate:   clock.Now(),
		OwnerId:     otherStudent.Name,
		Visibility:  classes,
		TrueAuction: false,
	}

	for i := 0; i < 5; i++ {
		err = MakeAuctionImpl(db, otherStudent, request, false)
		require.Nil(t, err)
		clock.TickOne(time.Second)
		request.EndDate = clock.Now().Add(time.Minute)
	}

	clock.TickOne(time.Second)
	request.OwnerId = seller.Name
	request.EndDate = clock.Now().Add(time.Minute)

	err = MakeAuctionImpl(db, seller, request, false)
	require.Nil(t, err)

	auctions, err := getStudentAuctions(db, &clock, seller)
	require.Nil(t, err)
	require.Equal(t, 6, len(auctions)) // 5 auctions + 1 auction

	_, err = placeBid(db, &clock, buyer, auctions[0].Id.Format(time.RFC3339Nano), 10, &NoopSSEService{})
	require.Nil(t, err)

	err = deleteStudent(db, &clock, buyer.Name)
	require.Nil(t, err)

	resp, err := getStudentAuctions(db, &clock, seller)
	require.Nil(t, err)
	//if the lead bidder is deleted then just delete the auction
	require.Equal(t, 5, len(resp))

}
