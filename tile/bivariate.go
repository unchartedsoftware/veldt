package tile

import (
	"fmt"
	"strings"

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

// TileBounds computes and returns the tile bounds for the provided tile coord.
func (b *Bivariate) TileBounds(coord *binning.TileCoord) *geometry.Bounds {
	if b.tileBounds == nil {
		b.tileBounds = binning.GetTileBounds(coord, b.globalBounds)
	}
	return b.tileBounds
}

// BinSizeX computes and returns the size of a bin across the x axis for the
// provided tile coord.
func (b *Bivariate) BinSizeX(coord *binning.TileCoord) float64 {
	return b.TileBounds(coord).RangeX() / float64(b.Resolution)
}

// BinSizeY computes and returns the size of a bin across the x axis for the
// provided tile coord.
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
	rang := bounds.RangeX()
	if bounds.Left > bounds.Right {
		return binning.MaxTileResolution - (((x - bounds.Right) / rang) * binning.MaxTileResolution)
	}
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
	rang := bounds.RangeY()
	if bounds.Bottom > bounds.Top {
		return binning.MaxTileResolution - (((y - bounds.Top) / rang) * binning.MaxTileResolution)
	}
	return ((y - bounds.Bottom) / rang) * binning.MaxTileResolution
}

// GetXY given a data hit, returns the corresponding coord within the range of
// [0 : 256) for the tile.
func (b *Bivariate) GetXY(coord *binning.TileCoord, hit map[string]interface{}) (float64, float64, bool) {
	// get X / Y of the data
	x, y, ok := b.getPixel(hit)
	if !ok {
		return 0, 0, false
	}
	// convert to tile pixel coords in the range [0 - 256)
	tx := b.GetX(coord, x)
	ty := b.GetY(coord, y)
	// return position in tile coords
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

func getInterface(arg map[string]interface{}, path ...string) (interface{}, bool) {
	child := arg
	last := len(path) - 1
	var val interface{} = child
	for index, key := range path {
		// does a child exists?
		v, ok := child[key]
		if !ok {
			return nil, false
		}
		// is it the target?
		if index == last {
			val = v
			break
		}
		// if not, does it have children to traverse?
		c, ok := v.(map[string]interface{})
		if !ok {
			return nil, false
		}
		child = c
	}
	return val, true
}

func getFloat64(arg map[string]interface{}, path ...string) (float64, bool) {
	i, ok := getInterface(arg, path...)
	if !ok {
		return 0, false
	}
	switch v := i.(type) {
	case float32:
		return float64(v), true
	case float64:
		return v, true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	}
	return 0, false
}

func (b *Bivariate) getPixel(hit map[string]interface{}) (float64, float64, bool) {
	xPath := strings.Split(b.XField, ".")
	yPath := strings.Split(b.YField, ".")
	x, ok := getFloat64(hit, xPath...)
	if !ok {
		return 0, 0, false
	}
	y, ok := getFloat64(hit, yPath...)
	if !ok {
		return 0, 0, false
	}
	return x, y, true
}
