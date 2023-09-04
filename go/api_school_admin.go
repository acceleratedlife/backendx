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
	"net/http"
	"strings"
)

// SchoolAdminApiController binds http requests to an api service and writes the service results to the http response
type SchoolAdminApiController struct {
	service      SchoolAdminApiServicer
	errorHandler ErrorHandler
}

// SchoolAdminApiOption for how the controller is set up.
type SchoolAdminApiOption func(*SchoolAdminApiController)

// WithSchoolAdminApiErrorHandler inject ErrorHandler into controller
func WithSchoolAdminApiErrorHandler(h ErrorHandler) SchoolAdminApiOption {
	return func(c *SchoolAdminApiController) {
		c.errorHandler = h
	}
}

// NewSchoolAdminApiController creates a default api controller
func NewSchoolAdminApiController(s SchoolAdminApiServicer, opts ...SchoolAdminApiOption) Router {
	controller := &SchoolAdminApiController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all of the api route for the SchoolAdminApiController
func (c *SchoolAdminApiController) Routes() Routes {
	return Routes{
		{
			"GetStudentCount",
			strings.ToUpper("Get"),
			"/api/schools/school/count",
			c.GetStudentCount,
		},
		{
			"SearchAdminTeacherClass",
			strings.ToUpper("Get"),
			"/api/classes/teachers",
			c.SearchAdminTeacherClass,
		},
	}
}

// GetStudentCount - gets student count for a school
func (c *SchoolAdminApiController) GetStudentCount(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	schoolIdParam := query.Get("schoolId")
	result, err := c.service.GetStudentCount(r.Context(), schoolIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// SearchAdminTeacherClass - gets the teacher class of an admin and all the teacher that are its members
func (c *SchoolAdminApiController) SearchAdminTeacherClass(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	idParam := query.Get("_id")
	result, err := c.service.SearchAdminTeacherClass(r.Context(), idParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}
