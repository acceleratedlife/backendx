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

type InlineResponse20019Transaction struct {
	Account string `json:"account,omitempty"`

	Amount float32 `json:"amount,omitempty"`

	Balance float32 `json:"balance,omitempty"`

	ConversionRatio float32 `json:"conversionRatio,omitempty"`

	Description string `json:"description,omitempty"`

	Owner string `json:"owner,omitempty"`

	Type string `json:"type,omitempty"`

	UBucks float32 `json:"uBucks,omitempty"`

	AssetID InlineResponse20019TransactionAssetId `json:"assetID,omitempty"`
}

// AssertInlineResponse20019TransactionRequired checks if the required fields are not zero-ed
func AssertInlineResponse20019TransactionRequired(obj InlineResponse20019Transaction) error {
	if err := AssertInlineResponse20019TransactionAssetIdRequired(obj.AssetID); err != nil {
		return err
	}
	return nil
}

// AssertRecurseInlineResponse20019TransactionRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of InlineResponse20019Transaction (e.g. [][]InlineResponse20019Transaction), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseInlineResponse20019TransactionRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aInlineResponse20019Transaction, ok := obj.(InlineResponse20019Transaction)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertInlineResponse20019TransactionRequired(aInlineResponse20019Transaction)
	})
}
