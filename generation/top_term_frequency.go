package generation

import (
	"fmt"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/util/json"
)

// TopTermFrequency represents a tiling generator that produces heatmaps.
type TopTermFrequency struct {
	Tiling *param.Bivariate
	TermField string
	SortField string
	Sort      string
	From      int64
	To        int64
	Interval  int64
}

// SetTopTermFrequencyParams sets the params for the specific generator.
func SetTopTermFrequencyParams(arg interface{}, coord *binning.TileCoord, params map[string]interface{}) error {
	generator, ok := arg.(*TopTermFrequency)
	if !ok {
		return fmt.Errorf("`%v` is not of type `*TopTermFrequency`", arg)
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
	generator.TermField = termField
	generator.SortField = sortField
	generator.Sort = sort
	generator.From = int64(from)
	generator.To = int64(to)
	generator.Interval = int64(interval)
	return nil
}
