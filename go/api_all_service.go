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

// AllApiService is a service that implements the logic for the AllApiServicer
// This service should implement the business logic for every endpoint for the AllApi API.
// Include any external packages or services that will be required by this service.
type AllApiService struct {
}

// NewAllApiService creates a default api service
func NewAllApiService() AllApiServicer {
	return &AllApiService{}
}

// AuthUser - return authenticated user details
func (s *AllApiService) AuthUser(ctx context.Context) (ImplResponse, error) {
	// TODO - update AuthUser with the required logic for this service method.
	// Add api_all_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, ResponseAuth2{}) or use other options such as http.Ok ...
	//return Response(200, ResponseAuth2{}), nil

	//TODO: Uncomment the next line to return response Response(404, ResponseAuth4{}) or use other options such as http.Ok ...
	//return Response(404, ResponseAuth4{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("AuthUser method not implemented")
}

// ConfirmEmail - confirm email
func (s *AllApiService) ConfirmEmail(ctx context.Context, token string) (ImplResponse, error) {
	// TODO - update ConfirmEmail with the required logic for this service method.
	// Add api_all_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, InlineResponse20017{}) or use other options such as http.Ok ...
	//return Response(200, InlineResponse20017{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("ConfirmEmail method not implemented")
}

// ExchangeRate - returns exchange rate between 2 buck accounts
func (s *AllApiService) ExchangeRate(ctx context.Context, sellCurrency string, buyCurrency string) (ImplResponse, error) {
	// TODO - update ExchangeRate with the required logic for this service method.
	// Add api_all_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []InlineResponse2004{}) or use other options such as http.Ok ...
	//return Response(200, []InlineResponse2004{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("ExchangeRate method not implemented")
}

// Login - logging in
func (s *AllApiService) Login(ctx context.Context, requestLogin RequestLogin) (ImplResponse, error) {
	// TODO - update Login with the required logic for this service method.
	// Add api_all_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, ResponseLogin2{}) or use other options such as http.Ok ...
	//return Response(200, ResponseLogin2{}), nil

	//TODO: Uncomment the next line to return response Response(401, ResponseLogin4{}) or use other options such as http.Ok ...
	//return Response(401, ResponseLogin4{}), nil

	//TODO: Uncomment the next line to return response Response(404, ResponseLogin4{}) or use other options such as http.Ok ...
	//return Response(404, ResponseLogin4{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("Login method not implemented")
}

// Logout - logout
func (s *AllApiService) Logout(ctx context.Context, id string) (ImplResponse, error) {
	// TODO - update Logout with the required logic for this service method.
	// Add api_all_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, ResponseRegister2{}) or use other options such as http.Ok ...
	//return Response(200, ResponseRegister2{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("Logout method not implemented")
}

// ResetPassword - reset password
func (s *AllApiService) ResetPassword(ctx context.Context, usersResetPasswordBody UsersResetPasswordBody) (ImplResponse, error) {
	// TODO - update ResetPassword with the required logic for this service method.
	// Add api_all_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, InlineResponse20017{}) or use other options such as http.Ok ...
	//return Response(200, InlineResponse20017{}), nil

	//TODO: Uncomment the next line to return response Response(401, InlineResponse401{}) or use other options such as http.Ok ...
	//return Response(401, InlineResponse401{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("ResetPassword method not implemented")
}

// SearchAccount - searches for a account
func (s *AllApiService) SearchAccount(ctx context.Context, id string) (ImplResponse, error) {
	// TODO - update SearchAccount with the required logic for this service method.
	// Add api_all_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []Account{}) or use other options such as http.Ok ...
	//return Response(200, []Account{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("SearchAccount method not implemented")
}

// SearchBucks - searches bucks
func (s *AllApiService) SearchBucks(ctx context.Context, id string) (ImplResponse, error) {
	// TODO - update SearchBucks with the required logic for this service method.
	// Add api_all_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, Buck{}) or use other options such as http.Ok ...
	//return Response(200, Buck{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("SearchBucks method not implemented")
}

// SearchClass - searches for a class
func (s *AllApiService) SearchClass(ctx context.Context, id string) (ImplResponse, error) {
	// TODO - update SearchClass with the required logic for this service method.
	// Add api_all_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, InlineResponse2007{}) or use other options such as http.Ok ...
	//return Response(200, InlineResponse2007{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("SearchClass method not implemented")
}

// SearchSchool - searches for a school
func (s *AllApiService) SearchSchool(ctx context.Context, id string) (ImplResponse, error) {
	// TODO - update SearchSchool with the required logic for this service method.
	// Add api_all_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []School{}) or use other options such as http.Ok ...
	//return Response(200, []School{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("SearchSchool method not implemented")
}

// SearchStudent - return one student
func (s *AllApiService) SearchStudent(ctx context.Context, id RequestUser) (ImplResponse, error) {
	// TODO - update SearchStudent with the required logic for this service method.
	// Add api_all_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, User{}) or use other options such as http.Ok ...
	//return Response(200, User{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("SearchStudent method not implemented")
}

// SearchStudentBuck - returns all buck accounts for specific user
func (s *AllApiService) SearchStudentBuck(ctx context.Context, userId string) (ImplResponse, error) {
	// TODO - update SearchStudentBuck with the required logic for this service method.
	// Add api_all_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []InlineResponse2003{}) or use other options such as http.Ok ...
	//return Response(200, []InlineResponse2003{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("SearchStudentBuck method not implemented")
}

// SearchStudents - return all students from a school
func (s *AllApiService) SearchStudents(ctx context.Context) (ImplResponse, error) {
	// TODO - update SearchStudents with the required logic for this service method.
	// Add api_all_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []UserNoHistory{}) or use other options such as http.Ok ...
	//return Response(200, []UserNoHistory{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("SearchStudents method not implemented")
}

// UserEdit - edit a user
func (s *AllApiService) UserEdit(ctx context.Context, usersUserBody UsersUserBody) (ImplResponse, error) {
	// TODO - update UserEdit with the required logic for this service method.
	// Add api_all_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, User{}) or use other options such as http.Ok ...
	//return Response(200, User{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("UserEdit method not implemented")
}
