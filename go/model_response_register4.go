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

type ResponseRegister4 struct {
	Message string `json:"message,omitempty"`
}

// AssertResponseRegister4Required checks if the required fields are not zero-ed
func AssertResponseRegister4Required(obj ResponseRegister4) error {
	return nil
}

// AssertRecurseResponseRegister4Required recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of ResponseRegister4 (e.g. [][]ResponseRegister4), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseResponseRegister4Required(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aResponseRegister4, ok := obj.(ResponseRegister4)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertResponseRegister4Required(aResponseRegister4)
	})
}
