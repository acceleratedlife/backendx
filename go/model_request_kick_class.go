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

type RequestKickClass struct {

	KickId string `json:"kick_id,omitempty"`

	Id string `json:"_id,omitempty"`
}

// AssertRequestKickClassRequired checks if the required fields are not zero-ed
func AssertRequestKickClassRequired(obj RequestKickClass) error {
	return nil
}

// AssertRecurseRequestKickClassRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of RequestKickClass (e.g. [][]RequestKickClass), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseRequestKickClassRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aRequestKickClass, ok := obj.(RequestKickClass)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertRequestKickClassRequired(aRequestKickClass)
	})
}
