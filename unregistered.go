package main

import (
	"context"
	openapi "github.com/acceleratedlife/backend/go"
	bolt "go.etcd.io/bbolt"
)

type UnregisteredApiServiceImpl struct {
	db *bolt.DB
}

func (u *UnregisteredApiServiceImpl) Register(ctx context.Context, register openapi.RequestRegister) (openapi.ImplResponse, error) {

	role, pathId, err := RoleByAddCode(u.db, register.AddCode)
	if err != nil {
		return openapi.Response(404,
			openapi.ResponseRegister4{
				Message: err.Error(),
			}), nil
	}

	if role == UserRoleTeacher {
		newUser := UserInfo{
			Name:        register.Email,
			FirstName:   register.FirstName,
			LastName:    register.LastName,
			Email:       register.Email,
			Confirmed:   false,
			PasswordSha: EncodePassword(register.Password),
			SchoolId:    pathId.schoolId,
			Role:        role,
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
		newUser := UserInfo{
			Name:        register.Email,
			FirstName:   register.FirstName,
			LastName:    register.LastName,
			Email:       register.Email,
			Confirmed:   false,
			PasswordSha: EncodePassword(register.Password),
			SchoolId:    pathId.schoolId,
			Role:        role,
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
func NewUnregisteredApiServiceImpl(db *bolt.DB) openapi.UnregisteredApiServicer {
	return &UnregisteredApiServiceImpl{
		db: db,
	}
}
