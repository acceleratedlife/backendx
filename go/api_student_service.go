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
	"errors"
	"net/http"
)

// StudentApiService is a service that implements the logic for the StudentApiServicer
// This service should implement the business logic for every endpoint for the StudentApi API.
// Include any external packages or services that will be required by this service.
type StudentApiService struct {
}

// NewStudentApiService creates a default api service
func NewStudentApiService() StudentApiServicer {
	return &StudentApiService{}
}

// AuctionBid - auction bid
func (s *StudentApiService) AuctionBid(ctx context.Context, auctionsPlaceBidBody AuctionsPlaceBidBody) (ImplResponse, error) {
	// TODO - update AuctionBid with the required logic for this service method.
	// Add api_student_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []InlineResponse2006{}) or use other options such as http.Ok ...
	//return Response(200, []InlineResponse2006{}), nil

	//TODO: Uncomment the next line to return response Response(400, InlineResponse400{}) or use other options such as http.Ok ...
	//return Response(400, InlineResponse400{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("AuctionBid method not implemented")
}

// BuckConvert - When a student is converting between 2 bucks
func (s *StudentApiService) BuckConvert(ctx context.Context, userId string, transactionsConversionTransactionBody TransactionsConversionTransactionBody) (ImplResponse, error) {
	// TODO - update BuckConvert with the required logic for this service method.
	// Add api_student_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, InlineResponse2009{}) or use other options such as http.Ok ...
	//return Response(200, InlineResponse2009{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("BuckConvert method not implemented")
}

// CryptoConvert - When a student is converting between 2 uBucks and Cryptos
func (s *StudentApiService) CryptoConvert(ctx context.Context, userId string, transactionCryptoTransactionBody TransactionCryptoTransactionBody) (ImplResponse, error) {
	// TODO - update CryptoConvert with the required logic for this service method.
	// Add api_student_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, InlineResponse20010{}) or use other options such as http.Ok ...
	//return Response(200, InlineResponse20010{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("CryptoConvert method not implemented")
}

// SearchAuctionsStudent - searches auctions
func (s *StudentApiService) SearchAuctionsStudent(ctx context.Context, userId string) (ImplResponse, error) {
	// TODO - update SearchAuctionsStudent with the required logic for this service method.
	// Add api_student_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []InlineResponse2006{}) or use other options such as http.Ok ...
	//return Response(200, []InlineResponse2006{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("SearchAuctionsStudent method not implemented")
}

// SearchBuckTransaction - searches for buck transactions
func (s *StudentApiService) SearchBuckTransaction(ctx context.Context, id string) (ImplResponse, error) {
	// TODO - update SearchBuckTransaction with the required logic for this service method.
	// Add api_student_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, InlineResponse2007{}) or use other options such as http.Ok ...
	//return Response(200, InlineResponse2007{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("SearchBuckTransaction method not implemented")
}

// SearchCrypto - returns the given crypto price, how many are owned and how many ubucks the user has.
func (s *StudentApiService) SearchCrypto(ctx context.Context, userId string, name string) (ImplResponse, error) {
	// TODO - update SearchCrypto with the required logic for this service method.
	// Add api_student_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, InlineResponse2002{}) or use other options such as http.Ok ...
	//return Response(200, InlineResponse2002{}), nil

	//TODO: Uncomment the next line to return response Response(404, InlineResponse404{}) or use other options such as http.Ok ...
	//return Response(404, InlineResponse404{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("SearchCrypto method not implemented")
}

// SearchCryptoTransaction - searches for Crypto transactions
func (s *StudentApiService) SearchCryptoTransaction(ctx context.Context, id string) (ImplResponse, error) {
	// TODO - update SearchCryptoTransaction with the required logic for this service method.
	// Add api_student_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, InlineResponse2008{}) or use other options such as http.Ok ...
	//return Response(200, InlineResponse2008{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("SearchCryptoTransaction method not implemented")
}

// SearchStudentCrypto - returns all crypto accounts for specific user
func (s *StudentApiService) SearchStudentCrypto(ctx context.Context, userId string) (ImplResponse, error) {
	// TODO - update SearchStudentCrypto with the required logic for this service method.
	// Add api_student_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []Account{}) or use other options such as http.Ok ...
	//return Response(200, []Account{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("SearchStudentCrypto method not implemented")
}

// SearchStudentUbuck - searches accounts for UBuck for this student at this school
func (s *StudentApiService) SearchStudentUbuck(ctx context.Context, userId string) (ImplResponse, error) {
	// TODO - update SearchStudentUbuck with the required logic for this service method.
	// Add api_student_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, InlineResponse2001{}) or use other options such as http.Ok ...
	//return Response(200, InlineResponse2001{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("SearchStudentUbuck method not implemented")
}

// StudentAddClass - student adding self to class
func (s *StudentApiService) StudentAddClass(ctx context.Context, requestAddClass RequestAddClass) (ImplResponse, error) {
	// TODO - update StudentAddClass with the required logic for this service method.
	// Add api_student_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []ResponseMemberClass{}) or use other options such as http.Ok ...
	//return Response(200, []ResponseMemberClass{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("StudentAddClass method not implemented")
}
