package main

import (
	"context"
	"fmt"
	"sort"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/go-pkgz/auth/token"
	"github.com/go-pkgz/lgr"
	bolt "go.etcd.io/bbolt"
)

type StaffApiServiceImpl struct {
	db *bolt.DB
}

func (s StaffApiServiceImpl) SearchEvents(ctx context.Context, s2 string) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s StaffApiServiceImpl) DeleteAuction(ctx context.Context, s2 string) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *StaffApiServiceImpl) Deleteclass(ctx context.Context, body openapi.RequestUser) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}
	if userDetails.Role == UserRoleStudent {
		return openapi.Response(401, ""), nil
	}
	var resp openapi.Class

	err = a.db.Update(func(tx *bolt.Tx) error {
		classBucket, err := getClassAtSchoolTx(tx, userDetails.SchoolId, body.Id)
		if err != nil {
			return err
		}
		path := classBucket.Tx().DB().Path()
		println(path)
		root := classBucket.Root()
		println(root)
		studentsBucket := classBucket.Bucket([]byte(KeyStudents))
		members, err := studentsToSlice(studentsBucket)
		if err != nil {
			return err
		}
		resp = openapi.Class{
			Id:      body.Id,
			AddCode: string(classBucket.Get([]byte(KeyAddCode))),
			Period:  btoi32(classBucket.Get([]byte(KeyPeriod))),
			Name:    string(classBucket.Get([]byte(KeyName))),
			Members: members,
		}
		return nil
	})
	if err != nil {
		lgr.Printf("ERROR cannot collect classes from the school: %s %v", userDetails.SchoolId, err)
		return openapi.Response(500, "{}"), nil
	}
	return openapi.Response(200, resp), nil
}

func (a *StaffApiServiceImpl) EditClass(ctx context.Context, body openapi.RequestEditClass) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}
	if userDetails.Role == UserRoleStudent {
		return openapi.Response(401, ""), nil
	}
	var class openapi.Class
	err = a.db.Update(func(tx *bolt.Tx) error {
		classBucket, err := getClassAtSchoolTx(tx, userDetails.SchoolId, body.Id)
		if err != nil {
			return err
		}
		studentsBucket := classBucket.Bucket([]byte(KeyStudents))
		members, err := studentsToSlice(studentsBucket)
		if err != nil {
			return err
		}
		err = classBucket.Put([]byte(KeyName), []byte(body.Name))
		if err != nil {
			return err
		}
		err = classBucket.Put([]byte(KeyPeriod), itob32(int32(body.Period)))
		if err != nil {
			return err
		}
		class = openapi.Class{
			Id:      body.Id,
			OwnerId: string(classBucket.Get([]byte("OwnerId"))), //OnwerId is legacy, not sure if it is needed
			Period:  btoi32(classBucket.Get([]byte(KeyPeriod))),
			Name:    string(classBucket.Get([]byte(KeyName))),
			AddCode: string(classBucket.Get([]byte(KeyAddCode))),
			Members: members,
		}
		return nil
	})
	if err != nil {
		lgr.Printf("ERROR cannot edit class with Id: %s %v", body.Id, err)
		return openapi.Response(500, "{}"), nil
	}
	return openapi.Response(200, class), nil
}

func (s *StaffApiServiceImpl) KickClass(ctx context.Context, body openapi.RequestKickClass) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s StaffApiServiceImpl) MakeAuction(ctx context.Context, s2 string, body openapi.AuctionsBody) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StaffApiServiceImpl) MakeClass(ctx context.Context, request openapi.RequestMakeClass) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(s.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}
	if userDetails.Role != UserRoleTeacher {
		return openapi.Response(401, ""), nil
	}

	classes, err := s.MakeClassImpl(userDetails, request)

	if err != nil {
		lgr.Printf("ERROR class not created: %v", err)
		return openapi.Response(500, "{}"), nil
	}

	return openapi.Response(200, classes), nil
}

func (s StaffApiServiceImpl) PayTransaction(ctx context.Context, body openapi.TransactionsPayTransactionBody) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s StaffApiServiceImpl) PayTransactions(ctx context.Context, body openapi.TransactionsPayTransactionsBody) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s StaffApiServiceImpl) SearchAllBucks(ctx context.Context, s2 string) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s StaffApiServiceImpl) SearchAuctionsTeacher(ctx context.Context, s2 string) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StaffApiServiceImpl) SearchClasses(ctx context.Context, query openapi.RequestUser) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(s.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}
	if userDetails.Role == UserRoleStudent {
		return openapi.Response(401, ""), nil
	}

	var data []openapi.Class
	if userDetails.Role == UserRoleTeacher {
		data = getTeacherClasses(s.db, userDetails.SchoolId, userDetails.Name)
	} else if userDetails.Role == UserRoleAdmin {
		data = getSchoolClasses(s.db, userDetails.SchoolId)
	}

	if data == nil {
		return openapi.ImplResponse{}, err
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].Period < data[j].Period
	})
	return openapi.Response(200,
		data), nil

}

func (s StaffApiServiceImpl) SearchTransactions(ctx context.Context, s2 string) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

// NewStaffApiServiceImpl creates a default api service
func NewStaffApiServiceImpl(db *bolt.DB) openapi.StaffApiServicer {
	return &StaffApiServiceImpl{
		db: db,
	}
}

func studentsToSlice(students *bolt.Bucket) ([]string, error) {
	Members := make([]string, 0)
	cstudents := students.Cursor()
	for k, v := cstudents.First(); k != nil; k, v = cstudents.Next() {
		if v == nil {
			return nil, fmt.Errorf("members of class not found")
		}
		Members = append(Members, string(k))
	}
	return Members, nil
}
