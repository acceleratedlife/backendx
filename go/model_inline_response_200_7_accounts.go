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

type InlineResponse2007Accounts struct {

	Basis float32 `json:"basis,omitempty"`

	CurrentPrice float32 `json:"currentPrice,omitempty"`

	Name string `json:"name,omitempty"`

	Quantity float32 `json:"quantity,omitempty"`
}

// AssertInlineResponse2007AccountsRequired checks if the required fields are not zero-ed
func AssertInlineResponse2007AccountsRequired(obj InlineResponse2007Accounts) error {
	return nil
}

// AssertRecurseInlineResponse2007AccountsRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of InlineResponse2007Accounts (e.g. [][]InlineResponse2007Accounts), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseInlineResponse2007AccountsRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aInlineResponse2007Accounts, ok := obj.(InlineResponse2007Accounts)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertInlineResponse2007AccountsRequired(aInlineResponse2007Accounts)
	})
}
