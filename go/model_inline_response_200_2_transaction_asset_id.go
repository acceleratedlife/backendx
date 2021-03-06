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

type InlineResponse2002TransactionAssetId struct {

	Name string `json:"name,omitempty"`

	Id string `json:"_id,omitempty"`
}

// AssertInlineResponse2002TransactionAssetIdRequired checks if the required fields are not zero-ed
func AssertInlineResponse2002TransactionAssetIdRequired(obj InlineResponse2002TransactionAssetId) error {
	return nil
}

// AssertRecurseInlineResponse2002TransactionAssetIdRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of InlineResponse2002TransactionAssetId (e.g. [][]InlineResponse2002TransactionAssetId), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseInlineResponse2002TransactionAssetIdRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aInlineResponse2002TransactionAssetId, ok := obj.(InlineResponse2002TransactionAssetId)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertInlineResponse2002TransactionAssetIdRequired(aInlineResponse2002TransactionAssetId)
	})
}
