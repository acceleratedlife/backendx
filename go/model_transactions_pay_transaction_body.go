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

type TransactionsPayTransactionBody struct {

	Owner string `json:"owner,omitempty"`

	Description string `json:"description,omitempty"`

	Amount float32 `json:"amount,omitempty"`

	Student string `json:"student,omitempty"`

	Kind string `json:"kind,omitempty"`
}

// AssertTransactionsPayTransactionBodyRequired checks if the required fields are not zero-ed
func AssertTransactionsPayTransactionBodyRequired(obj TransactionsPayTransactionBody) error {
	return nil
}

// AssertRecurseTransactionsPayTransactionBodyRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of TransactionsPayTransactionBody (e.g. [][]TransactionsPayTransactionBody), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseTransactionsPayTransactionBodyRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aTransactionsPayTransactionBody, ok := obj.(TransactionsPayTransactionBody)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertTransactionsPayTransactionBodyRequired(aTransactionsPayTransactionBody)
	})
}
