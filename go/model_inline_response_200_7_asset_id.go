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

type InlineResponse2007AssetId struct {

	Name string `json:"name,omitempty"`
}

// AssertInlineResponse2007AssetIdRequired checks if the required fields are not zero-ed
func AssertInlineResponse2007AssetIdRequired(obj InlineResponse2007AssetId) error {
	return nil
}

// AssertRecurseInlineResponse2007AssetIdRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of InlineResponse2007AssetId (e.g. [][]InlineResponse2007AssetId), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseInlineResponse2007AssetIdRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aInlineResponse2007AssetId, ok := obj.(InlineResponse2007AssetId)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertInlineResponse2007AssetIdRequired(aInlineResponse2007AssetId)
	})
}
