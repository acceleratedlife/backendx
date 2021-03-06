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

type InlineResponse20018AssetId struct {
	Name string `json:"name,omitempty"`
}

// AssertInlineResponse20018AssetIdRequired checks if the required fields are not zero-ed
func AssertInlineResponse20018AssetIdRequired(obj InlineResponse20018AssetId) error {
	return nil
}

// AssertRecurseInlineResponse20018AssetIdRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of InlineResponse20018AssetId (e.g. [][]InlineResponse20018AssetId), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseInlineResponse20018AssetIdRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aInlineResponse20018AssetId, ok := obj.(InlineResponse20018AssetId)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertInlineResponse20018AssetIdRequired(aInlineResponse20018AssetId)
	})
}
