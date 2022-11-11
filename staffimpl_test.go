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
		itemBucket, err := getMarketItemRx(tx, teacher, id)
		require.Nil(t, err)

		item, err := packageMarketItemRx(tx, teacher, itemBucket)
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
		itemBucket, err := getMarketItemRx(tx, teacher, id)
		require.Nil(t, err)

		item, err := packageMarketItemRx(tx, teacher, itemBucket)
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
	require.Equal(t, "Insufficient funds", err.Error())

	_ = db.View(func(tx *bolt.Tx) error {
		itemBucket, err := getMarketItemRx(tx, teacher, id)
		require.Nil(t, err)

		item, err := packageMarketItemRx(tx, teacher, itemBucket)
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
		Count: 3,
		Cost:  4,
	}

	id, err := makeMarketItem(db, &clock, teacher, body)
	require.Nil(t, err)

	purchaseId, err := buyMarketItem(db, &clock, student0, teacher, id)
	require.Nil(t, err)

	_ = db.Update(func(tx *bolt.Tx) error {
		itemBucket, err := getMarketItemRx(tx, teacher, id)
		if err != nil {
			return err
		}
		err = marketItemResolveTx(itemBucket, purchaseId)
		require.Nil(t, err)

		itemBucket, err = getMarketItemRx(tx, teacher, id)
		require.Nil(t, err)

		item, err := packageMarketItemRx(tx, teacher, itemBucket)
		require.Nil(t, err)

		require.Equal(t, 0, len(item.Buyers))

		return nil
	})

}
