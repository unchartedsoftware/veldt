package tile

import (
	"fmt"
	"encoding/json"

)

// Micro represents a tile that returns individual data points with optional
// included attributes.
type MicroPair struct {
	Micro
	X2Field    string
	Y2Field    string
	X2Included bool
	Y2Included bool
}

// ParseIncludes parses the included attributes to ensure they include the raw
// data coordinates.
func (m *MicroPair) ParseIncludes(includes []string, xField string, yField string, x2Field string, y2Field string) []string {

	fmt.Printf("<><><><><><><><><> \n")
	fmt.Printf("	<><><> passed in xField(%s) and yField(%s) \n", xField, yField)
	fmt.Printf("	<><><> tile.MicroPair pre-parse Includes (%d) = %v \n", len(includes), includes)
	// store x / y field
	includes = m.Micro.ParseIncludes(includes, xField, yField)
	fmt.Printf("	<><><> tile.MicroPair Micro.ParseIncludes (%d) = %v \n", len(includes), includes)
	m.X2Field = x2Field
	m.Y2Field = y2Field

	// ensure that the x2 / y2 field are included
	if !existsIn(x2Field, includes) {
		includes = append(includes, x2Field)
	} else {
		m.X2Included = true
	}
	if !existsIn(y2Field, includes) {
		includes = append(includes, y2Field)
	} else {
		m.Y2Included = true
	}
	fmt.Printf("	<><><> tile.MicroPair FINAL Includes (%d) = %v \n", len(includes), includes)
	fmt.Printf("<><><><><><><><><> \n")
	return includes
}

// Encode will encode the tile results based on the LOD property.
func (m *MicroPair) Encode(hits []map[string]interface{}, points []float32) ([]byte, error) {
	emptyHits := true
	// remove any non-included fields from hits
	if !m.XIncluded || !m.YIncluded || !m.X2Included || !m.Y2Included {
		for _, hit := range hits {
			// remove fields if they weren't explicitly included
			if !m.XIncluded {
				delete(hit, m.xField)
			}
			if !m.YIncluded {
				delete(hit, m.yField)
			}
			if !m.X2Included {
				delete(hit, m.X2Field)
			}
			if !m.Y2Included {
				delete(hit, m.Y2Field)
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
	//if m.LOD > 0 {
	//	// NOTE: during LOD points are sorted by morton code, therefore we sort
	//	// the hits by morton code as well to ensure both arrays align by index.
	//	sortHitsArray(hits, points, 4)
	//	// sort points and get offsets
	//	sorted, offsets := LOD(points, m.LOD)
	//	return json.Marshal(map[string]interface{}{
	//		"points":  sorted,
	//		"offsets": offsets,
	//		"hits":    hits,
	//	})
	//}
	// encode without LOD
	return json.Marshal(map[string]interface{}{
		"points": points,
		"hits":   hits,
	})
}
