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
	"time"
)

type ResponseAuctionStudent struct {

	Id string `json:"_id,omitempty"`

	Bid float32 `json:"bid,omitempty"`

	Description string `json:"description,omitempty"`

	EndDate time.Time `json:"endDate,omitempty"`

	StartDate time.Time `json:"startDate,omitempty"`

	ItemNumber string `json:"itemNumber,omitempty"`

	Owner ResponseAuctionStudentOwner `json:"owner,omitempty"`

	Winner ResponseAuctionStudentWinner `json:"winner,omitempty"`
}

// AssertResponseAuctionStudentRequired checks if the required fields are not zero-ed
func AssertResponseAuctionStudentRequired(obj ResponseAuctionStudent) error {
	if err := AssertResponseAuctionStudentOwnerRequired(obj.Owner); err != nil {
		return err
	}
	if err := AssertResponseAuctionStudentWinnerRequired(obj.Winner); err != nil {
		return err
	}
	return nil
}

// AssertRecurseResponseAuctionStudentRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of ResponseAuctionStudent (e.g. [][]ResponseAuctionStudent), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseResponseAuctionStudentRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aResponseAuctionStudent, ok := obj.(ResponseAuctionStudent)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertResponseAuctionStudentRequired(aResponseAuctionStudent)
	})
}
