package tile

import (
	"fmt"
	"math"

	"github.com/unchartedsoftware/veldt/binning"
	"github.com/unchartedsoftware/veldt/util/json"
)

// Bivariate represents the parameters required for any bivariate tile.
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

// GetXBin given an x value, returns the corresponding bin.
func (b *Bivariate) GetXBin(x int64) int {
	bounds := b.Bounds
	fx := float64(x)
	var bin int64
	if bounds.BottomLeft.X > bounds.TopRight.X {
		bin = int64(float64(b.Resolution-1) - ((fx - bounds.TopRight.X) / b.BinSizeX))
	} else {
		bin = int64((fx - bounds.BottomLeft.X) / b.BinSizeX)
	}
	return b.clampBin(bin)
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

// GetYBin given a y value, returns the corresponding bin.
func (b *Bivariate) GetYBin(y int64) int {
	bounds := b.Bounds
	fy := float64(y)
	var bin int64
	if bounds.BottomLeft.Y > bounds.TopRight.Y {
		bin = int64(float64(b.Resolution-1) - ((fy - bounds.TopRight.Y) / b.BinSizeY))
	} else {
		bin = int64((fy - bounds.BottomLeft.Y) / b.BinSizeY)
	}
	return b.clampBin(bin)
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

// GetXY given a data hit, returns the corresponding coord within the range of
// [0 : 256) for the tile.
func (b *Bivariate) GetXY(hit map[string]interface{}) (float32, float32, bool) {
	// get x / y fields from data
	ix, ok := hit[b.XField]
	if !ok {
		return 0, 0, false
	}
	iy, ok := hit[b.YField]
	if !ok {
		return 0, 0, false
	}
	// get X / Y of the data
	x, y, ok := castPixel(ix, iy)
	if !ok {
		return 0, 0, false
	}
	// convert to tile pixel coords in the range [0 - 256)
	tx := b.GetX(x)
	ty := b.GetY(y)
	// return position in tile coords with 2 decimal places
	return toFixed(float32(tx), 2), toFixed(float32(ty), 2), true
}

func (b *Bivariate) clampBin(bin int64) int {
	if bin > int64(b.Resolution)-1 {
		return b.Resolution - 1
	}
	if bin < 0 {
		return 0
	}
	return int(bin)
}

func toFixed(num float32, precision int) float32 {
	output := math.Pow(10, float64(precision))
	return float32(math.Floor(float64(num)*output+0.5)) / float32(output)
}

func castPixel(x interface{}, y interface{}) (float64, float64, bool) {
	xfval, xok := x.(float64)
	yfval, yok := y.(float64)
	if xok && yok {
		return xfval, yfval, true
	}
	xival, xok := x.(int64)
	yival, yok := y.(int64)
	if xok && yok {
		return float64(xival), float64(yival), true
	}
	return 0, 0, false
}
