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

type TransactionsConversionTransactionBody struct {
	AccountFrom string `json:"accountFrom,omitempty"`

	AccountTo string `json:"accountTo,omitempty"`

	Amount float32 `json:"amount,omitempty"`
}

// AssertTransactionsConversionTransactionBodyRequired checks if the required fields are not zero-ed
func AssertTransactionsConversionTransactionBodyRequired(obj TransactionsConversionTransactionBody) error {
	return nil
}

// AssertRecurseTransactionsConversionTransactionBodyRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of TransactionsConversionTransactionBody (e.g. [][]TransactionsConversionTransactionBody), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseTransactionsConversionTransactionBodyRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aTransactionsConversionTransactionBody, ok := obj.(TransactionsConversionTransactionBody)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertTransactionsConversionTransactionBodyRequired(aTransactionsConversionTransactionBody)
	})
}
