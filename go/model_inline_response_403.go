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

type InlineResponse403 struct {

	Error string `json:"error,omitempty"`

	Message string `json:"message,omitempty"`
}

// AssertInlineResponse403Required checks if the required fields are not zero-ed
func AssertInlineResponse403Required(obj InlineResponse403) error {
	return nil
}

// AssertRecurseInlineResponse403Required recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of InlineResponse403 (e.g. [][]InlineResponse403), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseInlineResponse403Required(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aInlineResponse403, ok := obj.(InlineResponse403)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertInlineResponse403Required(aInlineResponse403)
	})
}
