package meta

import (
	"fmt"
)

// Generator represents a function which takes a meta request and returns a byte
// slice of marshalled meta data.
type Generator func(metaReq *Request) ([]byte, error)

// registry contains all tiling function implementations.
var (
	registry = make(map[string]Generator)
)

// Register registers a meta generator under the provided type id string.
func Register(typeID string, meta Generator) {
	registry[typeID] = meta
}

// GetGeneratorByType when given a string id will return the registered
// meta generator.
func GetGeneratorByType(typeID string) (Generator, error) {
	gen, ok := registry[typeID]
	if !ok {
		return nil, fmt.Errorf("Meta type '%s' is not recognized", typeID)
	}
	return gen, nil
}
