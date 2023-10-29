package main

import (
	"context"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	openapi "github.com/acceleratedlife/backend/go"
	bolt "go.etcd.io/bbolt"
)

type UnregisteredApiServiceImpl struct {
	db    *bolt.DB
	clock Clock
}

func (s *UnregisteredApiServiceImpl) GetCryptos(ctx context.Context) (openapi.ImplResponse, error) {

	cryptos, err := getCryptos(s.db)

	if err != nil {
		return openapi.Response(404, ""), err
	}

	return openapi.Response(200, cryptos), nil
}

func (s *UnregisteredApiServiceImpl) ResetStaffPassword(ctx context.Context, body openapi.RequestUser) (openapi.ImplResponse, error) {
	staffDetails, err := getUserInLocalStore(s.db, body.Id)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	if staffDetails.Role == UserRoleStudent {
		return openapi.Response(401, ""), nil
	}

	err = s.db.Update(func(tx *bolt.Tx) error {
		resp, err := resetPasswordTx(tx, staffDetails, 3)
		if err != nil {
			return err
		}

		err = sendEmail(staffDetails, resp.Password)
		return err
	})

	if err != nil {
		return openapi.Response(400, ""), err
	}

	return openapi.Response(200, ""), nil
}

func (u *UnregisteredApiServiceImpl) Register(ctx context.Context, register openapi.RequestRegister) (openapi.ImplResponse, error) {

	_, err := getUserInLocalStore(u.db, register.Email)
	if err == nil {
		return openapi.Response(404,
			openapi.ResponseRegister4{
				Message: "That email is already in use",
			}), nil
	}

	role, pathId, err := RoleByAddCode(u.db, register.AddCode, u.clock)
	if err != nil {
		return openapi.Response(404,
			openapi.ResponseRegister4{
				Message: err.Error(),
			}), nil
	}

	emailPattern, err := regexp.Compile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`)
	if err != nil {
		return openapi.Response(404,
			openapi.ResponseRegister4{
				Message: err.Error(),
			}), nil
	}

	if role == UserRoleTeacher {
		if !emailPattern.MatchString(register.Email) {
			return openapi.Response(404,
				openapi.ResponseRegister4{
					Message: "That is not a valid email address",
				}), nil
		}

		newUser := UserInfo{
			Name:        register.Email,
			FirstName:   register.FirstName,
			LastName:    register.LastName,
			Email:       register.Email,
			Confirmed:   false,
			PasswordSha: EncodePassword(register.Password),
			SchoolId:    pathId.schoolId,
			Role:        role,
			Settings: TeacherSettings{
				CurrencyLock: false,
			},
		}
		err = createTeacher(u.db, newUser)

		if err != nil {
			return openapi.Response(404,
				openapi.ResponseRegister4{
					Message: err.Error(),
				}), nil
		}

		return openapi.Response(200,
			openapi.ResponseRegister2{
				Success: true,
			}), nil
	}

	if role == UserRoleStudent {
		if emailPattern.MatchString(register.Email) {
			return openapi.Response(404,
				openapi.ResponseRegister4{
					Message: "Students cannot register with an email address",
				}), nil
		}

		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		newUser := UserInfo{
			Name:        register.Email,
			FirstName:   register.FirstName,
			LastName:    string(register.LastName[0]) + strconv.Itoa(r.Intn(10000)),
			Email:       register.Email,
			Confirmed:   false,
			PasswordSha: EncodePassword(register.Password),
			SchoolId:    pathId.schoolId,
			Role:        role,
			Job:         getJobId(u.db, KeyJobs),
		}
		err = createStudent(u.db, newUser, pathId)

		if err != nil {
			return openapi.Response(404,
				openapi.ResponseRegister4{
					Message: err.Error(),
				}), nil
		}
		return openapi.Response(200,
			openapi.ResponseRegister2{
				Success: true,
			}), nil
	}

	if role == UserRoleAdmin {
		return openapi.Response(403,
			openapi.ResponseRegister4{
				Message: "not allowed",
			}), nil

	}

	return openapi.Response(404,
		openapi.ResponseRegister4{
			Message: "not implemented",
		}), nil
}

// NewUnregisteredApiServiceImpl creates a default api service
func NewUnregisteredApiServiceImpl(db *bolt.DB, clock Clock) openapi.UnregisteredApiServicer {
	return &UnregisteredApiServiceImpl{
		db:    db,
		clock: clock,
	}
}
