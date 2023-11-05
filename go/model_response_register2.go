/*
 * AL
 *
 * This is a simple API
 *
 * API version: 1.0.2
 * Contact: you@your-company.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type ResponseRegister2 struct {

	Success bool `json:"success,omitempty"`
}

// AssertResponseRegister2Required checks if the required fields are not zero-ed
func AssertResponseRegister2Required(obj ResponseRegister2) error {
	return nil
}

// AssertRecurseResponseRegister2Required recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of ResponseRegister2 (e.g. [][]ResponseRegister2), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseResponseRegister2Required(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aResponseRegister2, ok := obj.(ResponseRegister2)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertResponseRegister2Required(aResponseRegister2)
	})
}
