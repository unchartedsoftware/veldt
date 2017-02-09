package tile

import (
	"fmt"

	"github.com/unchartedsoftware/veldt/binning"
	"github.com/unchartedsoftware/veldt/geometry"
	"github.com/unchartedsoftware/veldt/util/json"
)

const (
	numDecimals = 2
)

// Bivariate represents the parameters required for any bivariate tile.
type Bivariate struct {
	XField       string
	YField       string
	Resolution   int
	tileBounds   *geometry.Bounds
	globalBounds *geometry.Bounds
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
	// get resolution
	resolution := json.GetIntDefault(params, 256, "resolution")
	// set attributes
	b.XField = xField
	b.YField = yField
	b.Resolution = resolution
	// clear tile bounds
	b.tileBounds = nil
	// get the global bounds
	b.globalBounds = &geometry.Bounds{}
	return b.globalBounds.Parse(params)
}

func (b *Bivariate) TileBounds(coord *binning.TileCoord) *geometry.Bounds {
	if b.tileBounds == nil {
		b.tileBounds = binning.GetTileBounds(coord, b.globalBounds)
	}
	return b.tileBounds
}

func (b *Bivariate) BinSizeX(coord *binning.TileCoord) float64 {
	return b.TileBounds(coord).RangeX() / float64(b.Resolution)
}

func (b *Bivariate) BinSizeY(coord *binning.TileCoord) float64 {
	return b.TileBounds(coord).RangeY() / float64(b.Resolution)
}

// GetXBin given an x value, returns the corresponding bin.
func (b *Bivariate) GetXBin(coord *binning.TileCoord, x float64) int {
	bounds := b.TileBounds(coord)
	binSize := b.BinSizeX(coord)
	var bin int64
	if bounds.Left > bounds.Right {
		bin = int64(float64(b.Resolution-1) - ((x - bounds.Right) / binSize))
	} else {
		bin = int64((x - bounds.Left) / binSize)
	}
	return b.clampBin(bin)
}

// GetX given an x value, returns the corresponding coord within the range of
// [0 : 256) for the tile.
func (b *Bivariate) GetX(coord *binning.TileCoord, x float64) float64 {
	bounds := b.TileBounds(coord)
	if bounds.Left > bounds.Right {
		rang := bounds.Left - bounds.Right
		return binning.MaxTileResolution - (((x - bounds.Right) / rang) * binning.MaxTileResolution)
	}
	rang := bounds.Right - bounds.Left
	return ((x - bounds.Left) / rang) * binning.MaxTileResolution
}

// GetYBin given a y value, returns the corresponding bin.
func (b *Bivariate) GetYBin(coord *binning.TileCoord, y float64) int {
	bounds := b.TileBounds(coord)
	binSize := b.BinSizeY(coord)
	var bin int64
	if bounds.Bottom > bounds.Top {
		bin = int64(float64(b.Resolution-1) - ((y - bounds.Top) / binSize))
	} else {
		bin = int64((y - bounds.Bottom) / binSize)
	}
	return b.clampBin(bin)
}

// GetY given an y value, returns the corresponding coord within the range of
// [0 : 256) for the tile.
func (b *Bivariate) GetY(coord *binning.TileCoord, y float64) float64 {
	bounds := b.TileBounds(coord)
	if bounds.Bottom > bounds.Top {
		rang := bounds.Bottom - bounds.Top
		return binning.MaxTileResolution - (((y - bounds.Top) / rang) * binning.MaxTileResolution)
	}
	rang := bounds.Top - bounds.Bottom
	return ((y - bounds.Bottom) / rang) * binning.MaxTileResolution
}

// GetXY given a data hit, returns the corresponding coord within the range of
// [0 : 256) for the tile.
func (b *Bivariate) GetXY(coord *binning.TileCoord, hit map[string]interface{}) (float64, float64, bool) {
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
	tx := b.GetX(coord, x)
	ty := b.GetY(coord, y)
	// return position in tile coords with 2 decimal places
	return tx, ty, true
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
