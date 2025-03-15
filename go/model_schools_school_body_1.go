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

type SchoolsSchoolBody1 struct {

	Name string `json:"name,omitempty"`

	City string `json:"city,omitempty"`

	Zip int32 `json:"zip,omitempty"`
}

// AssertSchoolsSchoolBody1Required checks if the required fields are not zero-ed
func AssertSchoolsSchoolBody1Required(obj SchoolsSchoolBody1) error {
	return nil
}

// AssertRecurseSchoolsSchoolBody1Required recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of SchoolsSchoolBody1 (e.g. [][]SchoolsSchoolBody1), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseSchoolsSchoolBody1Required(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aSchoolsSchoolBody1, ok := obj.(SchoolsSchoolBody1)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertSchoolsSchoolBody1Required(aSchoolsSchoolBody1)
	})
}
