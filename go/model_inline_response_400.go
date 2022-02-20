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

type InlineResponse400 struct {
	Error string `json:"error,omitempty"`
}

// AssertInlineResponse400Required checks if the required fields are not zero-ed
func AssertInlineResponse400Required(obj InlineResponse400) error {
	return nil
}

// AssertRecurseInlineResponse400Required recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of InlineResponse400 (e.g. [][]InlineResponse400), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseInlineResponse400Required(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aInlineResponse400, ok := obj.(InlineResponse400)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertInlineResponse400Required(aInlineResponse400)
	})
}
