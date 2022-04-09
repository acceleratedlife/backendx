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

type Account struct {

	Id string `json:"_id"`

	OwnerId string `json:"owner_id,omitempty"`

	Kind string `json:"kind"`

	TypeId string `json:"type_id,omitempty"`

	Value float32 `json:"value"`

	Basis float32 `json:"basis"`

	History [][]int32 `json:"history"`
}

// AssertAccountRequired checks if the required fields are not zero-ed
func AssertAccountRequired(obj Account) error {
	elements := map[string]interface{}{
		"_id": obj.Id,
		"kind": obj.Kind,
		"value": obj.Value,
		"basis": obj.Basis,
		"history": obj.History,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertRecurseAccountRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of Account (e.g. [][]Account), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseAccountRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aAccount, ok := obj.(Account)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertAccountRequired(aAccount)
	})
}
