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

type ResponseLottoLatest struct {

	Jackpot int32 `json:"jackpot,omitempty"`

	Odds int32 `json:"odds,omitempty"`

	Winner string `json:"winner,omitempty"`
}

// AssertResponseLottoLatestRequired checks if the required fields are not zero-ed
func AssertResponseLottoLatestRequired(obj ResponseLottoLatest) error {
	return nil
}

// AssertRecurseResponseLottoLatestRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of ResponseLottoLatest (e.g. [][]ResponseLottoLatest), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseResponseLottoLatestRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aResponseLottoLatest, ok := obj.(ResponseLottoLatest)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertResponseLottoLatestRequired(aResponseLottoLatest)
	})
}
