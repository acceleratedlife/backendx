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

type InlineResponse2004Transaction struct {

	Account string `json:"account,omitempty"`

	Amount float32 `json:"amount,omitempty"`

	Balance float32 `json:"balance,omitempty"`

	ConversionRatio float32 `json:"conversionRatio,omitempty"`

	Description string `json:"description,omitempty"`

	Owner string `json:"owner,omitempty"`

	Type string `json:"type,omitempty"`

	UBucks float32 `json:"uBucks,omitempty"`

	AssetID InlineResponse2004TransactionAssetId `json:"assetID,omitempty"`
}

// AssertInlineResponse2004TransactionRequired checks if the required fields are not zero-ed
func AssertInlineResponse2004TransactionRequired(obj InlineResponse2004Transaction) error {
	if err := AssertInlineResponse2004TransactionAssetIdRequired(obj.AssetID); err != nil {
		return err
	}
	return nil
}

// AssertRecurseInlineResponse2004TransactionRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of InlineResponse2004Transaction (e.g. [][]InlineResponse2004Transaction), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseInlineResponse2004TransactionRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aInlineResponse2004Transaction, ok := obj.(InlineResponse2004Transaction)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertInlineResponse2004TransactionRequired(aInlineResponse2004Transaction)
	})
}
