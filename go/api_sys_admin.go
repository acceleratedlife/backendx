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

// SysAdminApiController binds http requests to an api service and writes the service results to the http response
type SysAdminApiController struct {
	service SysAdminApiServicer
	errorHandler ErrorHandler
}

// SysAdminApiOption for how the controller is set up.
type SysAdminApiOption func(*SysAdminApiController)

// WithSysAdminApiErrorHandler inject ErrorHandler into controller
func WithSysAdminApiErrorHandler(h ErrorHandler) SysAdminApiOption {
	return func(c *SysAdminApiController) {
		c.errorHandler = h
	}
}

// NewSysAdminApiController creates a default api controller
func NewSysAdminApiController(s SysAdminApiServicer, opts ...SysAdminApiOption) Router {
	controller := &SysAdminApiController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all of the api route for the SysAdminApiController
func (c *SysAdminApiController) Routes() Routes {
	return Routes{ 
		{
			"CreateBuck",
			strings.ToUpper("Post"),
			"/api/bucks/buck",
			c.CreateBuck,
		},
		{
			"DeleteAccount",
			strings.ToUpper("Delete"),
			"/api/accounts/account",
			c.DeleteAccount,
		},
		{
			"DeleteBuck",
			strings.ToUpper("Delete"),
			"/api/bucks/buck",
			c.DeleteBuck,
		},
		{
			"DeleteSchool",
			strings.ToUpper("Delete"),
			"/api/schools/school",
			c.DeleteSchool,
		},
		{
			"Deletetransaction",
			strings.ToUpper("Delete"),
			"/api/transactions/transaction",
			c.Deletetransaction,
		},
		{
			"EditAccount",
			strings.ToUpper("Put"),
			"/api/accounts/account",
			c.EditAccount,
		},
		{
			"EditBuck",
			strings.ToUpper("Put"),
			"/api/bucks/buck",
			c.EditBuck,
		},
		{
			"EditSchool",
			strings.ToUpper("Put"),
			"/api/schools/school",
			c.EditSchool,
		},
		{
			"MakeAccount",
			strings.ToUpper("Post"),
			"/api/accounts/account",
			c.MakeAccount,
		},
		{
			"MakeSchool",
			strings.ToUpper("Post"),
			"/api/schools/school",
			c.MakeSchool,
		},
		{
			"SearchSchools",
			strings.ToUpper("Get"),
			"/api/schools",
			c.SearchSchools,
		},
		{
			"SearchTransaction",
			strings.ToUpper("Get"),
			"/api/transactions/transaction",
			c.SearchTransaction,
		},
	}
}

// CreateBuck - create buck
func (c *SysAdminApiController) CreateBuck(w http.ResponseWriter, r *http.Request) {
	bucksBuckBody1Param := BucksBuckBody1{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&bucksBuckBody1Param); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertBucksBuckBody1Required(bucksBuckBody1Param); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.CreateBuck(r.Context(), bucksBuckBody1Param)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// DeleteAccount - delete school
func (c *SysAdminApiController) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	idParam := query.Get("_id")
	result, err := c.service.DeleteAccount(r.Context(), idParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// DeleteBuck - delete buck
func (c *SysAdminApiController) DeleteBuck(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	idParam := query.Get("_id")
	result, err := c.service.DeleteBuck(r.Context(), idParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// DeleteSchool - delete school
func (c *SysAdminApiController) DeleteSchool(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	idParam := query.Get("_id")
	result, err := c.service.DeleteSchool(r.Context(), idParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// Deletetransaction - delete transaction
func (c *SysAdminApiController) Deletetransaction(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	idParam := query.Get("_id")
	result, err := c.service.Deletetransaction(r.Context(), idParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// EditAccount - edit account
func (c *SysAdminApiController) EditAccount(w http.ResponseWriter, r *http.Request) {
	accountsAccountBodyParam := AccountsAccountBody{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&accountsAccountBodyParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertAccountsAccountBodyRequired(accountsAccountBodyParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.EditAccount(r.Context(), accountsAccountBodyParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// EditBuck - edit buck
func (c *SysAdminApiController) EditBuck(w http.ResponseWriter, r *http.Request) {
	bucksBuckBodyParam := BucksBuckBody{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&bucksBuckBodyParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertBucksBuckBodyRequired(bucksBuckBodyParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.EditBuck(r.Context(), bucksBuckBodyParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// EditSchool - edit school
func (c *SysAdminApiController) EditSchool(w http.ResponseWriter, r *http.Request) {
	schoolsSchoolBodyParam := SchoolsSchoolBody{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&schoolsSchoolBodyParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertSchoolsSchoolBodyRequired(schoolsSchoolBodyParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.EditSchool(r.Context(), schoolsSchoolBodyParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// MakeAccount - make account
func (c *SysAdminApiController) MakeAccount(w http.ResponseWriter, r *http.Request) {
	accountsAccountBody1Param := AccountsAccountBody1{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&accountsAccountBody1Param); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertAccountsAccountBody1Required(accountsAccountBody1Param); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.MakeAccount(r.Context(), accountsAccountBody1Param)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// MakeSchool - make a new school
func (c *SysAdminApiController) MakeSchool(w http.ResponseWriter, r *http.Request) {
	schoolsSchoolBody1Param := SchoolsSchoolBody1{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&schoolsSchoolBody1Param); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertSchoolsSchoolBody1Required(schoolsSchoolBody1Param); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.MakeSchool(r.Context(), schoolsSchoolBody1Param)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// SearchSchools - searches schools
func (c *SysAdminApiController) SearchSchools(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	zipParam, err := parseInt32Parameter(query.Get("zip"), true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	result, err := c.service.SearchSchools(r.Context(), zipParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// SearchTransaction - searches for a transaction
func (c *SysAdminApiController) SearchTransaction(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	idParam := query.Get("_id")
	result, err := c.service.SearchTransaction(r.Context(), idParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}
