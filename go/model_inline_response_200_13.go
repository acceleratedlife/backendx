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

type InlineResponse20013 struct {
	Owner ClassesmemberOwner `json:"owner,omitempty"`

	Period int32 `json:"period,omitempty"`

	Id string `json:"_id,omitempty"`
}

// AssertInlineResponse20013Required checks if the required fields are not zero-ed
func AssertInlineResponse20013Required(obj InlineResponse20013) error {
	if err := AssertClassesmemberOwnerRequired(obj.Owner); err != nil {
		return err
	}
	return nil
}

// AssertRecurseInlineResponse20013Required recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of InlineResponse20013 (e.g. [][]InlineResponse20013), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseInlineResponse20013Required(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aInlineResponse20013, ok := obj.(InlineResponse20013)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertInlineResponse20013Required(aInlineResponse20013)
	})
}
