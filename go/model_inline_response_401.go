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

type InlineResponse401 struct {

	Error string `json:"error,omitempty"`

	Message string `json:"message,omitempty"`
}

// AssertInlineResponse401Required checks if the required fields are not zero-ed
func AssertInlineResponse401Required(obj InlineResponse401) error {
	return nil
}

// AssertRecurseInlineResponse401Required recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of InlineResponse401 (e.g. [][]InlineResponse401), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseInlineResponse401Required(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aInlineResponse401, ok := obj.(InlineResponse401)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertInlineResponse401Required(aInlineResponse401)
	})
}
