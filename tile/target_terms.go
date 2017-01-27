package tile

import (
	"fmt"

	"github.com/unchartedsoftware/veldt/util/json"
)

// TargetTerms represents a tile which returns counts for specific terms in a
// provided field.
type TargetTerms struct {
	TermsField string
	Terms      []string
}

// Parse parses the provided JSON object and populates the tiles attributes.
func (t *TargetTerms) Parse(params map[string]interface{}) error {
	termsField, ok := json.GetString(params, "termsField")
	if !ok {
		return fmt.Errorf("`termsField` parameter missing from tile")
	}
	terms, ok := json.GetStringArray(params, "terms")
	if !ok {
		return fmt.Errorf("`terms` parameter missing from tile")
	}
	t.TermsField = termsField
	t.Terms = terms
	return nil
}
