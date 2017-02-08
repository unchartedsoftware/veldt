package tile

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/unchartedsoftware/veldt/binning"
	"github.com/unchartedsoftware/veldt/geometry"
	jsonUtil "github.com/unchartedsoftware/veldt/util/json"
)

// Edge represents a tile that returns individual data edges with optional
// included attributes.
type Edge struct {
	LOD int
	// src
	SrcXField      string
	SrcYField      string
	isSrcXIncluded bool
	isSrcYIncluded bool
	// dst
	DstXField      string
	DstYField      string
	isDstXIncluded bool
	isDstYIncluded bool
	// Bounds
	WorldBounds *geometry.Bounds
	TileBounds  *geometry.Bounds
}

// Parse parses the provided JSON object and populates the structs attributes.
func (e *Edge) Parse(params map[string]interface{}) error {
	// src x and y fields
	srcXField, ok := jsonUtil.GetString(params, "srcXField")
	if !ok {
		return fmt.Errorf("`srcXField` parameter missing from tile")
	}
	srcYField, ok := jsonUtil.GetString(params, "srcYField")
	if !ok {
		return fmt.Errorf("`srcYField` parameter missing from tile")
	}
	// dst x and y fields
	dstXField, ok := jsonUtil.GetString(params, "dstXField")
	if !ok {
		return fmt.Errorf("`dstXField` parameter missing from tile")
	}
	dstYField, ok := jsonUtil.GetString(params, "dstYField")
	if !ok {
		return fmt.Errorf("`dstYField` parameter missing from tile")
	}
	worldBounds, err := geometry.NewBoundsByParse(params)
	if err != nil {
		return fmt.Errorf("left, right, bottom or top missing from tile")
	}

	// set attributes
	e.WorldBounds = worldBounds
	e.SrcXField = srcXField
	e.SrcYField = srcYField
	e.DstXField = dstXField
	e.DstYField = dstYField

	return nil
}

// ParseIncludes parses the included attributes to ensure they include the raw
// data coordinates.
func (e *Edge) ParseIncludes(includes []string) []string {
	// src includes
	if !existsIn(e.SrcXField, includes) {
		includes = append(includes, e.SrcXField)
	}
	if !existsIn(e.SrcYField, includes) {
		includes = append(includes, e.SrcYField)
	}
	// dst includes
	if !existsIn(e.DstXField, includes) {
		includes = append(includes, e.DstXField)
	}
	if !existsIn(e.DstYField, includes) {
		includes = append(includes, e.DstYField)
	}
	return includes
}

// GetX given an x value, returns the corresponding coord within the range of
// [0 : 2^zoom * 256) for the tile.
func (e *Edge) GetX(x float64, zoom uint32) float64 {
	extent := binning.MaxTileResolution * math.Pow(2, float64(zoom))
	bounds := e.TileBounds
	if bounds.BottomLeft().X > bounds.TopRight().X {
		rang := bounds.BottomLeft().X - bounds.TopRight().X
		return extent - (((x - bounds.TopRight().X) / rang) * extent)
	}
	rang := bounds.TopRight().X - bounds.BottomLeft().X
	return ((x - bounds.BottomLeft().X) / rang) * extent
}

// GetY given an y value, returns the corresponding coord within the range of
// [0 : 2^zoom * 256) for the tile.
func (e *Edge) GetY(y float64, zoom uint32) float64 {
	extent := binning.MaxTileResolution * math.Pow(2, float64(zoom))
	bounds := e.TileBounds
	if bounds.BottomLeft().Y > bounds.TopRight().Y {
		rang := bounds.BottomLeft().Y - bounds.TopRight().Y
		return extent - (((y - bounds.TopRight().Y) / rang) * extent)
	}
	rang := bounds.TopRight().Y - bounds.BottomLeft().Y
	return ((y - bounds.BottomLeft().Y) / rang) * extent
}

// GetSrcXY given a data hit, returns the corresponding coord within the range of
// [0 : 2^zoom * 256) for the tile.
func (e *Edge) getXY(hit map[string]interface{}, xField string, yField string, zoom uint32) (float32, float32, bool) {
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
	tx := e.GetX(x, zoom)
	ty := e.GetY(y, zoom)
	// return position in tile coords with 2 decimal places
	return toFixed(float32(tx), 2), toFixed(float32(ty), 2), true
}

// GetSrcXY given a data hit, returns the corresponding coord within the range of
// [0 : 2^zoom * 256) for the tile.
func (e *Edge) GetSrcXY(hit map[string]interface{}, zoom uint32) (float32, float32, bool) {
	return e.getXY(hit, e.SrcXField, e.SrcYField, zoom)
}

// GetDstXY given a data hit, returns the corresponding coord within the range of
// [0 : 2^zoom * 256) for the tile.
func (e *Edge) GetDstXY(hit map[string]interface{}, zoom uint32) (float32, float32, bool) {
	return e.getXY(hit, e.DstXField, e.DstYField, zoom)
}

// Encode will encode the tile results
func (e *Edge) Encode(hits []map[string]interface{}, points []float32) ([]byte, error) {
	emptyHits := true
	// remove any non-included fields from hits
	if !e.isSrcXIncluded || !e.isSrcYIncluded {
		for _, hit := range hits {
			// remove fields if they weren't explicitly included
			if !e.isSrcXIncluded {
				delete(hit, e.SrcXField)
			}
			if !e.isSrcYIncluded {
				delete(hit, e.SrcYField)
			}
			if !e.isDstXIncluded {
				delete(hit, e.DstXField)
			}
			if !e.isDstYIncluded {
				delete(hit, e.DstYField)
			}
			if emptyHits && len(hit) > 0 {
				emptyHits = false
			}
		}
	}

	// if no hit contains any data, occlude them from response
	if emptyHits {
		// no point returning an array of empty hits
		hits = nil
	}

	// encode using LOD
	if e.LOD > 0 {
		// NOTE: during LOD points are sorted by morton code, therefore we sort
		// the hits by morton code as well to ensure both arrays align by index.
		sortHitsArray(hits, points)
		// sort points and get offsets
		sorted, offsets := LOD(points, e.LOD)
		return json.Marshal(map[string]interface{}{
			"points":  sorted,
			"offsets": offsets,
			"hits":    hits,
		})
	}
	// encode without LOD
	return json.Marshal(map[string]interface{}{
		"points": points,
		"hits":   hits,
	})
}
