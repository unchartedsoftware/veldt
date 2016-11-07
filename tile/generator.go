package tile

import (
	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/query"
)

// Generator represents an interface for generating tile data.
type Generator interface {
	GetTile(string) ([]byte, error)
}

// GeneratorConstructor represents a function to instantiate a new generator
// from a tile request.
type GeneratorConstructor func(*binning.TileCoord, map[string]interface{}) (Generator, error)
