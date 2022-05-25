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

type AuctionsauctionVisibility struct {

	Name string `json:"name,omitempty"`
}

// AssertAuctionsauctionVisibilityRequired checks if the required fields are not zero-ed
func AssertAuctionsauctionVisibilityRequired(obj AuctionsauctionVisibility) error {
	return nil
}

// AssertRecurseAuctionsauctionVisibilityRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of AuctionsauctionVisibility (e.g. [][]AuctionsauctionVisibility), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseAuctionsauctionVisibilityRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aAuctionsauctionVisibility, ok := obj.(AuctionsauctionVisibility)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertAuctionsauctionVisibilityRequired(aAuctionsauctionVisibility)
	})
}
