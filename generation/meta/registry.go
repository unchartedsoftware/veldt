package meta

import (
	"fmt"
)

var (
	// registry contains all registered meta data generator constructors.
	registry = make(map[string]GeneratorConstructor)
)

// Generator represents an interface for generating meta data.
type Generator interface {
	GetMeta(*Request) ([]byte, error)
}

// GeneratorConstructor represents a function to instantiate a new generator
// from a meta data request.
type GeneratorConstructor func(*Request) (Generator, error)

// Register registers a meta data generator under the provided type id string.
func Register(typeID string, gen GeneratorConstructor) {
	registry[typeID] = gen
}

// GetGenerator instantiates a meta data generator from a meta data request.
func GetGenerator(metaReq *Request) (Generator, error) {
	ctor, ok := registry[metaReq.Type]
	if !ok {
		return nil, fmt.Errorf("Meta type '%s' is not recognized in request %s", metaReq.Type, metaReq.String())
	}
	return ctor(metaReq)
}
