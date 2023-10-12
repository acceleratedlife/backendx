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

import (
	"time"
)

type ResponseCdTransaction struct {

	CreatedAt time.Time `json:"createdAt"`

	Mature time.Time `json:"mature"`

	Principal float32 `json:"principal"`

	Value float32 `json:"value"`
}

// AssertResponseCdTransactionRequired checks if the required fields are not zero-ed
func AssertResponseCdTransactionRequired(obj ResponseCdTransaction) error {
	elements := map[string]interface{}{
		"createdAt": obj.CreatedAt,
		"mature": obj.Mature,
		"principal": obj.Principal,
		"value": obj.Value,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertRecurseResponseCdTransactionRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of ResponseCdTransaction (e.g. [][]ResponseCdTransaction), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseResponseCdTransactionRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aResponseCdTransaction, ok := obj.(ResponseCdTransaction)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertResponseCdTransactionRequired(aResponseCdTransaction)
	})
}