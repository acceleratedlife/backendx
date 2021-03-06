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

type ClassesAddClassBody struct {
	AddCode string `json:"addCode,omitempty"`

	Id string `json:"_id,omitempty"`
}

// AssertClassesAddClassBodyRequired checks if the required fields are not zero-ed
func AssertClassesAddClassBodyRequired(obj ClassesAddClassBody) error {
	return nil
}

// AssertRecurseClassesAddClassBodyRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of ClassesAddClassBody (e.g. [][]ClassesAddClassBody), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseClassesAddClassBodyRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aClassesAddClassBody, ok := obj.(ClassesAddClassBody)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertClassesAddClassBodyRequired(aClassesAddClassBody)
	})
}
