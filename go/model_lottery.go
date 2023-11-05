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

type Lottery struct {

	Odds int32 `json:"odds"`

	Jackpot int32 `json:"jackpot"`

	Number int32 `json:"number"`

	Winner string `json:"winner,omitempty"`
}

// AssertLotteryRequired checks if the required fields are not zero-ed
func AssertLotteryRequired(obj Lottery) error {
	elements := map[string]interface{}{
		"odds": obj.Odds,
		"jackpot": obj.Jackpot,
		"number": obj.Number,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertRecurseLotteryRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of Lottery (e.g. [][]Lottery), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseLotteryRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aLottery, ok := obj.(Lottery)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertLotteryRequired(aLottery)
	})
}
