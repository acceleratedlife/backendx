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

type ResponseMemberClass struct {

	Owner ResponseMemberClassOwner `json:"owner,omitempty"`

	Period int32 `json:"period,omitempty"`

	Id string `json:"_id,omitempty"`
}

// AssertResponseMemberClassRequired checks if the required fields are not zero-ed
func AssertResponseMemberClassRequired(obj ResponseMemberClass) error {
	if err := AssertResponseMemberClassOwnerRequired(obj.Owner); err != nil {
		return err
	}
	return nil
}

// AssertRecurseResponseMemberClassRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of ResponseMemberClass (e.g. [][]ResponseMemberClass), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseResponseMemberClassRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aResponseMemberClass, ok := obj.(ResponseMemberClass)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertResponseMemberClassRequired(aResponseMemberClass)
	})
}
