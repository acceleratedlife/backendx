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

type ResponseCrypto struct {

	Searched string `json:"searched,omitempty"`

	Usd float32 `json:"usd,omitempty"`

	Owned float32 `json:"owned,omitempty"`

	UBuck float32 `json:"UBuck,omitempty"`

	Basis float32 `json:"basis,omitempty"`
}

// AssertResponseCryptoRequired checks if the required fields are not zero-ed
func AssertResponseCryptoRequired(obj ResponseCrypto) error {
	return nil
}

// AssertRecurseResponseCryptoRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of ResponseCrypto (e.g. [][]ResponseCrypto), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseResponseCryptoRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aResponseCrypto, ok := obj.(ResponseCrypto)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertResponseCryptoRequired(aResponseCrypto)
	})
}
