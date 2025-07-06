package main

import (
	"context"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/go-pkgz/auth/token"
	bolt "go.etcd.io/bbolt"
)

type SysAdminApiServiceImpl struct {
	db         *bolt.DB
	jwtService *token.Service
}

func (s SysAdminApiServiceImpl) DeleteAccount(ctx context.Context, s2 string) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s SysAdminApiServiceImpl) DeleteSchool(ctx context.Context, s2 string) (openapi.ImplResponse, error) {
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

func (a SysAdminApiServiceImpl) GetSchools(ctx context.Context) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	if userDetails.Role != UserRoleSysAdmin {
		return openapi.Response(401, ""), nil
	}

	resp, err := getSchools(a.db)

	if err != nil {
		return openapi.Response(400, nil), err
	}

	return openapi.Response(200, resp), nil
}

func (a SysAdminApiServiceImpl) GetSchoolUsers(ctx context.Context, schoolId string) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	if userDetails.Role != UserRoleSysAdmin {
		return openapi.Response(401, ""), nil
	}

	resp, err := getSchoolUsers(a.db, schoolId)

	if err != nil {
		return openapi.Response(400, nil), err
	}

	return openapi.Response(200, resp), nil
}

func (a SysAdminApiServiceImpl) ImpersonateUser(ctx context.Context, user openapi.RequestImpersonate) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	if userDetails.Role != UserRoleSysAdmin {
		return openapi.Response(401, ""), nil
	}

	tgt, err := getUserInLocalStore(a.db, user.UserId)
	if err != nil {
		return openapi.Response(400, nil), err
	}

	jwtStr, xsrf, err := makeToken(a.jwtService, tgt)

	if err != nil {
		return openapi.Response(400, nil), err
	}

	return openapi.Response(200, openapi.ResponseImpersonate{
		Token: jwtStr,
		Xsrf:  xsrf,
	}), nil
}

func NewSysAdminApiServiceImpl(db *bolt.DB, jwt *token.Service) openapi.SysAdminApiServicer {
	return &SysAdminApiServiceImpl{
		db:         db,
		jwtService: jwt,
	}
}
