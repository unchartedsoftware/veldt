package tile

import (
	"fmt"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/util/json"
)

type Bivariate struct {
	XField     string
	YField     string
	Left       float64
	Right      float64
	Bottom     float64
	Top        float64
	Resolution int
	Bounds     *binning.Bounds
	BinSizeX   float64
	BinSizeY   float64
}

// Parse parses the provided JSON object and populates the tiles attributes.
func (b *Bivariate) Parse(params map[string]interface{}) error {
	// get x and y fields
	xField, ok := json.GetString(params, "xField")
	if !ok {
		return fmt.Errorf("`xField` parameter missing from tile")
	}
	yField, ok := json.GetString(params, "yField")
	if !ok {
		return fmt.Errorf("`yField` parameter missing from tile")
	}
	// get left, right, bottom, top extrema
	left, ok := json.GetNumber(params, "left")
	if !ok {
		return fmt.Errorf("`left` parameter missing from tile")
	}
	right, ok := json.GetNumber(params, "right")
	if !ok {
		return fmt.Errorf("`right` parameter missing from tile")
	}
	bottom, ok := json.GetNumber(params, "bottom")
	if !ok {
		return fmt.Errorf("`bottom` parameter missing from tile")
	}
	top, ok := json.GetNumber(params, "top")
	if !ok {
		return fmt.Errorf("`top` parameter missing from tile")
	}
	// get resolution
	resolution := json.GetNumberDefault(params, 256, "resolution")
	// set attributes
	b.XField = xField
	b.YField = yField
	b.Left = left
	b.Right = right
	b.Bottom = bottom
	b.Top = top
	b.Resolution = int(resolution)
	return nil
}

func (b *Bivariate) ClampBin(bin int64) int {
	if bin > int64(b.Resolution)-1 {
		return b.Resolution - 1
	}
	if bin < 0 {
		return 0
	}
	return int(bin)
}

func (b *Bivariate) GetXBin(x int64) int {
	bounds := b.Bounds
	fx := float64(x)
	var bin int64
	if bounds.BottomLeft.X > bounds.TopRight.X {
		bin = int64(float64(b.Resolution-1) - ((fx - bounds.TopRight.X) / b.BinSizeX))
	} else {
		bin = int64((fx - bounds.BottomLeft.X) / b.BinSizeX)
	}
	return b.ClampBin(bin)
}

// GetX given an x value, returns the corresponding coord within the range of
// [0 : 256) for the tile.
func (b *Bivariate) GetX(x float64) float64 {
	bounds := b.Bounds
	if bounds.BottomLeft.X > bounds.TopRight.X {
		rang := bounds.BottomLeft.X - bounds.TopRight.X
		return binning.MaxTileResolution - (((x - bounds.TopRight.X) / rang) * binning.MaxTileResolution)
	}
	rang := bounds.TopRight.X - bounds.BottomLeft.X
	return ((x - bounds.BottomLeft.X) / rang) * binning.MaxTileResolution
}

// GetYBin given an y value, returns the corresponding bin.
func (b *Bivariate) GetYBin(y int64) int {
	bounds := b.Bounds
	fy := float64(y)
	var bin int64
	if bounds.BottomLeft.Y > bounds.TopRight.Y {
		bin = int64(float64(b.Resolution-1) - ((fy - bounds.TopRight.Y) / b.BinSizeY))
	} else {
		bin = int64((fy - bounds.BottomLeft.Y) / b.BinSizeY)
	}
	return b.ClampBin(bin)
}

// GetY given an y value, returns the corresponding coord within the range of
// [0 : 256) for the tile.
func (b *Bivariate) GetY(y float64) float64 {
	bounds := b.Bounds
	if bounds.BottomLeft.Y > bounds.TopRight.Y {
		rang := bounds.BottomLeft.Y - bounds.TopRight.Y
		return binning.MaxTileResolution - (((y - bounds.TopRight.Y) / rang) * binning.MaxTileResolution)
	}
	rang := bounds.TopRight.Y - bounds.BottomLeft.Y
	return ((y - bounds.BottomLeft.Y) / rang) * binning.MaxTileResolution
}
