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

// StaffApiController binds http requests to an api service and writes the service results to the http response
type StaffApiController struct {
	service      StaffApiServicer
	errorHandler ErrorHandler
}

// StaffApiOption for how the controller is set up.
type StaffApiOption func(*StaffApiController)

// WithStaffApiErrorHandler inject ErrorHandler into controller
func WithStaffApiErrorHandler(h ErrorHandler) StaffApiOption {
	return func(c *StaffApiController) {
		c.errorHandler = h
	}
}

// NewStaffApiController creates a default api controller
func NewStaffApiController(s StaffApiServicer, opts ...StaffApiOption) Router {
	controller := &StaffApiController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all of the api route for the StaffApiController
func (c *StaffApiController) Routes() Routes {
	return Routes{
		{
			"DeleteAuction",
			strings.ToUpper("Delete"),
			"/api/auctions/auction",
			c.DeleteAuction,
		},
		{
			"DeleteStudent",
			strings.ToUpper("Delete"),
			"/api/users/user",
			c.DeleteStudent,
		},
		{
			"Deleteclass",
			strings.ToUpper("Delete"),
			"/api/classes/class",
			c.Deleteclass,
		},
		{
			"EditClass",
			strings.ToUpper("Put"),
			"/api/classes/class",
			c.EditClass,
		},
		{
			"KickClass",
			strings.ToUpper("Put"),
			"/api/classes/class/kick",
			c.KickClass,
		},
		{
			"MakeAuction",
			strings.ToUpper("Post"),
			"/api/auctions",
			c.MakeAuction,
		},
		{
			"MakeClass",
			strings.ToUpper("Post"),
			"/api/classes/class",
			c.MakeClass,
		},
		{
			"PayTransaction",
			strings.ToUpper("Post"),
			"/api/transactions/payTransaction",
			c.PayTransaction,
		},
		{
			"PayTransactions",
			strings.ToUpper("Post"),
			"/api/transactions/payTransactions",
			c.PayTransactions,
		},
		{
			"ResetPassword",
			strings.ToUpper("Post"),
			"/api/users/resetPassword",
			c.ResetPassword,
		},
		{
			"SearchAuctionsTeacher",
			strings.ToUpper("Get"),
			"/api/auctions",
			c.SearchAuctionsTeacher,
		},
		{
			"SearchClasses",
			strings.ToUpper("Get"),
			"/api/classes",
			c.SearchClasses,
		},
		{
			"SearchEvents",
			strings.ToUpper("Get"),
			"/api/events",
			c.SearchEvents,
		},
		{
			"SearchTransactions",
			strings.ToUpper("Get"),
			"/api/transactions",
			c.SearchTransactions,
		},
	}
}

// DeleteAuction - delete auction
func (c *StaffApiController) DeleteAuction(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	idParam := query.Get("_id")
	result, err := c.service.DeleteAuction(r.Context(), idParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// DeleteStudent - delete student
func (c *StaffApiController) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	idParam := query.Get("_id")
	result, err := c.service.DeleteStudent(r.Context(), idParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// Deleteclass - delete class
func (c *StaffApiController) Deleteclass(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	idParam := query.Get("_id")
	result, err := c.service.Deleteclass(r.Context(), idParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// EditClass - edit class
func (c *StaffApiController) EditClass(w http.ResponseWriter, r *http.Request) {
	requestEditClassParam := RequestEditClass{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&requestEditClassParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertRequestEditClassRequired(requestEditClassParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.EditClass(r.Context(), requestEditClassParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// KickClass - Kick member from class
func (c *StaffApiController) KickClass(w http.ResponseWriter, r *http.Request) {
	requestKickClassParam := RequestKickClass{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&requestKickClassParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertRequestKickClassRequired(requestKickClassParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.KickClass(r.Context(), requestKickClassParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// MakeAuction - make a new auction
func (c *StaffApiController) MakeAuction(w http.ResponseWriter, r *http.Request) {
	requestMakeAuctionParam := RequestMakeAuction{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&requestMakeAuctionParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertRequestMakeAuctionRequired(requestMakeAuctionParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.MakeAuction(r.Context(), requestMakeAuctionParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// MakeClass - make a new class
func (c *StaffApiController) MakeClass(w http.ResponseWriter, r *http.Request) {
	requestMakeClassParam := RequestMakeClass{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&requestMakeClassParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertRequestMakeClassRequired(requestMakeClassParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.MakeClass(r.Context(), requestMakeClassParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// PayTransaction - When a teacher or admin is paying/debting a student with their own bucks
func (c *StaffApiController) PayTransaction(w http.ResponseWriter, r *http.Request) {
	requestPayTransactionParam := RequestPayTransaction{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&requestPayTransactionParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertRequestPayTransactionRequired(requestPayTransactionParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.PayTransaction(r.Context(), requestPayTransactionParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// PayTransactions - When a teacher or admin is paying/debting an entire class
func (c *StaffApiController) PayTransactions(w http.ResponseWriter, r *http.Request) {
	requestPayTransactionsParam := RequestPayTransactions{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&requestPayTransactionsParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertRequestPayTransactionsRequired(requestPayTransactionsParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.PayTransactions(r.Context(), requestPayTransactionsParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// ResetPassword - reset password
func (c *StaffApiController) ResetPassword(w http.ResponseWriter, r *http.Request) {
	requestUserParam := RequestUser{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&requestUserParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertRequestUserRequired(requestUserParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.ResetPassword(r.Context(), requestUserParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// SearchAuctionsTeacher - searches auctions
func (c *StaffApiController) SearchAuctionsTeacher(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.SearchAuctionsTeacher(r.Context())
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// SearchClasses - searches for users classes
func (c *StaffApiController) SearchClasses(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	idParam := query.Get("_id")
	result, err := c.service.SearchClasses(r.Context(), idParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// SearchEvents - returns all of todays events for the users school
func (c *StaffApiController) SearchEvents(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.SearchEvents(r.Context())
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// SearchTransactions - all transactions a teacher is the owner of. Limit 30 sort by createdAt descending
func (c *StaffApiController) SearchTransactions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	idParam := query.Get("_id")
	result, err := c.service.SearchTransactions(r.Context(), idParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}
