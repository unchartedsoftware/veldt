package tiling

import (
	"fmt"
)

// Generator represents a function which takes a tile request and returns a byte
// slice of marshalled tile data.
type Generator func(tileReq *TileRequest) ([]byte, error)

// registry contains all tiling function implementations.
var registry = make(map[string]Generator)

// Register registers a tile generator under the provided type id string.
func Register(typeID string, gen Generator) {
	registry[typeID] = gen
}

// GetGeneratorByType when given a string id will return the registered
// tile generator.
func GetGeneratorByType(typeID string) (Generator, error) {
	gen, ok := registry[typeID]
	if !ok {
		return nil, fmt.Errorf("Tiling type '%s' is not recognized", typeID)
	}
	return gen, nil
}
