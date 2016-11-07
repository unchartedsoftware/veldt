package generation

import (
	"fmt"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/util/json"
)

// TopTermCount represents a tiling generator that produces heatmaps.
type TopTermCount struct {
	Tiling *param.Bivariate
	TermField string
	SortField string
	Sort      string
}

// SetTopTermCountParams sets the params for the specific generator.
func SetTopTermCountParams(arg interface{}, coord *binning.TileCoord, params map[string]interface{}) error {
	generator, ok := arg.(*TopTermCount)
	if !ok {
		return fmt.Errorf("`%v` is not of type `*TopTermCount`", arg)
	}
	err := SetBivariateParams(generator, coord, params)
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
	generator.TermField = termField
	generator.SortField = sortField
	generator.Sort = sort
	return nil
}
