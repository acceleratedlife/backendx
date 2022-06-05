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

// StaffApiService is a service that implements the logic for the StaffApiServicer
// This service should implement the business logic for every endpoint for the StaffApi API.
// Include any external packages or services that will be required by this service.
type StaffApiService struct {
}

// NewStaffApiService creates a default api service
func NewStaffApiService() StaffApiServicer {
	return &StaffApiService{}
}

// DeleteAuction - delete auction
func (s *StaffApiService) DeleteAuction(ctx context.Context, id RequestUser) (ImplResponse, error) {
	// TODO - update DeleteAuction with the required logic for this service method.
	// Add api_staff_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []Auction{}) or use other options such as http.Ok ...
	//return Response(200, []Auction{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("DeleteAuction method not implemented")
}

// DeleteUser - delete user
func (s *StaffApiService) DeleteUser(ctx context.Context, email string) (ImplResponse, error) {
	// TODO - update DeleteUser with the required logic for this service method.
	// Add api_staff_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("DeleteUser method not implemented")
}

// Deleteclass - delete class
func (s *StaffApiService) Deleteclass(ctx context.Context, id RequestUser) (ImplResponse, error) {
	// TODO - update Deleteclass with the required logic for this service method.
	// Add api_staff_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("Deleteclass method not implemented")
}

// EditClass - edit class
func (s *StaffApiService) EditClass(ctx context.Context, requestEditClass RequestEditClass) (ImplResponse, error) {
	// TODO - update EditClass with the required logic for this service method.
	// Add api_staff_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, Class{}) or use other options such as http.Ok ...
	//return Response(200, Class{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("EditClass method not implemented")
}

// KickClass - Kick member from class
func (s *StaffApiService) KickClass(ctx context.Context, requestKickClass RequestKickClass) (ImplResponse, error) {
	// TODO - update KickClass with the required logic for this service method.
	// Add api_staff_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("KickClass method not implemented")
}

// MakeAuction - make a new auction
func (s *StaffApiService) MakeAuction(ctx context.Context, requestMakeAuction RequestMakeAuction) (ImplResponse, error) {
	// TODO - update MakeAuction with the required logic for this service method.
	// Add api_staff_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []Auction{}) or use other options such as http.Ok ...
	//return Response(200, []Auction{}), nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("MakeAuction method not implemented")
}

// MakeClass - make a new class
func (s *StaffApiService) MakeClass(ctx context.Context, requestMakeClass RequestMakeClass) (ImplResponse, error) {
	// TODO - update MakeClass with the required logic for this service method.
	// Add api_staff_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []Class{}) or use other options such as http.Ok ...
	//return Response(200, []Class{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("MakeClass method not implemented")
}

// PayTransaction - When a teacher or admin is paying/debting a student with their own bucks
func (s *StaffApiService) PayTransaction(ctx context.Context, requestPayTransaction RequestPayTransaction) (ImplResponse, error) {
	// TODO - update PayTransaction with the required logic for this service method.
	// Add api_staff_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("PayTransaction method not implemented")
}

// PayTransactions - When a teacher or admin is paying/debting an entire class
func (s *StaffApiService) PayTransactions(ctx context.Context, requestPayTransactions RequestPayTransactions) (ImplResponse, error) {
	// TODO - update PayTransactions with the required logic for this service method.
	// Add api_staff_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(400, []ResponsePayTransactions{}) or use other options such as http.Ok ...
	//return Response(400, []ResponsePayTransactions{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("PayTransactions method not implemented")
}

// ResetPassword - reset password
func (s *StaffApiService) ResetPassword(ctx context.Context, usersResetPasswordBody UsersResetPasswordBody) (ImplResponse, error) {
	// TODO - update ResetPassword with the required logic for this service method.
	// Add api_staff_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, InlineResponse2006{}) or use other options such as http.Ok ...
	//return Response(200, InlineResponse2006{}), nil

	//TODO: Uncomment the next line to return response Response(401, InlineResponse401{}) or use other options such as http.Ok ...
	//return Response(401, InlineResponse401{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("ResetPassword method not implemented")
}

// SearchAllBucks - searches all bucks
func (s *StaffApiService) SearchAllBucks(ctx context.Context, id string) (ImplResponse, error) {
	// TODO - update SearchAllBucks with the required logic for this service method.
	// Add api_staff_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []Buck{}) or use other options such as http.Ok ...
	//return Response(200, []Buck{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("SearchAllBucks method not implemented")
}

// SearchAuctionsTeacher - searches auctions
func (s *StaffApiService) SearchAuctionsTeacher(ctx context.Context) (ImplResponse, error) {
	// TODO - update SearchAuctionsTeacher with the required logic for this service method.
	// Add api_staff_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []Auction{}) or use other options such as http.Ok ...
	//return Response(200, []Auction{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("SearchAuctionsTeacher method not implemented")
}

// SearchClasses - searches for users classes
func (s *StaffApiService) SearchClasses(ctx context.Context, id RequestUser) (ImplResponse, error) {
	// TODO - update SearchClasses with the required logic for this service method.
	// Add api_staff_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []Class{}) or use other options such as http.Ok ...
	//return Response(200, []Class{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("SearchClasses method not implemented")
}

// SearchEvents - returns all of todays events for the users school
func (s *StaffApiService) SearchEvents(ctx context.Context, id string) (ImplResponse, error) {
	// TODO - update SearchEvents with the required logic for this service method.
	// Add api_staff_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []ResponseEvents{}) or use other options such as http.Ok ...
	//return Response(200, []ResponseEvents{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("SearchEvents method not implemented")
}

// SearchTransactions - all transactions a teacher is the owner of. Limit 30 sort by createdAt descending
func (s *StaffApiService) SearchTransactions(ctx context.Context, id string) (ImplResponse, error) {
	// TODO - update SearchTransactions with the required logic for this service method.
	// Add api_staff_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []ResponseTransactions{}) or use other options such as http.Ok ...
	//return Response(200, []ResponseTransactions{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("SearchTransactions method not implemented")
}
