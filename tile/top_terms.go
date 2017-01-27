package tile

import (
	"fmt"

	"github.com/unchartedsoftware/veldt/util/json"
)

// TopTerms represents a tile which returns counts for the top most occurring
// terms in a provided field.
type TopTerms struct {
	TermsField string
	TermsCount int
}

// Parse parses the provided JSON object and populates the tiles attributes.
func (t *TopTerms) Parse(params map[string]interface{}) error {
	termsField, ok := json.GetString(params, "termsField")
	if !ok {
		return fmt.Errorf("`termsField` parameter missing from tile")
	}
	termsCount, ok := json.GetNumber(params, "termsCount")
	if !ok {
		return fmt.Errorf("`termsCount` parameter missing from tile")
	}
	t.TermsField = termsField
	t.TermsCount = int(termsCount)
	return nil
}
