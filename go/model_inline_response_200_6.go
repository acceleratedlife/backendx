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

type InlineResponse2006 struct {
	Id string `json:"_id,omitempty"`

	Bid float32 `json:"bid,omitempty"`

	Description string `json:"description,omitempty"`

	EndDate time.Time `json:"endDate,omitempty"`

	StartDate time.Time `json:"startDate,omitempty"`

	ItemNumber string `json:"itemNumber,omitempty"`

	Owner AuctionsstudentOwner `json:"owner,omitempty"`

	Winner AuctionsstudentWinner `json:"winner,omitempty"`
}

// AssertInlineResponse2006Required checks if the required fields are not zero-ed
func AssertInlineResponse2006Required(obj InlineResponse2006) error {
	if err := AssertAuctionsstudentOwnerRequired(obj.Owner); err != nil {
		return err
	}
	if err := AssertAuctionsstudentWinnerRequired(obj.Winner); err != nil {
		return err
	}
	return nil
}

// AssertRecurseInlineResponse2006Required recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of InlineResponse2006 (e.g. [][]InlineResponse2006), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseInlineResponse2006Required(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aInlineResponse2006, ok := obj.(InlineResponse2006)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertInlineResponse2006Required(aInlineResponse2006)
	})
}
