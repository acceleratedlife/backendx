package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
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
		message, err = placeBidtx(tx, a.clock, userDetails, body.Item, body.Bid)
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

	err = a.db.Update(func(tx *bolt.Tx) error {
		err := studentConvertTx(tx, a.clock, userDetails, decimal.NewFromFloat32(body.Amount), body.AccountFrom, body.AccountTo, true)
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
func (a *StudentApiServiceImpl) CryptoConvert(context.Context, string, openapi.TransactionCryptoTransactionBody) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
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
	err = a.db.View(func(tx *bolt.Tx) error {
		auctionsBucket, err := getAuctionsTx(tx, userDetails)
		if err != nil {
			return err
		}

		auctions, err := getStudentAuctionsTx(tx, auctionsBucket, userDetails)
		if err != nil {
			return err
		}

		for _, auction := range auctions {

			now := time.Now()
			if (auction.StartDate.Before(now) && auction.EndDate.After(now)) || auction.WinnerId.Id == userDetails.Name {
				iAuction := openapi.ResponseAuctionStudent{
					Id:          auction.Id,
					Bid:         float32(auction.Bid),
					Description: auction.Description,
					EndDate:     auction.EndDate,
					StartDate:   auction.StartDate,
					OwnerId: openapi.ResponseAuctionStudentOwnerId{
						Id:       auction.OwnerId.Id,
						LastName: auction.OwnerId.LastName,
					},
					WinnerId: openapi.ResponseAuctionStudentWinnerId{
						Id:        auction.WinnerId.Id,
						FirstName: auction.WinnerId.FirstName,
						LastName:  auction.WinnerId.LastName,
					},
				}

				resp = append(resp, iAuction)
			}
		}
		return nil
	})

	if err != nil {
		return openapi.Response(400, err), nil
	}

	return openapi.Response(200, resp), nil
}
func (a *StudentApiServiceImpl) SearchBuckTransaction(ctx context.Context) (openapi.ImplResponse, error) {
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

		student, err := getStudentBucketRoTx(tx, userDetails.Name)
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
func (a *StudentApiServiceImpl) SearchCrypto(context.Context, string) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}
func (a *StudentApiServiceImpl) SearchCryptoTransaction(context.Context, string) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}
func (a *StudentApiServiceImpl) SearchStudentCrypto(context.Context) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
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

//StudentAddClass(context.Context, RequestAddClass) (ImplResponse, error)

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

	_, pathId, err := RoleByAddCode(a.db, body.AddCode)
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

func placeBidtx(tx *bolt.Tx, clock Clock, userDetails UserInfo, item time.Time, bid int32) (message string, err error) {
	school, err := getSchoolBucketTx(tx, userDetails)
	if err != nil {
		return message, err
	}

	auctions := school.Bucket([]byte(KeyAuctions))
	if auctions == nil {
		return message, fmt.Errorf("cannot get auctions %s: %v", userDetails.Name, err)
	}

	// item = item.Truncate(time.Millisecond)
	auctionByte := auctions.Get([]byte(item.String()))
	if auctionByte == nil {
		return message, fmt.Errorf("cannot find the auction %s: %v", userDetails.Name, err)
	}

	var auction openapi.Auction
	err = json.Unmarshal(auctionByte, &auction)
	if err != nil {
		return message, err
	}

	if auction.WinnerId.Id == userDetails.Name {
		return message, fmt.Errorf("You are already winning this auction")
	}

	if auction.EndDate.Sub(time.Now()) <= 0 {
		return message, fmt.Errorf("Bid not accepted, auction expired")
	}

	if bid < auction.Bid {
		return message, fmt.Errorf("Failed to outbid, try refreshing")
	}

	if bid < auction.MaxBid {
		auction.Bid = bid + 1
		marshal, err := json.Marshal(auction)
		if err != nil {
			return message, err
		}

		err = auctions.Put([]byte(item.String()), marshal)
		if err != nil {
			return message, err
		}

		message = "You have been outbid"
		return message, err
	}

	err = chargeStudentUbuckTx(tx, clock, userDetails, decimal.NewFromInt32(bid), "Auction Bid "+strconv.Itoa(item.Second()), true)
	if err != nil {
		return message, err
	}

	if auction.WinnerId.Id != "" {
		loser, err := getUserInLocalStoreTx(tx, auction.WinnerId.Id)
		if err != nil {
			return message, err
		}
		err = addUbuck2StudentTx(tx, clock, loser, decimal.NewFromInt32(auction.MaxBid), "Auction Refund"+strconv.Itoa(item.Second()))
		if err != nil {
			return message, err
		}
	}

	auction.Bid = auction.MaxBid + 1
	auction.MaxBid = int32(bid)
	auction.WinnerId.Id = userDetails.Name

	if time.Until(auction.EndDate) < time.Minute {
		auction.EndDate = auction.EndDate.Add(time.Minute * 2)
	}

	marshal, err := json.Marshal(auction)
	if err != nil {
		return message, err
	}

	err = auctions.Put([]byte(item.String()), marshal)
	if err != nil {
		return message, err
	}

	return

}

func NewStudentApiServiceImpl(db *bolt.DB, clock Clock) openapi.StudentApiServicer {
	return &StudentApiServiceImpl{
		db:    db,
		clock: clock,
	}
}

//How to calculate netWorth
//givens:
//student accounts:
//uBucks 0, Kirill Bucks 5
//
//uBuck total currency: 1000
//Kirill Bucks total currency: 100
//conversion ratio 1000/100 = 10
//10 ubucks = 1 kirill buck
//networth = uBucks + Kirill bucks *10
//50 = 0 + (5*10)
//networth = 50
