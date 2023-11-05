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

type ResponseAuth2 struct {

	Email string `json:"email,omitempty"`

	FirstName string `json:"firstName,omitempty"`

	LastName string `json:"lastName,omitempty"`

	IsAdmin bool `json:"isAdmin"`

	IsAuth bool `json:"isAuth"`

	Role int32 `json:"role"`

	SchoolId string `json:"school_id,omitempty"`

	Id string `json:"_id,omitempty"`

	LottoPlay int32 `json:"lottoPlay,omitempty"`

	LottoWin int32 `json:"lottoWin,omitempty"`
}

// AssertResponseAuth2Required checks if the required fields are not zero-ed
func AssertResponseAuth2Required(obj ResponseAuth2) error {
	elements := map[string]interface{}{
		"isAdmin": obj.IsAdmin,
		"isAuth": obj.IsAuth,
		"role": obj.Role,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertRecurseResponseAuth2Required recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of ResponseAuth2 (e.g. [][]ResponseAuth2), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseResponseAuth2Required(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aResponseAuth2, ok := obj.(ResponseAuth2)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertResponseAuth2Required(aResponseAuth2)
	})
}
