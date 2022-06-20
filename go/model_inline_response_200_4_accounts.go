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

type InlineResponse2004Accounts struct {

	Basis float32 `json:"basis,omitempty"`

	CurrentPrice float32 `json:"currentPrice,omitempty"`

	Name string `json:"name,omitempty"`

	Quantity float32 `json:"quantity,omitempty"`
}

// AssertInlineResponse2004AccountsRequired checks if the required fields are not zero-ed
func AssertInlineResponse2004AccountsRequired(obj InlineResponse2004Accounts) error {
	return nil
}

// AssertRecurseInlineResponse2004AccountsRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of InlineResponse2004Accounts (e.g. [][]InlineResponse2004Accounts), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseInlineResponse2004AccountsRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aInlineResponse2004Accounts, ok := obj.(InlineResponse2004Accounts)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertInlineResponse2004AccountsRequired(aInlineResponse2004Accounts)
	})
}