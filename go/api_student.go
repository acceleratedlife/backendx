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

// StudentApiController binds http requests to an api service and writes the service results to the http response
type StudentApiController struct {
	service      StudentApiServicer
	errorHandler ErrorHandler
}

// StudentApiOption for how the controller is set up.
type StudentApiOption func(*StudentApiController)

// WithStudentApiErrorHandler inject ErrorHandler into controller
func WithStudentApiErrorHandler(h ErrorHandler) StudentApiOption {
	return func(c *StudentApiController) {
		c.errorHandler = h
	}
}

// NewStudentApiController creates a default api controller
func NewStudentApiController(s StudentApiServicer, opts ...StudentApiOption) Router {
	controller := &StudentApiController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all of the api route for the StudentApiController
func (c *StudentApiController) Routes() Routes {
	return Routes{
		{
			"AuctionBid",
			strings.ToUpper("Put"),
			"/api/auctions/placeBid",
			c.AuctionBid,
		},
		{
			"BuckConvert",
			strings.ToUpper("Post"),
			"/api/transactions/conversionTransaction",
			c.BuckConvert,
		},
		{
			"CryptoConvert",
			strings.ToUpper("Post"),
			"/api/transaction/cryptoTransaction",
			c.CryptoConvert,
		},
		{
			"SearchAuctionsStudent",
			strings.ToUpper("Get"),
			"/api/auctions/student",
			c.SearchAuctionsStudent,
		},
		{
			"SearchBuckTransaction",
			strings.ToUpper("Get"),
			"/api/transactions/buckTransactions",
			c.SearchBuckTransaction,
		},
		{
			"SearchCrypto",
			strings.ToUpper("Get"),
			"/api/accounts/crypto",
			c.SearchCrypto,
		},
		{
			"SearchCryptoTransaction",
			strings.ToUpper("Get"),
			"/api/transactions/cryptoTransactions",
			c.SearchCryptoTransaction,
		},
		{
			"SearchStudentCrypto",
			strings.ToUpper("Get"),
			"/api/accounts/allCrypto",
			c.SearchStudentCrypto,
		},
		{
			"SearchStudentUbuck",
			strings.ToUpper("Get"),
			"/api/accounts/account/student",
			c.SearchStudentUbuck,
		},
		{
			"StudentAddClass",
			strings.ToUpper("Put"),
			"/api/classes/addClass",
			c.StudentAddClass,
		},
	}
}

// AuctionBid - auction bid
func (c *StudentApiController) AuctionBid(w http.ResponseWriter, r *http.Request) {
	requestAuctionBidParam := RequestAuctionBid{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&requestAuctionBidParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertRequestAuctionBidRequired(requestAuctionBidParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.AuctionBid(r.Context(), requestAuctionBidParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// BuckConvert - When a student is converting between 2 bucks
func (c *StudentApiController) BuckConvert(w http.ResponseWriter, r *http.Request) {
	requestBuckConvertParam := RequestBuckConvert{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&requestBuckConvertParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertRequestBuckConvertRequired(requestBuckConvertParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.BuckConvert(r.Context(), requestBuckConvertParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// CryptoConvert - When a student is converting between 2 uBucks and Cryptos
func (c *StudentApiController) CryptoConvert(w http.ResponseWriter, r *http.Request) {
	userIdParam := r.Header.Get("user._id")
	transactionCryptoTransactionBodyParam := TransactionCryptoTransactionBody{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&transactionCryptoTransactionBodyParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertTransactionCryptoTransactionBodyRequired(transactionCryptoTransactionBodyParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.CryptoConvert(r.Context(), userIdParam, transactionCryptoTransactionBodyParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// SearchAuctionsStudent - searches auctions
func (c *StudentApiController) SearchAuctionsStudent(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.SearchAuctionsStudent(r.Context())
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// SearchBuckTransaction - searches for buck transactions
func (c *StudentApiController) SearchBuckTransaction(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.SearchBuckTransaction(r.Context())
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// SearchCrypto - returns the given crypto price, how many are owned and how many ubucks the user has.
func (c *StudentApiController) SearchCrypto(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	nameParam := query.Get("name")
	result, err := c.service.SearchCrypto(r.Context(), nameParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// SearchCryptoTransaction - searches for Crypto transactions
func (c *StudentApiController) SearchCryptoTransaction(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	idParam := query.Get("_id")
	result, err := c.service.SearchCryptoTransaction(r.Context(), idParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// SearchStudentCrypto - returns all crypto accounts for specific user
func (c *StudentApiController) SearchStudentCrypto(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.SearchStudentCrypto(r.Context())
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// SearchStudentUbuck - searches accounts for UBuck for this student at this school
func (c *StudentApiController) SearchStudentUbuck(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.SearchStudentUbuck(r.Context())
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// StudentAddClass - student adding self to class
func (c *StudentApiController) StudentAddClass(w http.ResponseWriter, r *http.Request) {
	requestAddClassParam := RequestAddClass{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&requestAddClassParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertRequestAddClassRequired(requestAddClassParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.StudentAddClass(r.Context(), requestAddClassParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}
