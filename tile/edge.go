package tile

import (
	"fmt"
	"strings"

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
	// requirements
	RequireSrc bool
	RequireDst bool
	// weight
	WeightField string
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
	// determine which points are required
	requireSrc := json.GetBoolDefault(params, true, "requireSrc")
	requireDst := json.GetBoolDefault(params, false, "requireDst")
	if !requireSrc && !requireDst {
		return fmt.Errorf("both `requireSrc` and `requireDst` set to false")
	}
	// get weight field
	weightField, ok := json.GetString(params, "weightField")
	if !ok {
		return fmt.Errorf("`weightField` parameter missing from tile")
	}
	// set attributes
	e.SrcXField = srcXField
	e.SrcYField = srcYField
	e.DstXField = dstXField
	e.DstYField = dstYField
	e.RequireSrc = requireSrc
	e.RequireDst = requireDst
	e.WeightField = weightField
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
// [0 : 256) for the tile.
func (e *Edge) GetSrcXY(coord *binning.TileCoord, hit map[string]interface{}) (float64, float64, bool) {
	return e.getXY(coord, hit, e.SrcXField, e.SrcYField)
}

// GetDstXY given a data hit, returns the corresponding coord within the range of
// [0 : 256) for the tile.
func (e *Edge) GetDstXY(coord *binning.TileCoord, hit map[string]interface{}) (float64, float64, bool) {
	return e.getXY(coord, hit, e.DstXField, e.DstYField)
}

// GetWeight given a data hit, returns the corresponding edge weight.
func (e *Edge) GetWeight(hit map[string]interface{}) (float64, bool) {
	weightPath := strings.Split(e.WeightField, ".")
	return getFloat64(hit, weightPath...)
}

// GetSrcXY given a data hit, returns the corresponding coord within the range of
// [0 : 256) for the tile.
func (e *Edge) getXY(coord *binning.TileCoord, hit map[string]interface{}, xField string, yField string) (float64, float64, bool) {
	// get X / Y of the data
	x, y, ok := e.getPixel(hit, xField, yField)
	if !ok {
		return 0, 0, false
	}
	// convert to tile pixel coords in the range [0 : 256)
	tx := e.GetX(coord, x)
	ty := e.GetY(coord, y)
	// return position in tile coords
	return tx, ty, true
}

func (e *Edge) getPixel(hit map[string]interface{}, xField string, yField string) (float64, float64, bool) {
	xPath := strings.Split(xField, ".")
	yPath := strings.Split(yField, ".")
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
