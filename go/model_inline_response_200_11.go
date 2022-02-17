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

type InlineResponse20011 struct {

	Owner ClassesmemberOwner `json:"owner,omitempty"`

	Period int32 `json:"period,omitempty"`

	Id string `json:"_id,omitempty"`
}

// AssertInlineResponse20011Required checks if the required fields are not zero-ed
func AssertInlineResponse20011Required(obj InlineResponse20011) error {
	if err := AssertClassesmemberOwnerRequired(obj.Owner); err != nil {
		return err
	}
	return nil
}

// AssertRecurseInlineResponse20011Required recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of InlineResponse20011 (e.g. [][]InlineResponse20011), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseInlineResponse20011Required(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aInlineResponse20011, ok := obj.(InlineResponse20011)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertInlineResponse20011Required(aInlineResponse20011)
	})
}
