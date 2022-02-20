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

type AuctionsstudentOwner struct {

	Id string `json:"_id,omitempty"`

	LastName string `json:"lastName,omitempty"`
}

// AssertAuctionsstudentOwnerRequired checks if the required fields are not zero-ed
func AssertAuctionsstudentOwnerRequired(obj AuctionsstudentOwner) error {
	return nil
}

// AssertRecurseAuctionsstudentOwnerRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of AuctionsstudentOwner (e.g. [][]AuctionsstudentOwner), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseAuctionsstudentOwnerRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aAuctionsstudentOwner, ok := obj.(AuctionsstudentOwner)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertAuctionsstudentOwnerRequired(aAuctionsstudentOwner)
	})
}
