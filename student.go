package main

import (
	"context"
	"fmt"
	"time"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/go-pkgz/auth/token"
	"github.com/go-pkgz/lgr"
	bolt "go.etcd.io/bbolt"
)

type StudentApiServiceImpl struct {
	db *bolt.DB
}

func (a *StudentApiServiceImpl) AuctionBid(ctx context.Context, auctionsPlaceBidBody openapi.RequestAuctionBid) (openapi.ImplResponse, error) {
	panic("implement me")
}

func (a *StudentApiServiceImpl) BuckConvert(context.Context, string, openapi.TransactionsConversionTransactionBody) (response openapi.ImplResponse, err error) {
	//TODO implement me
	panic("implement me")
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
			start, err := time.Parse(KeyTime, auction.StartDate)
			if err != nil {
				return fmt.Errorf("cannot get start time parse %s: %v", userDetails.Name, err)
			}

			end, err := time.Parse(KeyTime, auction.EndDate)
			if err != nil {
				return fmt.Errorf("cannot get end time parse %s: %v", userDetails.Name, err)
			}

			now := time.Now()
			if (start.Before(now) && end.After(now)) || auction.WinnerId.Id == userDetails.Name {
				iAuction := openapi.ResponseAuctionStudent{
					Id:          auction.Id,
					Bid:         float32(auction.Bid),
					Description: auction.Description,
					EndDate:     end,
					StartDate:   start,
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

		bAccounts := student.Bucket([]byte(KeybAccounts))
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

func NewStudentApiServiceImpl(db *bolt.DB) openapi.StudentApiServicer {
	return &StudentApiServiceImpl{
		db: db,
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
