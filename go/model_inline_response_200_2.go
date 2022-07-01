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

type InlineResponse2002 struct {

	Transaction InlineResponse2002Transaction `json:"transaction,omitempty"`

	Accounts []InlineResponse2002Accounts `json:"accounts,omitempty"`
}

// AssertInlineResponse2002Required checks if the required fields are not zero-ed
func AssertInlineResponse2002Required(obj InlineResponse2002) error {
	if err := AssertInlineResponse2002TransactionRequired(obj.Transaction); err != nil {
		return err
	}
	for _, el := range obj.Accounts {
		if err := AssertInlineResponse2002AccountsRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertRecurseInlineResponse2002Required recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of InlineResponse2002 (e.g. [][]InlineResponse2002), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseInlineResponse2002Required(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aInlineResponse2002, ok := obj.(InlineResponse2002)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertInlineResponse2002Required(aInlineResponse2002)
	})
}
