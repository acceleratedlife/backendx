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

import "github.com/shopspring/decimal"

type Crypto struct {
	Basis decimal.Decimal `json:"basis,omitempty"`

	CurrentPrice decimal.Decimal `json:"currentPrice,omitempty"`

	Name string `json:"name,omitempty"`

	Quantity decimal.Decimal `json:"quantity,omitempty"`
}

// AssertCryptoRequired checks if the required fields are not zero-ed
func AssertCryptoRequired(obj Crypto) error {
	return nil
}

// AssertRecurseCryptoRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of Crypto (e.g. [][]Crypto), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseCryptoRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aCrypto, ok := obj.(Crypto)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertCryptoRequired(aCrypto)
	})
}
