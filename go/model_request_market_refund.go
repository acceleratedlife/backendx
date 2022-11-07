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

type RequestMarketRefund struct {

	Id string `json:"_id,omitempty"`

	UserId string `json:"user_id,omitempty"`

	TeacherId string `json:"teacher_id,omitempty"`
}

// AssertRequestMarketRefundRequired checks if the required fields are not zero-ed
func AssertRequestMarketRefundRequired(obj RequestMarketRefund) error {
	return nil
}

// AssertRecurseRequestMarketRefundRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of RequestMarketRefund (e.g. [][]RequestMarketRefund), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseRequestMarketRefundRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aRequestMarketRefund, ok := obj.(RequestMarketRefund)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertRequestMarketRefundRequired(aRequestMarketRefund)
	})
}
