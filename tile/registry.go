package tile

import (
	"fmt"
)

var (
	// registry contains all registered tile generator constructors.
	registry = make(map[string]GeneratorConstructor)
)

// Register registers a tile generator under the provided type id string.
func Register(typeID string, gen GeneratorConstructor) {
	registry[typeID] = gen
}

// GetGenerator instantiates a tile generator from a tile request.
func GetGenerator(req *Request) (Generator, error) {
	ctor, ok := registry[req.Type]
	if !ok {
		return nil, fmt.Errorf("Tile type '%s' is not recognized in request %s", req.Type, req.String())
	}
	return ctor(req)
}
