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
	"encoding/json"
	"net/http"
	"strings"
)

// UnregisteredApiController binds http requests to an api service and writes the service results to the http response
type UnregisteredApiController struct {
	service      UnregisteredApiServicer
	errorHandler ErrorHandler
}

// UnregisteredApiOption for how the controller is set up.
type UnregisteredApiOption func(*UnregisteredApiController)

// WithUnregisteredApiErrorHandler inject ErrorHandler into controller
func WithUnregisteredApiErrorHandler(h ErrorHandler) UnregisteredApiOption {
	return func(c *UnregisteredApiController) {
		c.errorHandler = h
	}
}

// NewUnregisteredApiController creates a default api controller
func NewUnregisteredApiController(s UnregisteredApiServicer, opts ...UnregisteredApiOption) Router {
	controller := &UnregisteredApiController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all of the api route for the UnregisteredApiController
func (c *UnregisteredApiController) Routes() Routes {
	return Routes{
		{
			"Register",
			strings.ToUpper("Post"),
			"/api/users/register",
			c.Register,
		},
	}
}

// Register - Users register
func (c *UnregisteredApiController) Register(w http.ResponseWriter, r *http.Request) {
	requestRegisterParam := RequestRegister{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&requestRegisterParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertRequestRegisterRequired(requestRegisterParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.Register(r.Context(), requestRegisterParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}
