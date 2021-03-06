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

type InlineResponse20017 struct {
	Success bool `json:"success,omitempty"`
}

// AssertInlineResponse20017Required checks if the required fields are not zero-ed
func AssertInlineResponse20017Required(obj InlineResponse20017) error {
	return nil
}

// AssertRecurseInlineResponse20017Required recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of InlineResponse20017 (e.g. [][]InlineResponse20017), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseInlineResponse20017Required(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aInlineResponse20017, ok := obj.(InlineResponse20017)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertInlineResponse20017Required(aInlineResponse20017)
	})
}
