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

import (
	"time"
)

type UserNoHistory struct {
	Id string `json:"_id"`

	CollegeEnd time.Time `json:"collegeEnd,omitempty"`

	TransitionEnd time.Time `json:"transitionEnd,omitempty"`

	FirstName string `json:"firstName"`

	LastName string `json:"lastName"`

	Email string `json:"email"`

	Confirmed bool `json:"confirmed"`

	SchoolId string `json:"school_id,omitempty"`

	College bool `json:"college"`

	Children int32 `json:"children"`

	Income float32 `json:"income"`

	Role int32 `json:"role"`

	Rank int32 `json:"rank"`

	NetWorth float32 `json:"netWorth"`
}

// AssertUserNoHistoryRequired checks if the required fields are not zero-ed
func AssertUserNoHistoryRequired(obj UserNoHistory) error {
	elements := map[string]interface{}{
		"_id":       obj.Id,
		"firstName": obj.FirstName,
		"lastName":  obj.LastName,
		"email":     obj.Email,
		"confirmed": obj.Confirmed,
		"college":   obj.College,
		"children":  obj.Children,
		"income":    obj.Income,
		"role":      obj.Role,
		"rank":      obj.Rank,
		"netWorth":  obj.NetWorth,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertRecurseUserNoHistoryRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of UserNoHistory (e.g. [][]UserNoHistory), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseUserNoHistoryRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aUserNoHistory, ok := obj.(UserNoHistory)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertUserNoHistoryRequired(aUserNoHistory)
	})
}
