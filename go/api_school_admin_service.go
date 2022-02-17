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

// SchoolAdminApiService is a service that implements the logic for the SchoolAdminApiServicer
// This service should implement the business logic for every endpoint for the SchoolAdminApi API.
// Include any external packages or services that will be required by this service.
type SchoolAdminApiService struct {
}

// NewSchoolAdminApiService creates a default api service
func NewSchoolAdminApiService() SchoolAdminApiServicer {
	return &SchoolAdminApiService{}
}

// SearchAdminTeacherClass - gets the teacher class of an admin and all the teacher that are its members
func (s *SchoolAdminApiService) SearchAdminTeacherClass(ctx context.Context, id string) (ImplResponse, error) {
	// TODO - update SearchAdminTeacherClass with the required logic for this service method.
	// Add api_school_admin_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, ClassWithMembers{}) or use other options such as http.Ok ...
	//return Response(200, ClassWithMembers{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("SearchAdminTeacherClass method not implemented")
}
