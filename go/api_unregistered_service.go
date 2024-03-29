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

// UnregisteredApiService is a service that implements the logic for the UnregisteredApiServicer
// This service should implement the business logic for every endpoint for the UnregisteredApi API.
// Include any external packages or services that will be required by this service.
type UnregisteredApiService struct {
}

// NewUnregisteredApiService creates a default api service
func NewUnregisteredApiService() UnregisteredApiServicer {
	return &UnregisteredApiService{}
}

// GetCryptos - returns all cryptos current values
func (s *UnregisteredApiService) GetCryptos(ctx context.Context) (ImplResponse, error) {
	// TODO - update GetCryptos with the required logic for this service method.
	// Add api_unregistered_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []ResponseCryptoPrice{}) or use other options such as http.Ok ...
	//return Response(200, []ResponseCryptoPrice{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("GetCryptos method not implemented")
}

// Register - Users register
func (s *UnregisteredApiService) Register(ctx context.Context, requestRegister RequestRegister) (ImplResponse, error) {
	// TODO - update Register with the required logic for this service method.
	// Add api_unregistered_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, ResponseRegister2{}) or use other options such as http.Ok ...
	//return Response(200, ResponseRegister2{}), nil

	//TODO: Uncomment the next line to return response Response(404, ResponseRegister4{}) or use other options such as http.Ok ...
	//return Response(404, ResponseRegister4{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("Register method not implemented")
}

// ResetStaffPassword - reset staff password
func (s *UnregisteredApiService) ResetStaffPassword(ctx context.Context, requestUser RequestUser) (ImplResponse, error) {
	// TODO - update ResetStaffPassword with the required logic for this service method.
	// Add api_unregistered_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(401, {}) or use other options such as http.Ok ...
	//return Response(401, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("ResetStaffPassword method not implemented")
}
