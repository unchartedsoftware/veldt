package tile

import (
	"encoding/json"

	jsonutil "github.com/unchartedsoftware/veldt/util/json"
)

// MicroEdge represents a tile that returns individual data edges with optional
// included attributes.
type MicroEdge struct {
	LOD int
	// src
	srcXField    string
	srcYField    string
	srcXIncluded bool
	srcYIncluded bool
	// dst
	dstXField    string
	dstYField    string
	dstXIncluded bool
	dstYIncluded bool
}

// Parse parses the provided JSON object and populates the structs attributes.
func (e *MicroEdge) Parse(params map[string]interface{}) error {
	// parse LOD
	e.LOD = jsonutil.GetIntDefault(params, 0, "lod")
	return nil
}

// ParseIncludes parses the included attributes to ensure they include the raw
// data coordinates.
func (e *MicroEdge) ParseIncludes(includes []string, srcXField string, srcYField string, dstXField string, dstYField string) []string {
	// store x / y fields
	e.srcXField = srcXField
	e.srcYField = srcYField
	e.dstXField = dstXField
	e.dstYField = dstYField
	// src includes
	if !existsIn(e.srcXField, includes) {
		includes = append(includes, e.srcXField)
	} else {
		e.srcXIncluded = true
	}
	if !existsIn(e.srcYField, includes) {
		includes = append(includes, e.srcYField)
	} else {
		e.srcYIncluded = true
	}
	// dst includes
	if !existsIn(e.dstXField, includes) {
		includes = append(includes, e.dstXField)
	} else {
		e.dstXIncluded = true
	}
	if !existsIn(e.dstYField, includes) {
		includes = append(includes, e.dstYField)
	} else {
		e.dstYIncluded = true
	}
	return includes
}

// Encode will encode the tile results
func (e *MicroEdge) Encode(hits []map[string]interface{}, points []float32) ([]byte, error) {
	emptyHits := true
	// remove any non-included fields from hits
	if !e.srcXIncluded || !e.srcYIncluded ||
		!e.dstXIncluded || !e.dstYIncluded {
		for _, hit := range hits {
			// remove fields if they weren't explicitly included
			if !e.srcXIncluded {
				delete(hit, e.srcXField)
			}
			if !e.srcYIncluded {
				delete(hit, e.srcYField)
			}
			if !e.dstXIncluded {
				delete(hit, e.dstXField)
			}
			if !e.dstYIncluded {
				delete(hit, e.dstYField)
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
