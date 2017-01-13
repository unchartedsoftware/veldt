package tile

import (
	"github.com/unchartedsoftware/prism/util/json"
)

type Macro struct {
	LOD int
}

func (m *Macro) Parse(params map[string]interface{}) error {
	// parse LOD
	m.LOD = int(json.GetNumberDefault(params, 0, "lod"))
	return nil
}

func (m *Macro) Encode(points []float32) ([]byte, error) {
	// encode the results
	if m.LOD > 0 {
		return EncodeLOD(points, m.LOD), nil
	}
	return Encode(points), nil
}
