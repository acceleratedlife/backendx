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

type ResponseLogin2 struct {

	LoginSuccess bool `json:"loginSuccess,omitempty"`

	UserId string `json:"userId,omitempty"`
}

// AssertResponseLogin2Required checks if the required fields are not zero-ed
func AssertResponseLogin2Required(obj ResponseLogin2) error {
	return nil
}

// AssertRecurseResponseLogin2Required recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of ResponseLogin2 (e.g. [][]ResponseLogin2), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseResponseLogin2Required(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aResponseLogin2, ok := obj.(ResponseLogin2)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertResponseLogin2Required(aResponseLogin2)
	})
}
