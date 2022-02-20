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

type InlineResponse20012 struct {

	Owner ClassesmemberOwner `json:"owner,omitempty"`

	Period int32 `json:"period,omitempty"`

	Id string `json:"_id,omitempty"`
}

// AssertInlineResponse20012Required checks if the required fields are not zero-ed
func AssertInlineResponse20012Required(obj InlineResponse20012) error {
	if err := AssertClassesmemberOwnerRequired(obj.Owner); err != nil {
		return err
	}
	return nil
}

// AssertRecurseInlineResponse20012Required recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of InlineResponse20012 (e.g. [][]InlineResponse20012), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseInlineResponse20012Required(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aInlineResponse20012, ok := obj.(InlineResponse20012)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertInlineResponse20012Required(aInlineResponse20012)
	})
}
