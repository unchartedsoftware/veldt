package tile

import (
	"github.com/unchartedsoftware/veldt/util/json"
)

// Macro represents a tile which returns a point for any bin that contains a
// data point.
type Macro struct {
	LOD int
}

// Parse parses the provided JSON object and populates the structs attributes.
func (m *Macro) Parse(params map[string]interface{}) error {
	// parse LOD
	m.LOD = int(json.GetNumberDefault(params, 0, "lod"))
	return nil
}

// Encode will encode the tile results based on the LOD property.
func (m *Macro) Encode(points []float32) ([]byte, error) {
	// encode the results
	if m.LOD > 0 {
		return EncodeLOD(points, m.LOD), nil
	}
	return Encode(points), nil
}
