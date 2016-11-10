package generation

import (
	"fmt"
	"math"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/util/json"
)

// Bivariate represents a bivariate tile generator.
type Bivariate struct {
	XField     string
	YField     string
	Bounds     *binning.Bounds
	MinX       float64
	MaxX       float64
	MinY       float64
	MaxY       float64
	XRange     float64
	YRange     float64
	Resolution int
	BinSizeX   float64
	BinSizeY   float64
}

// Parse parses the provided JSON object and populates the tiles attributes.
func (b *Bivariate) Parse(coord *binning.TileCoord, params map[string]interface{}) error {
	// get x and y fields
	xField, ok := json.GetString(params, "xField")
	if !ok {
		return fmt.Errorf("`xField` parameter missing from tiling params")
	}
	yField, ok := json.GetString(params, "yField")
	if !ok {
		return fmt.Errorf("`yField` parameter missing from tiling params")
	}
	// get left, right, bottom, top extrema
	left, ok := json.GetNumber(params, "left")
	if !ok {
		return fmt.Errorf("`left` parameter missing from tiling params")
	}
	right, ok := json.GetNumber(params, "right")
	if !ok {
		return fmt.Errorf("`right` parameter missing from tiling params")
	}
	bottom, ok := json.GetNumber(params, "bottom")
	if !ok {
		return fmt.Errorf("`bottom` parameter missing from tiling params")
	}
	top, ok := json.GetNumber(params, "top")
	if !ok {
		return fmt.Errorf("`top` parameter missing from tiling params")
	}
	// get resolution
	resolution := json.GetNumberDefault(params, 256, "resolution")
	// get the tiles bounds
	extents := &binning.Bounds{
		TopLeft: &binning.Coord{
			X: left,
			Y: top,
		},
		BottomRight: &binning.Coord{
			X: right,
			Y: bottom,
		},
	}
	bounds := binning.GetTileBounds(coord, extents)
	// get bin size
	xRange := math.Abs(bounds.BottomRight.X - bounds.TopLeft.X)
	yRange := math.Abs(bounds.BottomRight.Y - bounds.TopLeft.Y)
	binSizeX := xRange / resolution
	binSizeY := yRange / resolution
	// create bivariate
	b.XField = xField
	b.YField = yField
	b.Bounds = bounds
	b.MinX = math.Min(bounds.TopLeft.X, bounds.BottomRight.X)
	b.MaxX = math.Max(bounds.TopLeft.X, bounds.BottomRight.X)
	b.MinY = math.Min(bounds.TopLeft.Y, bounds.BottomRight.Y)
	b.MaxY = math.Max(bounds.TopLeft.Y, bounds.BottomRight.Y)
	b.XRange = xRange
	b.YRange = yRange
	// add binning params
	b.Resolution = int(resolution)
	b.BinSizeX = binSizeX
	b.BinSizeY = binSizeY
	return nil
}
