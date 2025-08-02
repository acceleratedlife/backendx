package main

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/go-pkgz/auth/token"
	"github.com/go-pkgz/lgr"
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

func (a AllApiServiceImpl) ClearMessages(ctx context.Context) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	err = clearMessages(a.db, userDetails)

	if err != nil {
		return openapi.Response(400, nil), err
	}

	return openapi.Response(200, nil), nil

}

func (a AllApiServiceImpl) SearchTeachers(ctx context.Context) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	teachers, err := getTeachers(a.db, userDetails)
	if err != nil {
		return openapi.Response(400, ""), err
	}

	return openapi.Response(200, teachers), nil
}

func (a AllApiServiceImpl) SearchMarketItems(ctx context.Context, teacherId string) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	_, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	teacherDetails, err := getUserInLocalStore(a.db, teacherId)
	if err != nil {
		return openapi.Response(400, ""), err
	}

	items, err := getMarketItems(a.db, teacherDetails)
	if err != nil {
		return openapi.Response(400, ""), err
	}

	return openapi.Response(200, items), nil
}

func (a AllApiServiceImpl) IsPaused(ctx context.Context, id string) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	_, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	isPaused, err := isSchoolPaused(a.db, id)
	if err != nil {
		return openapi.Response(400, nil), err
	}

	if isPaused {
		return openapi.Response(201, nil), nil
	}
	return openapi.Response(200, nil), nil
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

	if userDetails.Role != UserRoleSysAdmin {
		isPaused, err := isSchoolPaused(a.db, userDetails.SchoolId)
		if err != nil {
			return openapi.Response(400, nil), err
		}
		if isPaused {
			return openapi.Response(400, map[string]string{
				"message": "your school is having an error. We are working on it and will have it back ASAP. Please try again later",
			}), nil
		}
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
			LottoPlay: userDetails.LottoPlay,
			LottoWin:  userDetails.LottoWin,
			Messages:  userDetails.Messages,
		}), nil
}

func (a AllApiServiceImpl) ConfirmEmail(ctx context.Context, s string) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AllApiServiceImpl) DeleteAuction(ctx context.Context, Id string) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	newTime, err := time.Parse(time.RFC3339, Id)
	if err != nil {
		return openapi.Response(400, nil), err
	}
	err = deleteAuction(a.db, userDetails, a.clock, newTime)

	if err != nil {
		lgr.Printf("ERROR cannot delete auction from the school: %s %v", userDetails.SchoolId, err)
		return openapi.Response(500, "{}"), err
	}

	return openapi.Response(200, nil), nil
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
		student, err := getStudentBucketRx(tx, userDetails.Name)
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
				Id:      to,
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

func (a *AllApiServiceImpl) MakeAuction(ctx context.Context, body openapi.RequestMakeAuction) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	isStaff := true
	if userDetails.Role == UserRoleStudent {
		settings, err := getSettings(a.db, userDetails)
		if err != nil {
			return openapi.Response(400, "{}"), err
		}

		if !settings.Student2student {
			return openapi.Response(400, ""), fmt.Errorf("disabled by administrator")
		}
		isStaff = false
	}

	err = MakeAuctionImpl(a.db, userDetails, body, isStaff)
	if err != nil {
		lgr.Printf("ERROR cannot make auctions from : %s %v", userDetails.Name, err)
		return openapi.Response(500, "{}"), err
	}

	return openapi.Response(200, nil), nil
}

func (a AllApiServiceImpl) SearchAccount(ctx context.Context, s string) (openapi.ImplResponse, error) {
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

func (s *AllApiServiceImpl) PayTransaction(ctx context.Context, body openapi.RequestPayTransaction) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(s.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	if userDetails.Role == UserRoleStudent {
		settings, err := getSettings(s.db, userDetails)
		if err != nil {
			return openapi.Response(400, ""), err
		}
		if !settings.Student2student {
			return openapi.Response(400, ""), fmt.Errorf("disabled by administrator")
		}

		err = executeStudentTransaction(s.db, s.clock, body.Amount, body.Student, userDetails, "")
		if err != nil {
			return openapi.Response(400, ""), err
		}

		if body.Amount > 0 {
			_, err = garnishHelper(s.db, s.clock, body, false)
			if err != nil {
				return openapi.Response(400, ""), err
			}
		}
	} else if userDetails.Role == UserRoleTeacher {
		err = executeTransaction(s.db, s.clock, body.Amount, body.Student, body.OwnerId, body.Description)
		if err != nil {
			return openapi.Response(400, ""), err
		}
		if body.Amount > 0 {
			_, err = garnishHelper(s.db, s.clock, body, true)
			if err != nil {
				return openapi.Response(400, ""), err
			}
		}
	} else {
		body.OwnerId = body.OwnerId[:1] + "." + body.OwnerId[1:]
		err = executeTransaction(s.db, s.clock, body.Amount, body.Student, body.OwnerId, body.Description)
		if err != nil {
			return openapi.Response(400, ""), err
		}
		if body.Amount > 0 {
			_, err = garnishHelper(s.db, s.clock, body, true)
			if err != nil {
				return openapi.Response(400, ""), err
			}
		}
	}

	if err != nil {
		return openapi.Response(400, ""), err
	}

	return openapi.Response(200, ""), nil
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
		searchedUser, err := getUserInLocalStoreTx(tx, Id)
		if err != nil {
			return err
		}

		nWorth := 0.0
		job := openapi.UserNoHistoryJob{}
		if searchedUser.Role == UserRoleStudent {
			nWorth, _ = StudentNetWorthTx(tx, searchedUser.Email).Float64()
			if searchedUser.College && searchedUser.CollegeEnd.IsZero() {
				job = getJobRx(tx, KeyCollegeJobs, searchedUser.Job)
			} else {
				job = getJobRx(tx, KeyJobs, searchedUser.Job)
			}

		}
		nUser := openapi.User{
			Id:               searchedUser.Email,
			CollegeEnd:       searchedUser.CollegeEnd,
			TransitionEnd:    searchedUser.TransitionEnd,
			FirstName:        searchedUser.FirstName,
			LastName:         searchedUser.LastName,
			Email:            searchedUser.Email,
			Confirmed:        searchedUser.Confirmed,
			SchoolId:         searchedUser.SchoolId,
			College:          searchedUser.College,
			Income:           searchedUser.Income,
			Role:             searchedUser.Role,
			Rank:             searchedUser.Rank,
			CareerTransition: searchedUser.CareerTransition,
			NetWorth:         float32(nWorth),
			Job:              job,
			TaxableIncome:    searchedUser.TaxableIncome,
		}
		resp = nUser

		return nil
	})

	if err != nil {
		lgr.Printf("ERROR cannot find the user: %s %v", Id, err)
		return openapi.Response(500, "{}"), err
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

func (s *AllApiServiceImpl) SearchClasses(ctx context.Context) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(s.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	var data []openapi.Class
	if userDetails.Role == UserRoleTeacher {
		data = getTeacherClasses(s.db, userDetails.SchoolId, userDetails.Name)
	} else if userDetails.Role == UserRoleAdmin {
		data = getSchoolClasses(s.db, userDetails.SchoolId)
	} else {
		data, err = getStudentClasses(s.db, userDetails)
	}

	if data == nil {
		return openapi.ImplResponse{}, err
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].Period < data[j].Period
	})
	return openapi.Response(200, data), nil

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
		student, err := getStudentBucketRx(tx, userDetails.Name)
		if err != nil {
			return err
		}

		accounts := student.Bucket([]byte(KeyAccounts))
		if accounts == nil {
			return fmt.Errorf("cannot find students buck accounts")
		}

		c := accounts.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if v != nil {
				continue
			}

			if string(k) != CurrencyUBuck && string(k) != KeyDebt && !strings.Contains(string(k), "@") {
				continue
			}

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

			account, err = getCBaccountDetailsRx(tx, userDetails, account) //is it really necessary to check if a cb account exists. This seems like it can be removed.
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
		CertificateOfDepositIfNeeded(a.db, a.clock, userDetails) //grows cds
		CollegeIfNeeded(a.db, a.clock, userDetails)              //assigns new job after college
		CareerIfNeeded(a.db, a.clock, userDetails)               //assigns new job
		DebtIfNeeded(a.db, a.clock, userDetails)                 //grows debt
		DailyPayIfNeeded(a.db, a.clock, userDetails)             //daily payment
		EventIfNeeded(a.db, a.clock, userDetails)                //generates event
		LotteryIfNeeded(a.db, a.clock, userDetails)              //grows lottery
	}

	var resp []openapi.UserNoHistory
	var ranked int
	err = a.db.View(func(tx *bolt.Tx) error {
		resp, ranked, err = getSchoolStudentsRx(tx, userDetails)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		lgr.Printf("ERROR cannot collect students from the school: %s %v", userDetails.SchoolId, err)
		return openapi.Response(500, "{}"), nil
	}
	return openapi.Response(200, resp[:ranked]), nil
}

func (a *AllApiServiceImpl) UserEdit(ctx context.Context, body openapi.RequestUserEdit) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	editDetails, err := getUserInLocalStore(a.db, body.Email)
	if err != nil {
		return openapi.Response(404, nil), err
	}

	if userDetails.Email != body.Email { //someone trying to edit someone else
		if userDetails.Role == UserRoleStudent { //this must not be a student
			return openapi.Response(401, ""), fmt.Errorf("students cannot edit other users")
		}
		if userDetails.Role == editDetails.Role { //teacher trying to edit other teacher or admin trying to edit other admin
			return openapi.Response(401, ""), fmt.Errorf("you are staff but you can't edit someone else with your same status")
		}
		userDetails, err = getUserInLocalStore(a.db, body.Email)
		if err != nil {
			return openapi.Response(404, nil), err
		}
	}

	err = userEdit(a.db, a.clock, userDetails, body)

	if err != nil {
		return openapi.Response(500, nil), err
	}

	userDetails, err = getUserInLocalStore(a.db, body.Email)
	if err != nil {
		return openapi.Response(500, nil), err
	}

	nWorth := 0.0
	if userDetails.Role == 0 {
		nWorth, _ = StudentNetWorth(a.db, userDetails.Name).Float64()
	}

	resp := openapi.User{
		Id:               userDetails.Name,
		Email:            userDetails.Email,
		CollegeEnd:       userDetails.CollegeEnd,
		TransitionEnd:    userDetails.TransitionEnd,
		FirstName:        userDetails.FirstName,
		LastName:         userDetails.LastName,
		Confirmed:        userDetails.Confirmed,
		SchoolId:         userDetails.SchoolId,
		CareerTransition: userDetails.CareerTransition,
		College:          userDetails.College,
		Income:           userDetails.Income,
		Role:             userDetails.Role,
		Rank:             userDetails.Rank,
		NetWorth:         float32(nWorth),
		TaxableIncome:    userDetails.TaxableIncome,
	}
	return openapi.Response(200, resp), nil //this is incomplete
}

// return each all bucks with their conversion ratio
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

		if Id != CurrencyUBuck && Id != KeyDebt && !strings.Contains(Id, "@") {
			continue
		}

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
