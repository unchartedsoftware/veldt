package tiling

import (
	"fmt"
)

// TileGenerator represents a function which takes a tile request and returns a byte
// slice of marshalled tile data.
type TileGenerator func(tileReq *TileRequest) ([]byte, error)

// MetaGenerator represents a function which takes a meta request and returns a byte
// slice of marshalled meta data.
type MetaGenerator func(metaReq *MetaRequest) ([]byte, error)

// GeneratorPair represents both the tile and meta data generators for a particular
// type.
type GeneratorPair struct {
	Tile TileGenerator
	Meta MetaGenerator
}

// registry contains all tiling function implementations.
var (
	registry = make(map[string]GeneratorPair)
)

// Register registers a tile generator under the provided type id string.
func Register(typeID string, tile TileGenerator, meta MetaGenerator) {
	registry[typeID] = GeneratorPair{
		Tile: tile,
		Meta: meta,
	}
}

// GetTileGeneratorByType when given a string id will return the registered
// tile generator.
func GetTileGeneratorByType(typeID string) (TileGenerator, error) {
	gen, ok := registry[typeID]
	if !ok {
		return nil, fmt.Errorf("Tile type '%s' is not recognized", typeID)
	}
	return gen.Tile, nil
}

// GetTileGeneratorByType when given a string id will return the registered
// meta generator.
func GetMetaGeneratorByType(typeID string) (MetaGenerator, error) {
	gen, ok := registry[typeID]
	if !ok {
		return nil, fmt.Errorf("Meta type '%s' is not recognized", typeID)
	}
	return gen.Meta, nil
}
