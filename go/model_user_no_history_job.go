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

type UserNoHistoryJob struct {
	Title string `json:"title,omitempty"`

	Pay int32 `json:"pay,omitempty"`

	Description string `json:"description,omitempty"`
}

// AssertUserNoHistoryJobRequired checks if the required fields are not zero-ed
func AssertUserNoHistoryJobRequired(obj UserNoHistoryJob) error {
	return nil
}

// AssertRecurseUserNoHistoryJobRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of UserNoHistoryJob (e.g. [][]UserNoHistoryJob), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseUserNoHistoryJobRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aUserNoHistoryJob, ok := obj.(UserNoHistoryJob)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertUserNoHistoryJobRequired(aUserNoHistoryJob)
	})
}
