package main

import (
	"context"
	"encoding/json"

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

func (a AllApiServiceImpl) ExchangeRate(ctx context.Context, s string, s2 string) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a AllApiServiceImpl) Logout(ctx context.Context, s string) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a AllApiServiceImpl) ResetPassword(ctx context.Context, body openapi.UsersResetPasswordBody) (openapi.ImplResponse, error) {
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

func (a *AllApiServiceImpl) SearchClass(ctx context.Context, query openapi.RequestUser) (openapi.ImplResponse, error) {
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
		classBucket, _, err := getClassAtSchoolTx(tx, userDetails.SchoolId, query.Id)
		if err != nil {
			return err
		}
		resp.Id = query.Id
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
	panic("implement me")
}

func (a *AllApiServiceImpl) SearchStudent(ctx context.Context, query openapi.RequestUser) (openapi.ImplResponse, error) {
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
		user, err := getUserInLocalStoreTx(tx, query.Id)
		if err != nil {
			return err
		}

		nWorth, _ := StudentNetWorthTx(tx, user.Email).Float64()
		nUser := openapi.User{
			Id: user.Email,
			//CollegeEnd:    time.Time{},
			//TransitionEnd: time.Time{},
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Confirmed: user.Confirmed,
			SchoolId:  user.SchoolId,
			//College:       false,
			//Children:      0,
			Income:   10,
			Role:     1,
			Rank:     2,
			NetWorth: float32(nWorth),
		}
		resp = nUser

		return nil
	})

	if err != nil {
		lgr.Printf("ERROR cannot find the user: %s %v", query.Id, err)
		return openapi.Response(500, "{}"), nil
	}
	return openapi.Response(200, resp), nil
}

func (a AllApiServiceImpl) SearchStudentBuck(ctx context.Context, s string) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
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
		teachers := school.Bucket([]byte(KeyTeachers))
		if teachers == nil {
			lgr.Printf("WARN no teachers bucket in school %s", userDetails.SchoolId)
			return nil
		}

		studentsId := make(map[string]int, 0)

		iterateBuckets(teachers, func(teacher *bolt.Bucket, _ []byte) {
			iterateBuckets(teacher, func(class *bolt.Bucket, _ []byte) {
				students := class.Bucket([]byte(KeyStudents))
				if students == nil {
					return
				}
				c := students.Cursor()

				for k, _ := c.First(); k != nil; k, _ = c.Next() {
					studentsId[string(k)] = 0
				}

			})
		})

		users := tx.Bucket([]byte(KeyUsers))

		for k, _ := range studentsId {
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
				//CollegeEnd:    time.Time{},
				//TransitionEnd: time.Time{},
				FirstName: student.FirstName,
				LastName:  student.LastName,
				Email:     student.Email,
				Confirmed: student.Confirmed,
				SchoolId:  student.SchoolId,
				//College:       false,
				//Children:      0,
				Income:   10,
				Role:     student.Role,
				Rank:     2,
				NetWorth: float32(nWorth),
			}

			resp = append(resp, nUser)

		}

		return nil
	})

	if err != nil {
		lgr.Printf("ERROR cannot collect students from the school: %s %v", userDetails.SchoolId, err)
		return openapi.Response(500, "{}"), nil
	}
	return openapi.Response(200, resp), nil
}

func (a AllApiServiceImpl) UserEdit(ctx context.Context, body openapi.UsersUserBody) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

// NewAllApiServiceImpl provides real api
func NewAllApiServiceImpl(db *bolt.DB, clock Clock) openapi.AllApiServicer {
	return &AllApiServiceImpl{
		db:    db,
		clock: clock,
	}
}
