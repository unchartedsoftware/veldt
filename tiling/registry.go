package tiling

import (
	"fmt"
)

// Generator represents a function which takes a tile request and returns a byte
// slice of marshalled tile data.
type Generator func(tileReq *TileRequest) ([]byte, error)

// registry contains all tiling function implementations.
var registry = make(map[string]Generator)

// Register registers a tiling function under the provided type id string.
func Register(typeID string, tileFunc Generator) {
	registry[typeID] = tileFunc
}

// GetTilingFuncByType when given a tiling type id will return the registered
// tiling function.
func GetTilingFuncByType(typeID string) (Generator, error) {
	tileFunc, ok := registry[typeID]
	if !ok {
		return nil, fmt.Errorf("Tiling type '%s' is not recognized", typeID)
	}
	return tileFunc, nil
}
