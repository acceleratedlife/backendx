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

type ResponseAuth4 struct {

	IsAuth bool `json:"isAuth,omitempty"`

	Error bool `json:"error,omitempty"`
}

// AssertResponseAuth4Required checks if the required fields are not zero-ed
func AssertResponseAuth4Required(obj ResponseAuth4) error {
	return nil
}

// AssertRecurseResponseAuth4Required recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of ResponseAuth4 (e.g. [][]ResponseAuth4), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseResponseAuth4Required(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aResponseAuth4, ok := obj.(ResponseAuth4)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertResponseAuth4Required(aResponseAuth4)
	})
}
