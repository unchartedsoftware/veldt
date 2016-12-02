package tile

import (
	"fmt"

	"github.com/unchartedsoftware/prism/util/json"
)

type TopHits struct {
	SortField     string
	SortOrder     string
	HitsCount     int
	IncludeFields []string
}

// Parse parses the provided JSON object and populates the tiles attributes.
func (t *TopHits) Parse(params map[string]interface{}) error {
	sortField, ok := json.GetString(params, "sortField")
	if !ok {
		return fmt.Errorf("`sortField` parameter missing from tile")
	}
	sortOrder := json.GetStringDefault(params, "desc", "sortOrder")
	hitsCount, ok := json.GetNumber(params, "hitsCount")
	if !ok {
		return fmt.Errorf("`hitsCount` parameter missing from tile")
	}
	includeFields, ok := json.GetStringArray(params, "includeFields")
	if !ok {
		includeFields = nil
	}
	t.SortField = sortField
	t.SortOrder = sortOrder
	t.HitsCount = int(hitsCount)
	t.IncludeFields = includeFields
	return nil
}
