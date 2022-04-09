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

type ClassesClassBody struct {

	Name string `json:"name,omitempty"`

	Period int32 `json:"period,omitempty"`

	Id string `json:"_id,omitempty"`
}

// AssertClassesClassBodyRequired checks if the required fields are not zero-ed
func AssertClassesClassBodyRequired(obj ClassesClassBody) error {
	return nil
}

// AssertRecurseClassesClassBodyRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of ClassesClassBody (e.g. [][]ClassesClassBody), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseClassesClassBodyRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aClassesClassBody, ok := obj.(ClassesClassBody)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertClassesClassBodyRequired(aClassesClassBody)
	})
}
