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

type InlineResponse2009 struct {

	Deleted bool `json:"deleted,omitempty"`
}

// AssertInlineResponse2009Required checks if the required fields are not zero-ed
func AssertInlineResponse2009Required(obj InlineResponse2009) error {
	return nil
}

// AssertRecurseInlineResponse2009Required recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of InlineResponse2009 (e.g. [][]InlineResponse2009), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseInlineResponse2009Required(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aInlineResponse2009, ok := obj.(InlineResponse2009)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertInlineResponse2009Required(aInlineResponse2009)
	})
}
