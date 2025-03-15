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

type ResponseMemberClassOwner struct {

	FirstName string `json:"firstName,omitempty"`

	LastName string `json:"lastName,omitempty"`

	Id string `json:"_id,omitempty"`
}

// AssertResponseMemberClassOwnerRequired checks if the required fields are not zero-ed
func AssertResponseMemberClassOwnerRequired(obj ResponseMemberClassOwner) error {
	return nil
}

// AssertRecurseResponseMemberClassOwnerRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of ResponseMemberClassOwner (e.g. [][]ResponseMemberClassOwner), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseResponseMemberClassOwnerRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aResponseMemberClassOwner, ok := obj.(ResponseMemberClassOwner)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertResponseMemberClassOwnerRequired(aResponseMemberClassOwner)
	})
}
