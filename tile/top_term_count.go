package generation

import (
	"fmt"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/util/json"
)

// TopTermCount represents a tiling generator that produces heatmaps.
type TopTermCount struct {
	Bivariate
	TermField string
	SortField string
	Sort      string
}

// Parse parses the provided JSON object and populates the tiles attributes.
func (t *TopTermCount) Parse(params map[string]interface{}) error {
	err := t.Bivariate.Parse(params)
	if err != nil {
		return err
	}
	termField, ok := json.GetString(params, "termField")
	if !ok {
		return fmt.Errorf("`termField` parameter missing from tiling params")
	}
	sortField, ok := json.GetString(params, "sortField")
	if !ok {
		return fmt.Errorf("`sortField` parameter missing from tiling params")
	}
	sort := json.GetStringDefault(params, "desc", "sort")
	t.TermField = termField
	t.SortField = sortField
	t.Sort = sort
	return nil
}
