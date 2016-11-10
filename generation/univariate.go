package generation

import (
	"fmt"
	"math"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/util/json"
)

// Univariate represents a univariate tile generator.
type Univariate struct {
	Field      string
	Min        float64
	Max        float64
	Range      float64
	Resolution int
	BinSize    float64
}

// Parse parses the provided JSON object and populates the tiles attributes.
func (u *Univariate) Parse(coord *binning.TileCoord, params map[string]interface{}) error {
	// get x and y fields
	field, ok := json.GetString(params, "field")
	if !ok {
		return fmt.Errorf("`field` parameter missing from tiling params")
	}
	// get left, right, bottom, top extrema
	min, ok := json.GetNumber(params, "min")
	if !ok {
		return fmt.Errorf("`min` parameter missing from tiling params")
	}
	max, ok := json.GetNumber(params, "max")
	if !ok {
		return fmt.Errorf("`max` parameter missing from tiling params")
	}
	extrema := binning.GetTileExtrema(coord.X, coord.Z, &binning.Extrema{
		Min: min,
		Max: max,
	})
	// get resolution
	resolution := json.GetNumberDefault(params, 256, "resolution")
	// get bin size
	rang := math.Abs(extrema.Min - extrema.Max)
	binSize := rang / resolution
	// create univariate
	u.Field = field
	u.Min = extrema.Min
	u.Max = extrema.Max
	u.Range = rang
	// add binning params
	u.Resolution = int(resolution)
	u.BinSize = binSize
	return u, nil
}
