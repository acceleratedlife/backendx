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

type ResponseAuctionStudentInnerWinner struct {

	Id string `json:"_id,omitempty"`

	FirstName string `json:"firstName,omitempty"`

	LastName string `json:"lastName,omitempty"`
}

// AssertResponseAuctionStudentInnerWinnerRequired checks if the required fields are not zero-ed
func AssertResponseAuctionStudentInnerWinnerRequired(obj ResponseAuctionStudentInnerWinner) error {
	return nil
}

// AssertRecurseResponseAuctionStudentInnerWinnerRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of ResponseAuctionStudentInnerWinner (e.g. [][]ResponseAuctionStudentInnerWinner), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseResponseAuctionStudentInnerWinnerRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aResponseAuctionStudentInnerWinner, ok := obj.(ResponseAuctionStudentInnerWinner)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertResponseAuctionStudentInnerWinnerRequired(aResponseAuctionStudentInnerWinner)
	})
}
