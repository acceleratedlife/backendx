package main

import (
	"testing"
	"time"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/go-pkgz/lgr"
	"github.com/stretchr/testify/require"
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
