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

type ResponseResetPassword struct {

	Password string `json:"password"`
}

// AssertResponseResetPasswordRequired checks if the required fields are not zero-ed
func AssertResponseResetPasswordRequired(obj ResponseResetPassword) error {
	elements := map[string]interface{}{
		"password": obj.Password,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertRecurseResponseResetPasswordRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of ResponseResetPassword (e.g. [][]ResponseResetPassword), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseResponseResetPasswordRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aResponseResetPassword, ok := obj.(ResponseResetPassword)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertResponseResetPasswordRequired(aResponseResetPassword)
	})
}
