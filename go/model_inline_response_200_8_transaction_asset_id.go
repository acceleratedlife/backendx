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

type InlineResponse2008TransactionAssetId struct {

	Name string `json:"name,omitempty"`

	Id string `json:"_id,omitempty"`
}

// AssertInlineResponse2008TransactionAssetIdRequired checks if the required fields are not zero-ed
func AssertInlineResponse2008TransactionAssetIdRequired(obj InlineResponse2008TransactionAssetId) error {
	return nil
}

// AssertRecurseInlineResponse2008TransactionAssetIdRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of InlineResponse2008TransactionAssetId (e.g. [][]InlineResponse2008TransactionAssetId), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseInlineResponse2008TransactionAssetIdRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aInlineResponse2008TransactionAssetId, ok := obj.(InlineResponse2008TransactionAssetId)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertInlineResponse2008TransactionAssetIdRequired(aInlineResponse2008TransactionAssetId)
	})
}
