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

type ResponseAuctionStudentWinnerId struct {

	Id string `json:"_id,omitempty"`

	FirstName string `json:"firstName,omitempty"`

	LastName string `json:"lastName,omitempty"`
}

// AssertResponseAuctionStudentWinnerIdRequired checks if the required fields are not zero-ed
func AssertResponseAuctionStudentWinnerIdRequired(obj ResponseAuctionStudentWinnerId) error {
	return nil
}

// AssertRecurseResponseAuctionStudentWinnerIdRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of ResponseAuctionStudentWinnerId (e.g. [][]ResponseAuctionStudentWinnerId), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseResponseAuctionStudentWinnerIdRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aResponseAuctionStudentWinnerId, ok := obj.(ResponseAuctionStudentWinnerId)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertResponseAuctionStudentWinnerIdRequired(aResponseAuctionStudentWinnerId)
	})
}
