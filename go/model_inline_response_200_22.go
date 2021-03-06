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

type InlineResponse20022 struct {
	Deleted bool `json:"deleted,omitempty"`
}

// AssertInlineResponse20022Required checks if the required fields are not zero-ed
func AssertInlineResponse20022Required(obj InlineResponse20022) error {
	return nil
}

// AssertRecurseInlineResponse20022Required recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of InlineResponse20022 (e.g. [][]InlineResponse20022), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseInlineResponse20022Required(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aInlineResponse20022, ok := obj.(InlineResponse20022)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertInlineResponse20022Required(aInlineResponse20022)
	})
}
