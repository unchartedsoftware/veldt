package tile

import (
	"fmt"

	"github.com/unchartedsoftware/veldt/util/json"
)

// TopHits represents a tile which returns the top data points based on a
// provided field and sort order.
type TopHits struct {
	SortField     string
	SortOrder     string
	HitsCount     int
	IncludeFields []string
}

// Parse parses the provided JSON object and populates the tiles attributes.
func (t *TopHits) Parse(params map[string]interface{}) error {
	sortField := json.GetStringDefault(params, "", "sortField")
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
