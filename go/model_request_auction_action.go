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

type RequestAuctionAction struct {

	AuctionId string `json:"auctionId,omitempty"`
}

// AssertRequestAuctionActionRequired checks if the required fields are not zero-ed
func AssertRequestAuctionActionRequired(obj RequestAuctionAction) error {
	return nil
}

// AssertRecurseRequestAuctionActionRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of RequestAuctionAction (e.g. [][]RequestAuctionAction), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseRequestAuctionActionRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aRequestAuctionAction, ok := obj.(RequestAuctionAction)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertRequestAuctionActionRequired(aRequestAuctionAction)
	})
}