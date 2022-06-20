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

type ResponseCurrencyExchangeBuck struct {

	History []History `json:"history,omitempty"`

	Name string `json:"name,omitempty"`
}

// AssertResponseCurrencyExchangeBuckRequired checks if the required fields are not zero-ed
func AssertResponseCurrencyExchangeBuckRequired(obj ResponseCurrencyExchangeBuck) error {
	for _, el := range obj.History {
		if err := AssertHistoryRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertRecurseResponseCurrencyExchangeBuckRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of ResponseCurrencyExchangeBuck (e.g. [][]ResponseCurrencyExchangeBuck), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseResponseCurrencyExchangeBuckRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aResponseCurrencyExchangeBuck, ok := obj.(ResponseCurrencyExchangeBuck)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertResponseCurrencyExchangeBuckRequired(aResponseCurrencyExchangeBuck)
	})
}