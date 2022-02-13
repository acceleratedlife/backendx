package main

import (
	"context"
	openapi "github.com/acceleratedlife/backend/go"
	bolt "go.etcd.io/bbolt"
)

type SysAdminApiServiceImpl struct {
	db *bolt.DB
}

func (s SysAdminApiServiceImpl) CreateBuck(ctx context.Context, body1 openapi.BucksBuckBody1) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s SysAdminApiServiceImpl) DeleteAccount(ctx context.Context, s2 string) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s SysAdminApiServiceImpl) DeleteBuck(ctx context.Context, s2 string) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s SysAdminApiServiceImpl) DeleteSchool(ctx context.Context, s2 string) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s SysAdminApiServiceImpl) DeleteUser(ctx context.Context, s2 string) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s SysAdminApiServiceImpl) Deletetransaction(ctx context.Context, s2 string) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s SysAdminApiServiceImpl) EditAccount(ctx context.Context, body openapi.AccountsAccountBody) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s SysAdminApiServiceImpl) EditBuck(ctx context.Context, body openapi.BucksBuckBody) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s SysAdminApiServiceImpl) EditSchool(ctx context.Context, body openapi.SchoolsSchoolBody) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s SysAdminApiServiceImpl) MakeAccount(ctx context.Context, body1 openapi.AccountsAccountBody1) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s SysAdminApiServiceImpl) MakeSchool(ctx context.Context, body1 openapi.SchoolsSchoolBody1) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s SysAdminApiServiceImpl) SearchAllBucks(ctx context.Context, s2 string) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SysAdminApiServiceImpl) SearchSchools(ctx context.Context, zip int32) (openapi.ImplResponse, error) {

	res, err := schoolsByZip(s.db, zip)

	if res != nil {
		return openapi.Response(500, nil), err
	}
	return openapi.Response(200, res), nil
}

func (s SysAdminApiServiceImpl) SearchTransaction(ctx context.Context, s2 string) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewSysAdminApiServiceImpl(db *bolt.DB) openapi.SysAdminApiServicer {
	return &SysAdminApiServiceImpl{
		db: db,
	}
}
