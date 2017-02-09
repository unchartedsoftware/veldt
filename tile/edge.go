package tile

import (
	"fmt"

	"github.com/unchartedsoftware/veldt/binning"
	"github.com/unchartedsoftware/veldt/geometry"
	"github.com/unchartedsoftware/veldt/util/json"
)

// Edge represents a tile that returns individual data edges with optional
// included attributes.
type Edge struct {
	// src
	SrcXField string
	SrcYField string
	// dst
	DstXField string
	DstYField string
	// Bounds
	tileBounds   *geometry.Bounds
	globalBounds *geometry.Bounds
}

// Parse parses the provided JSON object and populates the structs attributes.
func (e *Edge) Parse(params map[string]interface{}) error {
	// src x and y fields
	srcXField, ok := json.GetString(params, "srcXField")
	if !ok {
		return fmt.Errorf("`srcXField` parameter missing from tile")
	}
	srcYField, ok := json.GetString(params, "srcYField")
	if !ok {
		return fmt.Errorf("`srcYField` parameter missing from tile")
	}
	// dst x and y fields
	dstXField, ok := json.GetString(params, "dstXField")
	if !ok {
		return fmt.Errorf("`dstXField` parameter missing from tile")
	}
	dstYField, ok := json.GetString(params, "dstYField")
	if !ok {
		return fmt.Errorf("`dstYField` parameter missing from tile")
	}
	// set attributes
	e.SrcXField = srcXField
	e.SrcYField = srcYField
	e.DstXField = dstXField
	e.DstYField = dstYField
	// clear tile bounds
	e.tileBounds = nil
	// get the global bounds
	e.globalBounds = &geometry.Bounds{}
	return e.globalBounds.Parse(params)
}

// TileBounds computes and returns the tile bounds for the provided tile coord.
func (e *Edge) TileBounds(coord *binning.TileCoord) *geometry.Bounds {
	if e.tileBounds == nil {
		e.tileBounds = binning.GetTileBounds(coord, e.globalBounds)
	}
	return e.tileBounds
}

// GetX given an x value, returns the corresponding coord within the range of
// [0 : 256) for the tile.
func (e *Edge) GetX(coord *binning.TileCoord, x float64) float64 {
	bounds := e.TileBounds(coord)
	rang := bounds.RangeX()
	if bounds.Left > bounds.Right {
		return binning.MaxTileResolution - (((x - bounds.Right) / rang) * binning.MaxTileResolution)
	}
	return ((x - bounds.Left) / rang) * binning.MaxTileResolution
}

// GetY given an y value, returns the corresponding coord within the range of
// [0 : 256) for the tile.
func (e *Edge) GetY(coord *binning.TileCoord, y float64) float64 {
	bounds := e.TileBounds(coord)
	rang := bounds.RangeY()
	if bounds.Bottom > bounds.Top {
		return binning.MaxTileResolution - (((y - bounds.Top) / rang) * binning.MaxTileResolution)
	}
	return ((y - bounds.Bottom) / rang) * binning.MaxTileResolution
}

// GetSrcXY given a data hit, returns the corresponding coord within the range of
// [0 : 2^zoom * 256) for the tile.
func (e *Edge) getXY(coord *binning.TileCoord, hit map[string]interface{}, xField string, yField string) (float64, float64, bool) {
	// get x / y fields from data
	ix, ok := hit[xField]
	if !ok {
		return 0, 0, false
	}
	iy, ok := hit[yField]
	if !ok {
		return 0, 0, false
	}
	// get X / Y of the data
	x, y, ok := castPixel(ix, iy)
	if !ok {
		return 0, 0, false
	}
	// convert to tile pixel coords in the range [0 : 2^zoom * 256)
	tx := e.GetX(coord, x)
	ty := e.GetY(coord, y)
	// return position in tile coords with 2 decimal places
	return tx, ty, true
}

// GetSrcXY given a data hit, returns the corresponding coord within the range of
// [0 : 2^zoom * 256) for the tile.
func (e *Edge) GetSrcXY(coord *binning.TileCoord, hit map[string]interface{}) (float64, float64, bool) {
	return e.getXY(coord, hit, e.SrcXField, e.SrcYField)
}

// GetDstXY given a data hit, returns the corresponding coord within the range of
// [0 : 2^zoom * 256) for the tile.
func (e *Edge) GetDstXY(coord *binning.TileCoord, hit map[string]interface{}) (float64, float64, bool) {
	return e.getXY(coord, hit, e.DstXField, e.DstYField)
}
