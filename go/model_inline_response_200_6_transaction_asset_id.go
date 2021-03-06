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

type InlineResponse2006TransactionAssetId struct {

	Name string `json:"name,omitempty"`

	Id string `json:"_id,omitempty"`
}

// AssertInlineResponse2006TransactionAssetIdRequired checks if the required fields are not zero-ed
func AssertInlineResponse2006TransactionAssetIdRequired(obj InlineResponse2006TransactionAssetId) error {
	return nil
}

// AssertRecurseInlineResponse2006TransactionAssetIdRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of InlineResponse2006TransactionAssetId (e.g. [][]InlineResponse2006TransactionAssetId), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseInlineResponse2006TransactionAssetIdRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aInlineResponse2006TransactionAssetId, ok := obj.(InlineResponse2006TransactionAssetId)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertInlineResponse2006TransactionAssetIdRequired(aInlineResponse2006TransactionAssetId)
	})
}
