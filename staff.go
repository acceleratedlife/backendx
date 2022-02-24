package main

import (
	"context"
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

func (s *StaffApiServiceImpl) Deleteclass(ctx context.Context, query openapi.RequestUser) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StaffApiServiceImpl) EditClass(ctx context.Context, body openapi.RequestEditClass) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
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
