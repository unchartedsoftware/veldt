package tile

import (
	"fmt"
)

// Generator represents a function which takes a tile request and returns a
// byte slice of marshalled tile data.
type Generator func(*Request) ([]byte, error)

// Hasher represents a function which takes a tile request and returns a
// unique string hash o.
type Hasher func(*Request) string

// Pair represents a tile generator and tile hasher pair.
type Pair struct {
	Hasher    Hasher
	Generator Generator
}

// registry contains all tiling function implementations.
var (
	registry = make(map[string]Pair)
)

// Register registers a tile generator under the provided type id string.
func Register(typeID string, tile Generator, hash Hasher) {
	registry[typeID] = Pair{
		Generator: tile,
		Hasher:    hash,
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

// GetHasherByType when given a string id will return the registered
// tile generator.
func GetHasherByType(typeID string) (Hasher, error) {
	pair, ok := registry[typeID]
	if !ok {
		return nil, fmt.Errorf("Tile type '%s' is not recognized", typeID)
	}
	return pair.Hasher, nil
}
