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

type UsersResetPasswordBody struct {
	Email string `json:"email,omitempty"`
}

// AssertUsersResetPasswordBodyRequired checks if the required fields are not zero-ed
func AssertUsersResetPasswordBodyRequired(obj UsersResetPasswordBody) error {
	return nil
}

// AssertRecurseUsersResetPasswordBodyRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of UsersResetPasswordBody (e.g. [][]UsersResetPasswordBody), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseUsersResetPasswordBodyRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aUsersResetPasswordBody, ok := obj.(UsersResetPasswordBody)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertUsersResetPasswordBodyRequired(aUsersResetPasswordBody)
	})
}
