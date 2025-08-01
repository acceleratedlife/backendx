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

type RequestMessage struct {

	Message string `json:"message"`

	UserId string `json:"user_id,omitempty"`

	SchoolId string `json:"school_id,omitempty"`
}

// AssertRequestMessageRequired checks if the required fields are not zero-ed
func AssertRequestMessageRequired(obj RequestMessage) error {
	elements := map[string]interface{}{
		"message": obj.Message,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertRecurseRequestMessageRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of RequestMessage (e.g. [][]RequestMessage), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseRequestMessageRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aRequestMessage, ok := obj.(RequestMessage)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertRequestMessageRequired(aRequestMessage)
	})
}
