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

type ClassKickBody struct {
	KickId string `json:"kick_id,omitempty"`

	Id string `json:"_id,omitempty"`
}

// AssertClassKickBodyRequired checks if the required fields are not zero-ed
func AssertClassKickBodyRequired(obj ClassKickBody) error {
	return nil
}

// AssertRecurseClassKickBodyRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of ClassKickBody (e.g. [][]ClassKickBody), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseClassKickBodyRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aClassKickBody, ok := obj.(ClassKickBody)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertClassKickBodyRequired(aClassKickBody)
	})
}
