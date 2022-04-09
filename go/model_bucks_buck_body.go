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

type BucksBuckBody struct {

	Name string `json:"name,omitempty"`

	Owner string `json:"owner,omitempty"`

	School string `json:"school,omitempty"`

	TotalCurrency float32 `json:"totalCurrency,omitempty"`

	FreeCurrency float32 `json:"freeCurrency,omitempty"`
}

// AssertBucksBuckBodyRequired checks if the required fields are not zero-ed
func AssertBucksBuckBodyRequired(obj BucksBuckBody) error {
	return nil
}

// AssertRecurseBucksBuckBodyRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of BucksBuckBody (e.g. [][]BucksBuckBody), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseBucksBuckBodyRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aBucksBuckBody, ok := obj.(BucksBuckBody)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertBucksBuckBodyRequired(aBucksBuckBody)
	})
}
