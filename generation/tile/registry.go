package tile

import (
	"fmt"
)

// TileGenerator represents a function which takes a tile request and returns a byte
// slice of marshalled tile data.
type TileGenerator func(tileReq *TileRequest) ([]byte, error)

// registry contains all tiling function implementations.
var (
	registry = make(map[string]TileGenerator)
)

// Register registers a tile generator under the provided type id string.
func Register(typeID string, tile TileGenerator) {
	registry[typeID] = tile
}

// GetGeneratorByType when given a string id will return the registered
// tile generator.
func GetGeneratorByType(typeID string) (TileGenerator, error) {
	gen, ok := registry[typeID]
	if !ok {
		return nil, fmt.Errorf("Tile type '%s' is not recognized", typeID)
	}
	return gen, nil
}
