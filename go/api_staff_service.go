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

// AuctionApprove - 
func (s *StaffApiService) AuctionApprove(ctx context.Context, requestAuctionAction RequestAuctionAction) (ImplResponse, error) {
	// TODO - update AuctionApprove with the required logic for this service method.
	// Add api_staff_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("AuctionApprove method not implemented")
}

// AuctionReject - 
func (s *StaffApiService) AuctionReject(ctx context.Context, id string) (ImplResponse, error) {
	// TODO - update AuctionReject with the required logic for this service method.
	// Add api_staff_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("AuctionReject method not implemented")
}

// AuctionsAll - get all auctions for school
func (s *StaffApiService) AuctionsAll(ctx context.Context) (ImplResponse, error) {
	// TODO - update AuctionsAll with the required logic for this service method.
	// Add api_staff_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []ResponseAuctionStudent{}) or use other options such as http.Ok ...
	//return Response(200, []ResponseAuctionStudent{}), nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("AuctionsAll method not implemented")
}

// DeleteStudent - delete student
func (s *StaffApiService) DeleteStudent(ctx context.Context, id string) (ImplResponse, error) {
	// TODO - update DeleteStudent with the required logic for this service method.
	// Add api_staff_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("DeleteStudent method not implemented")
}

// Deleteclass - delete class
func (s *StaffApiService) Deleteclass(ctx context.Context, id string) (ImplResponse, error) {
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
func (s *StaffApiService) ResetPassword(ctx context.Context, requestUser RequestUser) (ImplResponse, error) {
	// TODO - update ResetPassword with the required logic for this service method.
	// Add api_staff_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, ResponseResetPassword{}) or use other options such as http.Ok ...
	//return Response(200, ResponseResetPassword{}), nil

	//TODO: Uncomment the next line to return response Response(401, {}) or use other options such as http.Ok ...
	//return Response(401, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("ResetPassword method not implemented")
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

// SearchEvents - returns all of todays events for the users school
func (s *StaffApiService) SearchEvents(ctx context.Context) (ImplResponse, error) {
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
