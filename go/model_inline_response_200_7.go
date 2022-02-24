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

type InlineResponse2007 struct {

	Account string `json:"account,omitempty"`

	Owner string `json:"owner,omitempty"`

	Balance float32 `json:"balance,omitempty"`

	Description string `json:"description,omitempty"`

	ConversionRatio float32 `json:"conversionRatio,omitempty"`

	Amount float32 `json:"amount,omitempty"`

	UBucks float32 `json:"uBucks,omitempty"`

	Type string `json:"type,omitempty"`

	AssetID InlineResponse2007AssetId `json:"assetID,omitempty"`
}

// AssertInlineResponse2007Required checks if the required fields are not zero-ed
func AssertInlineResponse2007Required(obj InlineResponse2007) error {
	if err := AssertInlineResponse2007AssetIdRequired(obj.AssetID); err != nil {
		return err
	}
	return nil
}

// AssertRecurseInlineResponse2007Required recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of InlineResponse2007 (e.g. [][]InlineResponse2007), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseInlineResponse2007Required(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aInlineResponse2007, ok := obj.(InlineResponse2007)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertInlineResponse2007Required(aInlineResponse2007)
	})
}
