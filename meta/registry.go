package meta

import (
	"fmt"
)

var (
	registry = make(map[string]GeneratorConstructor)
)

// Register registers a meta data generator under the provided type id string.
func Register(typeID string, gen GeneratorConstructor) {
	registry[typeID] = gen
}

// GetGenerator instantiates a meta data generator from a meta data request.
func GetGenerator(req *Request) (Generator, error) {
	ctor, ok := registry[req.Type]
	if !ok {
		return nil, fmt.Errorf("Meta type '%s' is not recognized in request %s", req.Type, req.String())
	}
	return ctor(req)
}
