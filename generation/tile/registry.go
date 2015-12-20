package tile

import (
	"fmt"
)

var (
	// registry contains all tiling function implementations.
	registry = make(map[string]GeneratorConstructor)
)

// Param represents a single set of related tiling parameters.
type Param interface {
	GetHash() string
}

// Generator represents an interface for generating tile data.
type Generator interface {
	GetTile(*Request) ([]byte, error)
	GetParams() []Param
}

// GeneratorConstructor represents a function to instantiate a new generator
// from a tile request.
type GeneratorConstructor func(*Request) (Generator, error)

// Register registers a tile generator under the provided type id string.
func Register(typeID string, gen GeneratorConstructor) {
	registry[typeID] = gen
}

// GetGenerator instantiates a tile generator from a tile request.
func GetGenerator(tileReq *Request) (Generator, error) {
	ctor, ok := registry[tileReq.Type]
	if !ok {
		return nil, fmt.Errorf("Tile type '%s' is not recognized", tileReq.Type)
	}
	return ctor(tileReq)
}
