package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/go-pkgz/auth/token"
	"github.com/go-pkgz/lgr"
	"github.com/shopspring/decimal"
	bolt "go.etcd.io/bbolt"
)

type AllApiServiceImpl struct {
	db    *bolt.DB
	clock Clock
}

func (a *AllApiServiceImpl) Login(ctx context.Context, login openapi.RequestLogin) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AllApiServiceImpl) AuthUser(ctx context.Context) (user openapi.ImplResponse, err error) {

	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}
	return openapi.Response(200,
		openapi.ResponseAuth2{
			Email:     userDetails.Email,
			FirstName: userDetails.FirstName,
			LastName:  userDetails.LastName,
			IsAdmin:   userDetails.Role != 0,
			IsAuth:    true,
			Role:      userDetails.Role,
			SchoolId:  userDetails.SchoolId,
			Id:        userDetails.Name,
		}), nil
}

func (a AllApiServiceImpl) ConfirmEmail(ctx context.Context, s string) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a AllApiServiceImpl) ExchangeRate(ctx context.Context, from string, to string) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	var resp []openapi.ResponseCurrencyExchange
	err = a.db.View(func(tx *bolt.Tx) error {
		student, err := getStudentBucketRoTx(tx, userDetails.Name)
		if err != nil {
			return err
		}

		accounts := student.Bucket([]byte(KeyAccounts))
		if accounts == nil {
			return fmt.Errorf("cannot find students accounts")
		}

		fromName, err := getBuckNameTx(tx, from)
		if err != nil {
			return err
		}

		toName, err := getBuckNameTx(tx, to)
		if err != nil {
			return err
		}

		rate, err := xRateToBaseRx(tx, userDetails.SchoolId, from, to)
		if err != nil {
			return err
		}

		account := accounts.Bucket([]byte(from))
		if account == nil {
			resp = append(resp, openapi.ResponseCurrencyExchange{
				Conversion: float32(rate.InexactFloat64()),
				Balance:    0,
				Id:         from,
				Buck: openapi.ResponseCurrencyExchangeBuck{
					Name: fromName,
				},
			})
		} else {
			responseAccount, err := getStudentAccountRx(tx, account, from)
			if err != nil {
				return err
			}

			responseAccount.Conversion = float32(rate.InexactFloat64())
			responseAccount.Buck.Name = fromName

			resp = append(resp, responseAccount)
		}

		account = accounts.Bucket([]byte(to))
		if account == nil {
			resp = append(resp, openapi.ResponseCurrencyExchange{
				Balance: 0,
				Id:      from,
				Buck: openapi.ResponseCurrencyExchangeBuck{
					Name: toName,
				},
			})
		} else {
			responseAccount, err := getStudentAccountRx(tx, account, to)
			if err != nil {
				return err
			}

			responseAccount.Buck.Name = toName

			resp = append(resp, responseAccount)
		}

		return nil

	})

	if err != nil {
		return openapi.Response(400, nil), err
	}

	return openapi.Response(200, resp), nil

}

func (a AllApiServiceImpl) Logout(ctx context.Context, s string) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a AllApiServiceImpl) SearchAccount(ctx context.Context, s string) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a AllApiServiceImpl) SearchBucks(ctx context.Context, s string) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AllApiServiceImpl) SearchClass(ctx context.Context, Id string) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}
	var resp openapi.ClassWithMembers
	err = a.db.View(func(tx *bolt.Tx) error {
		classBucket, _, err := getClassAtSchoolTx(tx, userDetails.SchoolId, Id)
		if err != nil {
			return err
		}
		resp.Id = Id
		resp.AddCode = string(classBucket.Get([]byte(KeyAddCode)))
		// resp.OwnerId = string(class.Get([]byte("ownerId")))
		resp.Period = btoi32(classBucket.Get([]byte(KeyPeriod)))
		resp.Name = string(classBucket.Get([]byte(KeyName)))
		Members, err := PopulateClassMembers(tx, classBucket)
		if err != nil {
			return err
		}
		resp.Members = Members
		return nil
	})
	if err != nil {
		lgr.Printf("ERROR cannot collect classes from the school: %s %v", userDetails.SchoolId, err)
		return openapi.Response(500, "{}"), nil
	}
	return openapi.Response(200, resp), nil
}

func (a AllApiServiceImpl) SearchSchool(ctx context.Context, s string) (openapi.ImplResponse, error) {
	//TODO implement me
	//depricated
	panic("implement me")
}

func (a *AllApiServiceImpl) SearchStudent(ctx context.Context, Id string) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	_, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	var resp openapi.User
	err = a.db.View(func(tx *bolt.Tx) error {
		user, err := getUserInLocalStoreTx(tx, Id)
		if err != nil {
			return err
		}

		history, err := getStudentHistoryTX(tx, user.Name)
		if err != nil {
			return err
		}

		nWorth, _ := StudentNetWorthTx(tx, user.Email).Float64()
		nUser := openapi.User{
			Id:               user.Email,
			CollegeEnd:       user.CollegeEnd,
			TransitionEnd:    user.TransitionEnd,
			FirstName:        user.FirstName,
			LastName:         user.LastName,
			Email:            user.Email,
			Confirmed:        user.Confirmed,
			SchoolId:         user.SchoolId,
			College:          user.College,
			Children:         user.Children,
			Income:           user.Income,
			Role:             0,
			Rank:             user.Rank,
			History:          history,
			CareerTransition: user.CareerTransition,
			NetWorth:         float32(nWorth),
		}
		resp = nUser

		return nil
	})

	if err != nil {
		lgr.Printf("ERROR cannot find the user: %s %v", Id, err)
		return openapi.Response(500, "{}"), nil
	}
	return openapi.Response(200, resp), nil
}

func (s *AllApiServiceImpl) SearchAllBucks(ctx context.Context) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(s.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}
	var resp []openapi.Buck
	err = s.db.View(func(tx *bolt.Tx) error {
		bucks, err := getCBBucksRx(tx, userDetails.SchoolId)
		if err != nil {
			return err
		}

		resp = bucks

		return nil
	})

	if err != nil {
		return openapi.Response(400, nil), err
	}

	sort.Slice(resp, func(i, j int) bool {
		return strings.ToLower(resp[i].Name) < strings.ToLower(resp[j].Name)
	})

	return openapi.Response(200, resp), nil

}

func (s *AllApiServiceImpl) SearchStudentBucks(ctx context.Context) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(s.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	var resp []openapi.ResponseCurrencyExchange
	err = s.db.View(func(tx *bolt.Tx) error {
		student, err := getStudentBucketRoTx(tx, userDetails.Name)
		if err != nil {
			return err
		}

		accounts := student.Bucket([]byte(KeyAccounts))
		if accounts == nil {
			return fmt.Errorf("cannot find students buck accounts")
		}

		c := accounts.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			account, err := getStudentAccountRx(tx, accounts.Bucket(k), string(k))
			if err != nil {
				return err
			}

			if account.Balance <= 0 {
				continue
			}

			if account.Id == CurrencyUBuck {
				account.Buck.Name = "UBuck"
				account.Conversion = 1
			} else if account.Id == KeyDebt {
				account.Buck.Name = "Debt"
				account.Conversion = -1
			} else {
				conversion, err := xRateToBaseRx(tx, userDetails.SchoolId, account.Id, "")
				if err != nil {
					return err
				}
				account.Conversion = float32(conversion.InexactFloat64())
				owner, err := getUserInLocalStoreTx(tx, account.Id)
				if err != nil {
					return err
				}
				account.Buck.Name = owner.LastName + " Buck"
			}

			account, err = getCBaccountDetailsRx(tx, userDetails, account)
			if err != nil {
				return err
			}

			resp = append(resp, account)
		}

		return nil
	})

	if err != nil {
		return openapi.Response(400, nil), err
	}
	return openapi.Response(200, resp), nil
}

func (a *AllApiServiceImpl) SearchStudents(ctx context.Context) (openapi.ImplResponse, error) {

	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	if userDetails.Role == UserRoleStudent {
		DailyPayIfNeeded(a.db, a.clock, userDetails)
	}

	resp := make([]openapi.UserNoHistory, 0)
	err = a.db.View(func(tx *bolt.Tx) error {
		school, err := SchoolByIdTx(tx, userDetails.SchoolId)
		if err != nil {
			return err
		}

		students := school.Bucket([]byte(KeyStudents))
		if students == nil {
			return fmt.Errorf("cannot find students bucket")
		}

		c := students.Cursor()

		users := tx.Bucket([]byte(KeyUsers))

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			studentData := users.Get([]byte(k))
			var student UserInfo
			err = json.Unmarshal(studentData, &student)
			if err != nil {
				lgr.Printf("ERROR cannot unmarshal userInfo for %s", k)
				continue
			}
			if student.Role != UserRoleStudent {
				lgr.Printf("ERROR student %s has role %s", k, userData.Role)
				continue
			}

			nWorth, _ := StudentNetWorthTx(tx, student.Name).Float64()
			nUser := openapi.UserNoHistory{
				Id: student.Email,
				// CollegeEnd:    time.Time{},
				//TransitionEnd: time.Time{},
				FirstName: student.FirstName,
				LastName:  student.LastName,
				// Email:     student.Email,
				// Confirmed: student.Confirmed,
				// SchoolId:  student.SchoolId,
				//College:       false,
				//Children:      0,
				// Income:   10,
				// Role:     student.Role,
				Rank:     student.Rank,
				NetWorth: float32(nWorth),
			}

			resp = append(resp, nUser)

		}

		return nil

	})

	sort.SliceStable(resp, func(i, j int) bool {
		return resp[i].NetWorth > resp[j].NetWorth
	})

	for i := 0; i < len(resp); i++ {
		resp[i].Rank = int32(i + 1)
	}

	ranked, err2 := saveRanks(a.db, resp)
	if err2 != nil {
		lgr.Printf("ERROR saving students ranks: %s %v", userDetails.SchoolId, err)
	}

	if err != nil || err2 != nil {
		lgr.Printf("ERROR cannot collect students from the school: %s %v", userDetails.SchoolId, err)
		return openapi.Response(500, "{}"), nil
	}
	return openapi.Response(200, resp[:ranked]), nil
}

func (a *AllApiServiceImpl) UserEdit(ctx context.Context, body openapi.UsersUserBody) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	err = a.db.Update(func(tx *bolt.Tx) error {
		users := tx.Bucket([]byte(KeyUsers))
		if users == nil {
			return fmt.Errorf("users do not exist")
		}

		user := users.Get([]byte(userData.Name))

		if user == nil {
			return fmt.Errorf("user does not exist")
		}

		if body.FirstName != "" {
			userDetails.FirstName = body.FirstName
		}
		if body.LastName != "" {
			userDetails.LastName = body.LastName
		}
		if len(body.Password) > 5 {
			userDetails.PasswordSha = EncodePassword(body.Password)
		}
		if body.CareerTransition && !userDetails.CareerTransition {
			userDetails.CareerTransition = true
			userDetails.TransitionEnd = a.clock.Now().AddDate(0, 0, 7) //7 days
			userDetails.Income = userDetails.Income / 2
		}
		if body.College && !userDetails.College {
			cost := decimal.NewFromFloat(math.Floor(rand.Float64()*(8000-5000) + 8000))
			err := chargeStudentUbuckTx(tx, a.clock, userDetails, cost, "Paying for College", false)
			if err != nil {
				return fmt.Errorf("Failed to chargeStudentUbuckTx: %v", err)
			}
			userDetails.College = true
			userDetails.CollegeEnd = a.clock.Now().AddDate(0, 0, 28) //28 days
			userDetails.Income = userDetails.Income / 2
		}

		marshal, err := json.Marshal(userDetails)
		if err != nil {
			return fmt.Errorf("Failed to Marshal userDetails")
		}
		err = users.Put([]byte(userData.Name), marshal)
		if err != nil {
			return fmt.Errorf("Failed to Put userDetails")
		}
		return nil
	})

	if err != nil {
		return openapi.Response(500, nil), err
	}

	userDetails, err = getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(500, nil), err
	}

	history, err := getStudentHistory(a.db, userDetails.Name, userDetails.SchoolId)
	if err != nil {
		return openapi.Response(500, nil), err
	}

	nWorth, _ := StudentNetWorth(a.db, userDetails.Name).Float64()
	resp := openapi.User{
		Id:               userDetails.Name,
		Email:            userDetails.Email,
		CollegeEnd:       userDetails.CollegeEnd,
		TransitionEnd:    userDetails.TransitionEnd,
		FirstName:        userDetails.FirstName,
		LastName:         userDetails.LastName,
		History:          history,
		Confirmed:        userDetails.Confirmed,
		SchoolId:         userDetails.SchoolId,
		CareerTransition: userDetails.CareerTransition,
		College:          userDetails.College,
		Children:         userDetails.Children,
		Income:           userDetails.Income,
		Role:             userDetails.Role,
		Rank:             userDetails.Rank,
		NetWorth:         float32(nWorth),
	}
	return openapi.Response(200, resp), nil //this is incomplete
}

func getCBBucksRx(tx *bolt.Tx, schoolId string) (bucks []openapi.Buck, err error) {
	cb, err := getCbRx(tx, schoolId)
	if err != nil {
		return bucks, err
	}

	accounts := cb.Bucket([]byte(KeyAccounts))
	if accounts == nil {
		return bucks, fmt.Errorf("cannot get CB accounts")
	}

	c := accounts.Cursor()
	for k, _ := c.First(); k != nil; k, _ = c.Next() {

		Id := string(k)

		var teacher UserInfo
		var ratio float32
		if CurrencyUBuck == Id {
			teacher.LastName = "UBuck"
			ratio = 1
		} else if KeyDebt == Id {
			teacher.LastName = "Debt"
			ratio = -1
		} else {
			teacher, err = getUserInLocalStoreTx(tx, Id)
			if err != nil {
				return bucks, err
			}

			teacher.LastName = teacher.LastName + " Buck"
			rate, err := xRateToBaseRx(tx, schoolId, Id, "")
			if err != nil {
				return bucks, err
			}

			ratio = float32(rate.InexactFloat64())

		}

		bucks = append(bucks, openapi.Buck{
			Id:    Id,
			Name:  teacher.LastName,
			Ratio: ratio,
		})
	}

	return
}

// NewAllApiServiceImpl provides real api
func NewAllApiServiceImpl(db *bolt.DB, clock Clock) openapi.AllApiServicer {
	return &AllApiServiceImpl{
		db:    db,
		clock: clock,
	}
}
