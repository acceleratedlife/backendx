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

type InlineResponse2006Accounts struct {

	Basis float32 `json:"basis,omitempty"`

	CurrentPrice float32 `json:"currentPrice,omitempty"`

	Name string `json:"name,omitempty"`

	Quantity float32 `json:"quantity,omitempty"`
}

// AssertInlineResponse2006AccountsRequired checks if the required fields are not zero-ed
func AssertInlineResponse2006AccountsRequired(obj InlineResponse2006Accounts) error {
	return nil
}

// AssertRecurseInlineResponse2006AccountsRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of InlineResponse2006Accounts (e.g. [][]InlineResponse2006Accounts), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseInlineResponse2006AccountsRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aInlineResponse2006Accounts, ok := obj.(InlineResponse2006Accounts)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertInlineResponse2006AccountsRequired(aInlineResponse2006Accounts)
	})
}
