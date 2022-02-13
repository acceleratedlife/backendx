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

type InlineResponse2008 struct {
	Name string `json:"name,omitempty"`

	Owner string `json:"owner,omitempty"`

	Period int32 `json:"period,omitempty"`

	AddCode string `json:"addCode,omitempty"`

	Members []string `json:"members,omitempty"`
}

// AssertInlineResponse2008Required checks if the required fields are not zero-ed
func AssertInlineResponse2008Required(obj InlineResponse2008) error {
	return nil
}

// AssertRecurseInlineResponse2008Required recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of InlineResponse2008 (e.g. [][]InlineResponse2008), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseInlineResponse2008Required(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aInlineResponse2008, ok := obj.(InlineResponse2008)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertInlineResponse2008Required(aInlineResponse2008)
	})
}
