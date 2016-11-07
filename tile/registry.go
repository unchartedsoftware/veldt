package tile

import (
	"fmt"

	"github.com/unchartedsoftware/prism/binning"
)

var (
	registry = make(map[string]GeneratorConstructor)
)

// Register registers a tile generator under the provided type id string.
func Register(typeID string, gen GeneratorConstructor) {
	registry[typeID] = gen
}

// GetGenerator instantiates a tile generator from a tile request.
func GetGenerator(typeID string, coord *binning.TileCoord, params map[string]interface{}) (Generator, error) {
	ctor, ok := registry[typeID]
	if !ok {
		return nil, fmt.Errorf("query `%s` has not been registered",
			typeID)
	}
	return ctor(coord, params)
}
