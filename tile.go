package prism

import (
	"github.com/unchartedsoftware/prism/binning"
)

// Tile represents an interface for generating tile data.
type Tile interface {
	Create(string, binning.TileCoord) ([]byte, error)
	Parse(map[string]interface{}) error
}

// TileCtor represents a function that instantiates and returns a new tile
// data type.
type TileCtor func() (Meta, error)
