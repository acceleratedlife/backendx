/*
 * AL
 *
 * This is a simple API
 *
 * API version: 1.0.2
 * Contact: you@your-company.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type RequestPayTransaction struct {

	OwnerId string `json:"owner_id,omitempty"`

	Description string `json:"description,omitempty"`

	Amount float32 `json:"amount,omitempty"`

	Student string `json:"student,omitempty"`
}

// AssertRequestPayTransactionRequired checks if the required fields are not zero-ed
func AssertRequestPayTransactionRequired(obj RequestPayTransaction) error {
	return nil
}

// AssertRecurseRequestPayTransactionRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of RequestPayTransaction (e.g. [][]RequestPayTransaction), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseRequestPayTransactionRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aRequestPayTransaction, ok := obj.(RequestPayTransaction)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertRequestPayTransactionRequired(aRequestPayTransaction)
	})
}
