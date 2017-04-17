package tile

import (
	"github.com/unchartedsoftware/veldt/util/json"
)

// MacroEdge represents a tile that returns individual data edges with optional
// included attributes.
type MacroEdge struct {
	LOD int
}

// Parse parses the provided JSON object and populates the structs attributes.
func (e *MacroEdge) Parse(params map[string]interface{}) error {
	// parse LOD
	e.LOD = json.GetIntDefault(params, 0, "lod")
	return nil
}

// ParseIncludes parses the included attributes to ensure they include the raw
// data coordinates.
func (e *MacroEdge) ParseIncludes(includes []string, srcXField string, srcYField string, dstXField string, dstYField string, weightField string) []string {
	// src includes
	if !existsIn(srcXField, includes) {
		includes = append(includes, srcXField)
	}
	if !existsIn(srcYField, includes) {
		includes = append(includes, srcYField)
	}
	// dst includes
	if !existsIn(dstXField, includes) {
		includes = append(includes, dstXField)
	}
	if !existsIn(dstYField, includes) {
		includes = append(includes, dstYField)
	}
	// weight includes
	if !existsIn(weightField, includes) {
		includes = append(includes, weightField)
	}
	return includes
}

// Encode will encode the tile results
func (e *MacroEdge) Encode(edges []float32) ([]byte, error) {
	// encode the results
	if e.LOD > 0 {
		return EncodeEdgeLOD(edges, e.LOD), nil
	}
	return EncodeFloat32(edges), nil
}
