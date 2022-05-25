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

type ResponseSearchStudentUbuck struct {

	Value float32 `json:"value,omitempty"`
}

// AssertResponseSearchStudentUbuckRequired checks if the required fields are not zero-ed
func AssertResponseSearchStudentUbuckRequired(obj ResponseSearchStudentUbuck) error {
	return nil
}

// AssertRecurseResponseSearchStudentUbuckRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of ResponseSearchStudentUbuck (e.g. [][]ResponseSearchStudentUbuck), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseResponseSearchStudentUbuckRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aResponseSearchStudentUbuck, ok := obj.(ResponseSearchStudentUbuck)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertResponseSearchStudentUbuckRequired(aResponseSearchStudentUbuck)
	})
}
