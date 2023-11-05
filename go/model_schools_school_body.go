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

type SchoolsSchoolBody struct {

	Name string `json:"name,omitempty"`

	City string `json:"city,omitempty"`

	Zip int32 `json:"zip,omitempty"`
}

// AssertSchoolsSchoolBodyRequired checks if the required fields are not zero-ed
func AssertSchoolsSchoolBodyRequired(obj SchoolsSchoolBody) error {
	return nil
}

// AssertRecurseSchoolsSchoolBodyRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of SchoolsSchoolBody (e.g. [][]SchoolsSchoolBody), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseSchoolsSchoolBodyRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aSchoolsSchoolBody, ok := obj.(SchoolsSchoolBody)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertSchoolsSchoolBodyRequired(aSchoolsSchoolBody)
	})
}
