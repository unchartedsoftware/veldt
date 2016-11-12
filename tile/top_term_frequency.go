package tile

import (
	"fmt"

	"github.com/unchartedsoftware/prism/util/json"
)

// TopTermFrequency represents a tiling generator that produces heatmaps.
type TopTermFrequency struct {
	Bivariate
	TermField string
	SortField string
	Sort      string
	From      int64
	To        int64
	Interval  int64
}

// Parse parses the provided JSON object and populates the tiles attributes.
func (t *TopTermFrequency) Parse(params map[string]interface{}) error {
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
	from, ok := json.GetNumber(params, "from")
	if !ok {
		return fmt.Errorf("`sort` parameter missing from tiling params")
	}
	to, ok := json.GetNumber(params, "to")
	if !ok {
		return fmt.Errorf("`sort` parameter missing from tiling params")
	}
	interval, ok := json.GetNumber(params, "interval")
	if !ok {
		return fmt.Errorf("`sort` parameter missing from tiling params")
	}
	t.TermField = termField
	t.SortField = sortField
	t.Sort = sort
	t.From = int64(from)
	t.To = int64(to)
	t.Interval = int64(interval)
	return nil
}
