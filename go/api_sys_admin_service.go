/*
 * AL
 *
 * This is a simple API
 *
 * API version: 1.0.1
 * Contact: you@your-company.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"context"
	"errors"
	"net/http"
)

// SysAdminApiService is a service that implements the logic for the SysAdminApiServicer
// This service should implement the business logic for every endpoint for the SysAdminApi API.
// Include any external packages or services that will be required by this service.
type SysAdminApiService struct {
}

// NewSysAdminApiService creates a default api service
func NewSysAdminApiService() SysAdminApiServicer {
	return &SysAdminApiService{}
}

// CreateBuck - create buck
func (s *SysAdminApiService) CreateBuck(ctx context.Context, bucksBuckBody1 BucksBuckBody1) (ImplResponse, error) {
	// TODO - update CreateBuck with the required logic for this service method.
	// Add api_sys_admin_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, Buck{}) or use other options such as http.Ok ...
	//return Response(200, Buck{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("CreateBuck method not implemented")
}

// DeleteAccount - delete school
func (s *SysAdminApiService) DeleteAccount(ctx context.Context, id string) (ImplResponse, error) {
	// TODO - update DeleteAccount with the required logic for this service method.
	// Add api_sys_admin_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("DeleteAccount method not implemented")
}

// DeleteBuck - delete buck
func (s *SysAdminApiService) DeleteBuck(ctx context.Context, id string) (ImplResponse, error) {
	// TODO - update DeleteBuck with the required logic for this service method.
	// Add api_sys_admin_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, Buck{}) or use other options such as http.Ok ...
	//return Response(200, Buck{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("DeleteBuck method not implemented")
}

// DeleteSchool - delete school
func (s *SysAdminApiService) DeleteSchool(ctx context.Context, id string) (ImplResponse, error) {
	// TODO - update DeleteSchool with the required logic for this service method.
	// Add api_sys_admin_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("DeleteSchool method not implemented")
}

// DeleteUser - delete user
func (s *SysAdminApiService) DeleteUser(ctx context.Context, email string) (ImplResponse, error) {
	// TODO - update DeleteUser with the required logic for this service method.
	// Add api_sys_admin_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, InlineResponse20011{}) or use other options such as http.Ok ...
	//return Response(200, InlineResponse20011{}), nil

	//TODO: Uncomment the next line to return response Response(403, InlineResponse403{}) or use other options such as http.Ok ...
	//return Response(403, InlineResponse403{}), nil

	//TODO: Uncomment the next line to return response Response(404, InlineResponse4041{}) or use other options such as http.Ok ...
	//return Response(404, InlineResponse4041{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("DeleteUser method not implemented")
}

// Deletetransaction - delete transaction
func (s *SysAdminApiService) Deletetransaction(ctx context.Context, id string) (ImplResponse, error) {
	// TODO - update Deletetransaction with the required logic for this service method.
	// Add api_sys_admin_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, string{}) or use other options such as http.Ok ...
	//return Response(200, string{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("Deletetransaction method not implemented")
}

// EditAccount - edit account
func (s *SysAdminApiService) EditAccount(ctx context.Context, accountsAccountBody AccountsAccountBody) (ImplResponse, error) {
	// TODO - update EditAccount with the required logic for this service method.
	// Add api_sys_admin_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []Account{}) or use other options such as http.Ok ...
	//return Response(200, []Account{}), nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("EditAccount method not implemented")
}

// EditBuck - edit buck
func (s *SysAdminApiService) EditBuck(ctx context.Context, bucksBuckBody BucksBuckBody) (ImplResponse, error) {
	// TODO - update EditBuck with the required logic for this service method.
	// Add api_sys_admin_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, Buck{}) or use other options such as http.Ok ...
	//return Response(200, Buck{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("EditBuck method not implemented")
}

// EditSchool - edit school
func (s *SysAdminApiService) EditSchool(ctx context.Context, schoolsSchoolBody SchoolsSchoolBody) (ImplResponse, error) {
	// TODO - update EditSchool with the required logic for this service method.
	// Add api_sys_admin_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []School{}) or use other options such as http.Ok ...
	//return Response(200, []School{}), nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("EditSchool method not implemented")
}

// MakeAccount - make account
func (s *SysAdminApiService) MakeAccount(ctx context.Context, accountsAccountBody1 AccountsAccountBody1) (ImplResponse, error) {
	// TODO - update MakeAccount with the required logic for this service method.
	// Add api_sys_admin_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []Account{}) or use other options such as http.Ok ...
	//return Response(200, []Account{}), nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("MakeAccount method not implemented")
}

// MakeSchool - make a new school
func (s *SysAdminApiService) MakeSchool(ctx context.Context, schoolsSchoolBody1 SchoolsSchoolBody1) (ImplResponse, error) {
	// TODO - update MakeSchool with the required logic for this service method.
	// Add api_sys_admin_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, InlineResponse200{}) or use other options such as http.Ok ...
	//return Response(200, InlineResponse200{}), nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("MakeSchool method not implemented")
}

// SearchSchools - searches schools
func (s *SysAdminApiService) SearchSchools(ctx context.Context, zip int32) (ImplResponse, error) {
	// TODO - update SearchSchools with the required logic for this service method.
	// Add api_sys_admin_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []ResponseSchoolsInner{}) or use other options such as http.Ok ...
	//return Response(200, []ResponseSchoolsInner{}), nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("SearchSchools method not implemented")
}

// SearchTransaction - searches for a transaction
func (s *SysAdminApiService) SearchTransaction(ctx context.Context, id string) (ImplResponse, error) {
	// TODO - update SearchTransaction with the required logic for this service method.
	// Add api_sys_admin_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, Transaction{}) or use other options such as http.Ok ...
	//return Response(200, Transaction{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("SearchTransaction method not implemented")
}
