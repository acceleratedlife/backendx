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

type InlineResponse2003 struct {

	Converion float32 `json:"converion,omitempty"`

	History [][]int32 `json:"history,omitempty"`

	Bucks string `json:"bucks,omitempty"`

	Balance float32 `json:"balance,omitempty"`

	Id string `json:"_id,omitempty"`

	TypeId string `json:"type_id,omitempty"`
}

// AssertInlineResponse2003Required checks if the required fields are not zero-ed
func AssertInlineResponse2003Required(obj InlineResponse2003) error {
	return nil
}

// AssertRecurseInlineResponse2003Required recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of InlineResponse2003 (e.g. [][]InlineResponse2003), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseInlineResponse2003Required(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aInlineResponse2003, ok := obj.(InlineResponse2003)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertInlineResponse2003Required(aInlineResponse2003)
	})
}
