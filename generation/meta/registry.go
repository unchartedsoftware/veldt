package meta

import (
	"fmt"
)

// MetaGenerator represents a function which takes a meta request and returns a byte
// slice of marshalled meta data.
type MetaGenerator func(metaReq *MetaRequest) ([]byte, error)

// registry contains all tiling function implementations.
var (
	registry = make(map[string]MetaGenerator)
)

// Register registers a meta generator under the provided type id string.
func Register(typeID string, meta MetaGenerator) {
	registry[typeID] = meta
}

// GetGeneratorByType when given a string id will return the registered
// meta generator.
func GetGeneratorByType(typeID string) (MetaGenerator, error) {
	gen, ok := registry[typeID]
	if !ok {
		return nil, fmt.Errorf("Meta type '%s' is not recognized", typeID)
	}
	return gen, nil
}
