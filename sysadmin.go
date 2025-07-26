package main

import (
	"context"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/go-pkgz/auth/token"
	"github.com/go-pkgz/lgr"
	bolt "go.etcd.io/bbolt"
)

type SysAdminApiServiceImpl struct {
	clock      Clock
	db         *bolt.DB
	jwtService *token.Service
}

func (s SysAdminApiServiceImpl) DeleteAccount(ctx context.Context, s2 string) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s SysAdminApiServiceImpl) DeleteSchool(ctx context.Context, s2 string) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(s.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	if userDetails.Role != UserRoleSysAdmin {
		return openapi.Response(401, ""), nil
	}

	err = deleteSchool(s.db, s2)

	if err != nil {
		return openapi.Response(500, nil), err
	}

	return openapi.Response(200, nil), nil
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

func (s SysAdminApiServiceImpl) MakeSchool(ctx context.Context, body openapi.RequestMakeSchool) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(s.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	if userDetails.Role != UserRoleSysAdmin {
		return openapi.Response(401, ""), nil
	}

	_, passwords := constSlice()
	password := randomWords(1, 10, passwords)
	response := openapi.ResponseResetPassword{
		Password: password,
	}

	err = createNewSchool(s.db, s.clock, body, response.Password)

	if err != nil {
		return openapi.Response(500, nil), err
	}

	return openapi.Response(200, response), nil
}

func (s SysAdminApiServiceImpl) MessageAll(ctx context.Context, body openapi.RequestMessage) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(s.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	if userDetails.Role != UserRoleSysAdmin {
		return openapi.Response(401, ""), nil
	}

	err = message(s.db, &body.Message, nil, nil, true, true)
	if err != nil {
		return openapi.Response(500, nil), err
	}

	return openapi.Response(200, nil), nil
}

func (s SysAdminApiServiceImpl) MessageAllSchool(ctx context.Context, body openapi.RequestMessage) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(s.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	if userDetails.Role != UserRoleSysAdmin {
		return openapi.Response(401, ""), nil
	}

	err = message(s.db, &body.Message, nil, &body.SchoolId, true, true)
	if err != nil {
		return openapi.Response(500, nil), err
	}

	return openapi.Response(200, nil), nil
}

func (s SysAdminApiServiceImpl) MessageAllStaff(ctx context.Context, body openapi.RequestMessage) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(s.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	if userDetails.Role != UserRoleSysAdmin {
		return openapi.Response(401, ""), nil
	}

	err = message(s.db, &body.Message, nil, nil, false, true)
	if err != nil {
		return openapi.Response(500, nil), err
	}

	return openapi.Response(200, nil), nil
}

func (s SysAdminApiServiceImpl) MessageAllStudents(ctx context.Context, body openapi.RequestMessage) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(s.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	if userDetails.Role != UserRoleSysAdmin {
		return openapi.Response(401, ""), nil
	}

	err = message(s.db, &body.Message, nil, nil, true, false)
	if err != nil {
		return openapi.Response(500, nil), err
	}

	return openapi.Response(200, nil), nil
}

func (s SysAdminApiServiceImpl) MessageAllSchoolStaff(ctx context.Context, body openapi.RequestMessage) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(s.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	if userDetails.Role != UserRoleSysAdmin {
		return openapi.Response(401, ""), nil
	}

	err = message(s.db, &body.Message, nil, &body.SchoolId, false, true)
	if err != nil {
		return openapi.Response(500, nil), err
	}

	return openapi.Response(200, nil), nil
}

func (s SysAdminApiServiceImpl) MessageAllSchoolStudents(ctx context.Context, body openapi.RequestMessage) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(s.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	if userDetails.Role != UserRoleSysAdmin {
		return openapi.Response(401, ""), nil
	}
	err = message(s.db, &body.Message, nil, &body.SchoolId, true, false)
	if err != nil {
		return openapi.Response(500, nil), err
	}

	return openapi.Response(200, nil), nil
}

func (s SysAdminApiServiceImpl) MessageUser(ctx context.Context, body openapi.RequestMessage) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(s.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	if userDetails.Role != UserRoleSysAdmin {
		return openapi.Response(401, ""), nil
	}

	err = message(s.db, &body.Message, &body.UserId, nil, false, false)
	if err != nil {
		return openapi.Response(500, nil), err
	}

	return openapi.Response(200, nil), nil
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

	lgr.Printf("INFO %s is impersonating %s", userDetails.Name, tgt.Name)

	return openapi.Response(200, openapi.ResponseImpersonate{
		Token: jwtStr,
		Xsrf:  xsrf,
	}), nil
}

func NewSysAdminApiServiceImpl(db *bolt.DB, clock Clock, jwt *token.Service) openapi.SysAdminApiServicer {
	return &SysAdminApiServiceImpl{
		clock:      clock,
		db:         db,
		jwtService: jwt,
	}
}
