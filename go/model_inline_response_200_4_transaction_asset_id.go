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

type InlineResponse2004TransactionAssetId struct {

	Name string `json:"name,omitempty"`

	Id string `json:"_id,omitempty"`
}

// AssertInlineResponse2004TransactionAssetIdRequired checks if the required fields are not zero-ed
func AssertInlineResponse2004TransactionAssetIdRequired(obj InlineResponse2004TransactionAssetId) error {
	return nil
}

// AssertRecurseInlineResponse2004TransactionAssetIdRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of InlineResponse2004TransactionAssetId (e.g. [][]InlineResponse2004TransactionAssetId), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseInlineResponse2004TransactionAssetIdRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aInlineResponse2004TransactionAssetId, ok := obj.(InlineResponse2004TransactionAssetId)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertInlineResponse2004TransactionAssetIdRequired(aInlineResponse2004TransactionAssetId)
	})
}
