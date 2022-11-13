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
	"net/http"
	"errors"
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

	//TODO: Uncomment the next line to return response Response(200, InlineResponse2002{}) or use other options such as http.Ok ...
	//return Response(200, InlineResponse2002{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("ConfirmEmail method not implemented")
}

// DeleteAuction - delete auction
func (s *AllApiService) DeleteAuction(ctx context.Context, id string) (ImplResponse, error) {
	// TODO - update DeleteAuction with the required logic for this service method.
	// Add api_all_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("DeleteAuction method not implemented")
}

// ExchangeRate - returns exchange rate between 2 buck accounts
func (s *AllApiService) ExchangeRate(ctx context.Context, sellCurrency string, buyCurrency string) (ImplResponse, error) {
	// TODO - update ExchangeRate with the required logic for this service method.
	// Add api_all_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []ResponseCurrencyExchange{}) or use other options such as http.Ok ...
	//return Response(200, []ResponseCurrencyExchange{}), nil

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

// MakeAuction - make a new auction
func (s *AllApiService) MakeAuction(ctx context.Context, requestMakeAuction RequestMakeAuction) (ImplResponse, error) {
	// TODO - update MakeAuction with the required logic for this service method.
	// Add api_all_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("MakeAuction method not implemented")
}

// PayTransaction - When a teacher or admin is paying/debting a student with their own bucks
func (s *AllApiService) PayTransaction(ctx context.Context, requestPayTransaction RequestPayTransaction) (ImplResponse, error) {
	// TODO - update PayTransaction with the required logic for this service method.
	// Add api_all_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("PayTransaction method not implemented")
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

// SearchAllBucks - searches all bucks
func (s *AllApiService) SearchAllBucks(ctx context.Context) (ImplResponse, error) {
	// TODO - update SearchAllBucks with the required logic for this service method.
	// Add api_all_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []Buck{}) or use other options such as http.Ok ...
	//return Response(200, []Buck{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("SearchAllBucks method not implemented")
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

	//TODO: Uncomment the next line to return response Response(200, ClassWithMembers{}) or use other options such as http.Ok ...
	//return Response(200, ClassWithMembers{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("SearchClass method not implemented")
}

// SearchClasses - searches for users classes
func (s *AllApiService) SearchClasses(ctx context.Context) (ImplResponse, error) {
	// TODO - update SearchClasses with the required logic for this service method.
	// Add api_all_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []Class{}) or use other options such as http.Ok ...
	//return Response(200, []Class{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("SearchClasses method not implemented")
}

// SearchMarketItems - all market items relitive to this user
func (s *AllApiService) SearchMarketItems(ctx context.Context, id string) (ImplResponse, error) {
	// TODO - update SearchMarketItems with the required logic for this service method.
	// Add api_all_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []ResponseMarketItem{}) or use other options such as http.Ok ...
	//return Response(200, []ResponseMarketItem{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("SearchMarketItems method not implemented")
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
func (s *AllApiService) SearchStudent(ctx context.Context, id string) (ImplResponse, error) {
	// TODO - update SearchStudent with the required logic for this service method.
	// Add api_all_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, User{}) or use other options such as http.Ok ...
	//return Response(200, User{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("SearchStudent method not implemented")
}

// SearchStudentBucks - returns all buck accounts for specific user
func (s *AllApiService) SearchStudentBucks(ctx context.Context) (ImplResponse, error) {
	// TODO - update SearchStudentBucks with the required logic for this service method.
	// Add api_all_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []ResponseCurrencyExchange{}) or use other options such as http.Ok ...
	//return Response(200, []ResponseCurrencyExchange{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("SearchStudentBucks method not implemented")
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

// SearchTeachers - all the teachers that are at the same school of the user
func (s *AllApiService) SearchTeachers(ctx context.Context) (ImplResponse, error) {
	// TODO - update SearchTeachers with the required logic for this service method.
	// Add api_all_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []ResponseTeachers{}) or use other options such as http.Ok ...
	//return Response(200, []ResponseTeachers{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("SearchTeachers method not implemented")
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
