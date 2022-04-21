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

type InlineResponse20010TransactionAssetId struct {

	Name string `json:"name,omitempty"`

	Id string `json:"_id,omitempty"`
}

// AssertInlineResponse20010TransactionAssetIdRequired checks if the required fields are not zero-ed
func AssertInlineResponse20010TransactionAssetIdRequired(obj InlineResponse20010TransactionAssetId) error {
	return nil
}

// AssertRecurseInlineResponse20010TransactionAssetIdRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of InlineResponse20010TransactionAssetId (e.g. [][]InlineResponse20010TransactionAssetId), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseInlineResponse20010TransactionAssetIdRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aInlineResponse20010TransactionAssetId, ok := obj.(InlineResponse20010TransactionAssetId)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertInlineResponse20010TransactionAssetIdRequired(aInlineResponse20010TransactionAssetId)
	})
}
