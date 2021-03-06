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

type InlineResponse20014 struct {
	Owner ClassesmemberOwner `json:"owner,omitempty"`

	Period int32 `json:"period,omitempty"`

	Id string `json:"_id,omitempty"`
}

// AssertInlineResponse20014Required checks if the required fields are not zero-ed
func AssertInlineResponse20014Required(obj InlineResponse20014) error {
	if err := AssertClassesmemberOwnerRequired(obj.Owner); err != nil {
		return err
	}
	return nil
}

// AssertRecurseInlineResponse20014Required recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of InlineResponse20014 (e.g. [][]InlineResponse20014), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseInlineResponse20014Required(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aInlineResponse20014, ok := obj.(InlineResponse20014)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertInlineResponse20014Required(aInlineResponse20014)
	})
}
