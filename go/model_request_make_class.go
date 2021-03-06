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

type RequestMakeClass struct {

	Name string `json:"name,omitempty"`

	OwnerId string `json:"owner_id,omitempty"`

	Period int32 `json:"period,omitempty"`
}

// AssertRequestMakeClassRequired checks if the required fields are not zero-ed
func AssertRequestMakeClassRequired(obj RequestMakeClass) error {
	return nil
}

// AssertRecurseRequestMakeClassRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of RequestMakeClass (e.g. [][]RequestMakeClass), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseRequestMakeClassRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aRequestMakeClass, ok := obj.(RequestMakeClass)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertRequestMakeClassRequired(aRequestMakeClass)
	})
}
