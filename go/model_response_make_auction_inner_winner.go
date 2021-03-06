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

type ResponseMakeAuctionInnerWinner struct {

	FirstName string `json:"firstName,omitempty"`

	LastName string `json:"lastName,omitempty"`
}

// AssertResponseMakeAuctionInnerWinnerRequired checks if the required fields are not zero-ed
func AssertResponseMakeAuctionInnerWinnerRequired(obj ResponseMakeAuctionInnerWinner) error {
	return nil
}

// AssertRecurseResponseMakeAuctionInnerWinnerRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of ResponseMakeAuctionInnerWinner (e.g. [][]ResponseMakeAuctionInnerWinner), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseResponseMakeAuctionInnerWinnerRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aResponseMakeAuctionInnerWinner, ok := obj.(ResponseMakeAuctionInnerWinner)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertResponseMakeAuctionInnerWinnerRequired(aResponseMakeAuctionInnerWinner)
	})
}
