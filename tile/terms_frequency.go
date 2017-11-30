package tile

import (
"fmt"

"github.com/unchartedsoftware/veldt/util/json"
)

// TermsFrequency represents a tile which returns counts for each of the terms in a provided field.
// This should be used with fields with a finite set of values, as there are no limits provided.
// fieldType is an optional value representing the type of the field.  Currently only 'string' is supported, and all
// other fieldType values default to an array of strings.
type TermsFrequency struct {
	TermsField string
	FieldType  string
}

// Parse parses the provided JSON object and populates the tiles attributes.
func (t *TermsFrequency) Parse(params map[string]interface{}) error {
	termsField, ok := json.GetString(params, "termsField")
	if !ok {
		return fmt.Errorf("`termsField` parameter missing from tile")
	}
	fieldType, ok := json.GetString(params, "fieldType")
	if ok && fieldType == "string" {
		t.FieldType = fieldType
	}
	t.TermsField = termsField
	return nil
}
