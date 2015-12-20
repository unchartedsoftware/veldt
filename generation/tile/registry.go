package tile

import (
	"fmt"
)

// Generator represents a function which takes a tile request and returns a
// byte slice of marshalled tile data.
type Generator func(*Request, map[string]Param) ([]byte, error)

// Param represents a single set of related tiling parameters.
type Param interface {
	GetHash() string
}

// Params represents a function which takes a tile request and returns a
// map of parameters.
type Params func(*Request) map[string]Param

// Pair represents a tile generator and tile hasher pair.
type Pair struct {
	Params    Params
	Generator Generator
}

// registry contains all tiling function implementations.
var (
	registry = make(map[string]Pair)
)

// Register registers a tile generator under the provided type id string.
func Register(typeID string, tile Generator, params Params) {
	registry[typeID] = Pair{
		Generator: tile,
		Params:    params,
	}
}

// GetGeneratorByType when given a string id will return the registered
// tile generator.
func GetGeneratorByType(typeID string) (Generator, error) {
	pair, ok := registry[typeID]
	if !ok {
		return nil, fmt.Errorf("Tile type '%s' is not recognized", typeID)
	}
	return pair.Generator, nil
}

// GetParamsByType when given a string id will return the registered
// tile parameters.
func GetParamsByType(typeID string) (Params, error) {
	pair, ok := registry[typeID]
	if !ok {
		return nil, fmt.Errorf("Parameters type '%s' is not recognized", typeID)
	}
	return pair.Params, nil
}
