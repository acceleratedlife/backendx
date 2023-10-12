package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/go-pkgz/auth/token"
	"github.com/go-pkgz/lgr"
	"github.com/shopspring/decimal"
	bolt "go.etcd.io/bbolt"
)

type StudentApiServiceImpl struct {
	db    *bolt.DB
	clock Clock
}

func (a *StudentApiServiceImpl) BuyCD(ctx context.Context, CD_details openapi.RequestBuyCd) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	if userDetails.Role != UserRoleStudent {
		return openapi.Response(401, ""), nil
	}

	err = buyCD(a.db, a.clock, userDetails, CD_Details)

	if err != nil {
		return openapi.Response(400, nil), err
	}

	return openapi.Response(200, nil), nil

}

func (a *StudentApiServiceImpl) RefundCD(ctx context.Context, CD_id openapi.RequestUser) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	if userDetails.Role != UserRoleStudent {
		return openapi.Response(401, ""), nil
	}

	// resp, err := getStudentBuck(a.db, userDetails, CD_Details)

	if err != nil {
		return openapi.Response(400, nil), err
	}

	return openapi.Response(200, nil), nil

}

func (a *StudentApiServiceImpl) SearchCDS(ctx context.Context) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	if userDetails.Role != UserRoleStudent {
		return openapi.Response(401, ""), nil
	}

	// resp, err := getStudentBuck(a.db, userDetails, CD_Details)

	if err != nil {
		return openapi.Response(400, nil), err
	}

	return openapi.Response(200, nil), nil

}

func (a *StudentApiServiceImpl) SearchCDTransactions(ctx context.Context) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	if userDetails.Role != UserRoleStudent {
		return openapi.Response(401, ""), nil
	}

	// resp, err := getStudentBuck(a.db, userDetails, CD_Details)

	if err != nil {
		return openapi.Response(400, nil), err
	}

	return openapi.Response(200, nil), nil

}

func (a *StudentApiServiceImpl) SearchBuck(ctx context.Context, teacherId string) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	if userDetails.Role != UserRoleStudent {
		return openapi.Response(401, ""), nil
	}

	resp, err := getStudentBuck(a.db, userDetails, teacherId)

	if err != nil {
		return openapi.Response(400, nil), err
	}

	return openapi.Response(200, resp), nil

}

func (a *StudentApiServiceImpl) LatestLotto(ctx context.Context) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	lottery, err := getLottoLatest(a.db, userDetails)
	if err != nil {
		return openapi.Response(400, nil), err
	}

	resp := openapi.ResponseLottoLatest{
		Jackpot: lottery.Jackpot,
		Odds:    lottery.Odds,
	}

	return openapi.Response(200, resp), nil

}

func (a *StudentApiServiceImpl) LottoPurchase(ctx context.Context, tickets int32) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	winner, err := purchaseLotto(a.db, a.clock, userDetails, tickets)
	if err != nil {
		return openapi.Response(400, nil), err
	}

	return openapi.Response(200, winner), nil

}

func (a *StudentApiServiceImpl) PreviousLotto(ctx context.Context) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	lotteryPrev, err := getLottoPrevious(a.db, userDetails)
	if err != nil {
		return openapi.Response(400, nil), err
	}

	resp := openapi.ResponseLottoLatest{
		Jackpot: lotteryPrev.Jackpot,
		Winner:  lotteryPrev.Winner,
	}

	return openapi.Response(200, resp), nil

}

func (a *StudentApiServiceImpl) MarketItemBuy(ctx context.Context, body openapi.RequestMarketRefund) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	if userDetails.Role != UserRoleStudent {
		return openapi.Response(401, ""), nil
	}

	teacher, err := getUserInLocalStore(a.db, body.TeacherId)
	if err != nil {
		return openapi.Response(400, nil), err
	}

	_, err = buyMarketItem(a.db, a.clock, userDetails, teacher, body.Id)
	if err != nil {
		return openapi.Response(400, nil), err
	}

	return openapi.Response(200, nil), nil

}

func (a *StudentApiServiceImpl) AuctionBid(ctx context.Context, body openapi.RequestAuctionBid) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}
	if userDetails.Role != UserRoleStudent {
		return openapi.Response(401, ""), nil
	}

	var message string
	err = a.db.Update(func(tx *bolt.Tx) error {
		message, err = placeBidtx(tx, a.clock, userDetails, body.Item, int32(body.Bid))
		if err != nil {
			return err
		}

		return nil

	})

	if err != nil {
		return openapi.Response(400, nil), err
	}

	return openapi.Response(200, nil), fmt.Errorf(message)
}

func (a *StudentApiServiceImpl) BuckConvert(ctx context.Context, body openapi.RequestBuckConvert) (response openapi.ImplResponse, err error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}
	if userDetails.Role != UserRoleStudent {
		return openapi.Response(401, ""), nil
	}

	if body.AccountFrom == body.AccountTo {
		return openapi.Response(400, nil), fmt.Errorf("can't convert same bucks")
	}

	if body.AccountFrom == KeyDebt {
		return openapi.Response(400, nil), fmt.Errorf("can't convert from debt account")
	}

	err = a.db.Update(func(tx *bolt.Tx) error {
		err := studentConvertTx(tx, a.clock, userDetails, decimal.NewFromFloat32(body.Amount), body.AccountFrom, body.AccountTo, "", true)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return openapi.Response(400, nil), err
	}

	return openapi.Response(200, nil), nil

}
func (a *StudentApiServiceImpl) CryptoConvert(ctx context.Context, body openapi.RequestCryptoConvert) (resp openapi.ImplResponse, err error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}
	if userDetails.Role != UserRoleStudent {
		return openapi.Response(401, ""), nil
	}

	err = cryptoTransaction(a.db, a.clock, userDetails, body)

	if err != nil {
		return openapi.Response(400, nil), err
	}

	return openapi.Response(200, nil), nil
}
func (a *StudentApiServiceImpl) SearchAuctionsStudent(ctx context.Context) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}
	if userDetails.Role != UserRoleStudent {
		return openapi.Response(401, ""), nil
	}

	var resp []openapi.ResponseAuctionStudent
	auctions, err := getStudentAuctions(a.db, userDetails)
	if err != nil {
		return openapi.Response(400, err), nil
	}

	for _, auction := range auctions {

		now := a.clock.Now()
		if (auction.Approved && auction.StartDate.Before(now) && auction.EndDate.After(now) && auction.OwnerId.Id != userDetails.Name) || auction.WinnerId.Id == userDetails.Name {
			iAuction := openapi.ResponseAuctionStudent{
				Id:          auction.Id,
				Active:      auction.Active,
				Bid:         float32(auction.Bid),
				Description: auction.Description,
				EndDate:     auction.EndDate,
				StartDate:   auction.StartDate,
				OwnerId: openapi.ResponseAuctionStudentOwnerId{
					Id:       auction.OwnerId.Id,
					LastName: auction.OwnerId.LastName,
				},
				WinnerId: openapi.ResponseAuctionStudentOwnerId{
					Id:        auction.WinnerId.Id,
					FirstName: auction.WinnerId.FirstName,
					LastName:  auction.WinnerId.LastName,
				},
			}

			resp = append(resp, iAuction)
		}
	}

	if err != nil {
		return openapi.Response(400, err), nil
	}

	return openapi.Response(200, resp), nil
}
func (a *StudentApiServiceImpl) SearchBuckTransactions(ctx context.Context) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}
	if userDetails.Role != UserRoleStudent {
		return openapi.Response(401, ""), nil
	}

	var resp []openapi.ResponseBuckTransaction
	err = a.db.View(func(tx *bolt.Tx) error {

		student, err := getStudentBucketRx(tx, userDetails.Name)
		if err != nil {
			return err
		}

		bAccounts := student.Bucket([]byte(KeyAccounts))
		if bAccounts == nil {
			return fmt.Errorf("failed to get bAccount")
		}

		resp, err = getStudentTransactionsTx(tx, bAccounts, userDetails)
		if err != nil {
			return err
		}

		return nil

	})

	if err != nil {
		return openapi.Response(400, nil), err
	}

	return openapi.Response(200, resp), nil

}
func (a *StudentApiServiceImpl) SearchCrypto(ctx context.Context, crypto string) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}
	if userDetails.Role != UserRoleStudent {
		return openapi.Response(401, ""), nil
	}

	var resp openapi.ResponseCrypto
	resp, err = getCryptoForStudentRequest(a.db, userDetails, crypto)

	if err != nil {
		lgr.Printf("ERROR cannot get Cryptos for: %s %v", userDetails.Name, err)
		return openapi.Response(500, "{}"), err
	}
	return openapi.Response(200, resp), nil

}
func (a *StudentApiServiceImpl) SearchCryptoTransaction(ctx context.Context) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}
	if userDetails.Role != UserRoleStudent {
		return openapi.Response(401, ""), nil
	}

	var resp []openapi.ResponseCryptoTransaction
	err = a.db.View(func(tx *bolt.Tx) error {
		resp, err = getStudentCryptoTransactionsRx(tx, userDetails)
		return err
	})

	if err != nil {
		lgr.Printf("ERROR cannot get transactions for: %s %v", userDetails.Name, err)
		return openapi.Response(500, "{}"), err
	}

	return openapi.Response(200, resp), nil
}
func (a *StudentApiServiceImpl) SearchStudentCrypto(ctx context.Context) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}
	if userDetails.Role != UserRoleStudent {
		return openapi.Response(401, ""), nil
	}

	var resp []CryptoDecimal
	err = a.db.View(func(tx *bolt.Tx) error {
		resp, err = getStudentCryptosRx(tx, userDetails)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		lgr.Printf("ERROR cannot get Cryptos for: %s %v", userDetails.Name, err)
		return openapi.Response(500, "{}"), err
	}
	return openapi.Response(200, resp), nil

}
func (a *StudentApiServiceImpl) SearchStudentUbuck(ctx context.Context) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}
	if userDetails.Role != UserRoleStudent {
		return openapi.Response(401, ""), nil
	}

	ubucks, err := getStudentUbuck(a.db, userDetails)
	if err != nil {
		lgr.Printf("ERROR cannot get Ubucks for: %s %v", userDetails.Name, err)
		return openapi.Response(500, "{}"), nil
	}
	return openapi.Response(200, ubucks), nil

}

func (a *StudentApiServiceImpl) StudentAddClass(ctx context.Context, body openapi.RequestAddClass) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}
	if userDetails.Role != UserRoleStudent {
		return openapi.Response(401, ""), nil
	}

	_, pathId, err := RoleByAddCode(a.db, body.AddCode, a.clock)
	if err != nil {
		return openapi.Response(404,
			openapi.ResponseRegister4{
				Message: err.Error(),
			}), nil
	}

	err = a.db.Update(func(tx *bolt.Tx) error {
		schools := tx.Bucket([]byte(KeySchools))
		school := schools.Bucket([]byte(userDetails.SchoolId))
		if school == nil {
			return fmt.Errorf("can't find school")
		}

		schoolClasses := school.Bucket([]byte(KeyClasses))
		if schoolClasses == nil {
			return fmt.Errorf("can't find school classes")
		}
		class := schoolClasses.Bucket([]byte(pathId.classId))
		if class != nil {
			studentsBucket, err := class.CreateBucketIfNotExists([]byte(KeyStudents))
			if err != nil {
				return err
			}
			if studentsBucket != nil {
				err = studentsBucket.Put([]byte(userDetails.Email), nil)
				if err != nil {
					return err
				}
				return nil
			}
		}

		teachers := school.Bucket([]byte(KeyTeachers))
		if teachers == nil {
			return fmt.Errorf("can't find teachers")
		}
		teacher := teachers.Bucket([]byte(pathId.teacherId))
		if teacher == nil {
			return fmt.Errorf("can't find teacher")
		}
		classesBucket := teacher.Bucket([]byte(KeyClasses))
		if classesBucket == nil {
			return fmt.Errorf("can't find classesBucket")
		}
		classBucket := classesBucket.Bucket([]byte(pathId.classId))
		if classBucket == nil {
			return fmt.Errorf("can't find class")
		}
		studentsBucket, err := classBucket.CreateBucketIfNotExists([]byte(KeyStudents))
		if err != nil {
			return err
		}
		err = studentsBucket.Put([]byte(userDetails.Email), nil)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		lgr.Printf("ERROR cannot edit class with Id: %s %v", body.Id, err)
		return openapi.Response(500, "{}"), nil
	}
	resp, err := classesWithOwnerDetails(a.db, userDetails.SchoolId, userDetails.Email)
	if err != nil {
		return openapi.Response(500, "{}"), nil
	}
	return openapi.Response(200, resp), nil

}

func placeBid(db *bolt.DB, clock Clock, userDetails UserInfo, item string, bid int32) (message string, err error) {
	err = db.Update(func(tx *bolt.Tx) error {
		message, err = placeBidtx(tx, clock, userDetails, item, bid)
		return err
	})

	return
}

func placeBidtx(tx *bolt.Tx, clock Clock, userDetails UserInfo, item string, bid int32) (message string, err error) {
	school, err := getSchoolBucketRx(tx, userDetails)
	if err != nil {
		return message, err
	}

	auctions := school.Bucket([]byte(KeyAuctions))
	if auctions == nil {
		return message, fmt.Errorf("cannot get auctions %s: %v", userDetails.Name, err)
	}

	newTime, err := time.Parse(time.RFC3339, item)
	if err != nil {
		return
	}

	item = newTime.Truncate(time.Millisecond).String()

	auctionByte := auctions.Get([]byte(item))
	if auctionByte == nil {
		return message, fmt.Errorf("cannot find the auction %s: %v", item, err)
	}

	var auction openapi.Auction
	err = json.Unmarshal(auctionByte, &auction)
	if err != nil {
		return message, err
	}

	if auction.WinnerId.Id == userDetails.Name {
		return message, fmt.Errorf("you are already winning this auction")
	}

	if clock.Now().After(auction.EndDate) {
		return message, fmt.Errorf("bid not accepted, auction expired")
	}

	if bid < auction.Bid {
		return message, fmt.Errorf("failed to outbid, try refreshing")
	}

	if bid < auction.MaxBid {
		auction.Bid = bid + 1
		marshal, err := json.Marshal(auction)
		if err != nil {
			return message, err
		}

		err = auctions.Put([]byte(item), marshal)
		if err != nil {
			return message, err
		}

		message = "You have been outbid"
		return message, err
	}

	if bid == auction.MaxBid {
		auction.Bid = auction.MaxBid
		marshal, err := json.Marshal(auction)
		if err != nil {
			return message, err
		}

		err = auctions.Put([]byte(item), marshal)
		if err != nil {
			return message, err
		}

		message = "You have not outbid the MaxBid"
		return message, err
	}

	err = chargeStudentUbuckTx(tx, clock, userDetails, decimal.NewFromInt32(bid), "Auction Bid "+item, true)
	if err != nil {
		return message, err
	}

	if auction.WinnerId.Id != "" {
		err = repayLosertx(tx, clock, auction.WinnerId.Id, auction.MaxBid, "Auction Refund "+item)
		if err != nil {
			return message, err
		}

	}

	auction.Bid = auction.MaxBid + 1
	auction.MaxBid = int32(bid)
	auction.WinnerId.Id = userDetails.Name

	if auction.TrueAuction && time.Until(auction.EndDate) < time.Minute {
		auction.EndDate = auction.EndDate.Add(time.Minute * 2)
	}

	marshal, err := json.Marshal(auction)
	if err != nil {
		return message, err
	}

	err = auctions.Put([]byte(item), marshal)
	if err != nil {
		return message, err
	}

	return

}

func repayLosertx(tx *bolt.Tx, clock Clock, winnerId string, amount int32, message string) (err error) {
	loser, err := getUserInLocalStoreTx(tx, winnerId)
	if err != nil {
		return err
	}
	err = addUbuck2StudentTx(tx, clock, loser, decimal.NewFromInt32(amount), message)
	if err != nil {
		return err
	}

	return
}

func NewStudentApiServiceImpl(db *bolt.DB, clock Clock) openapi.StudentApiServicer {
	return &StudentApiServiceImpl{
		db:    db,
		clock: clock,
	}
}
