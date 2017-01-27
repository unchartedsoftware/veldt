package veldt

import (
	"github.com/unchartedsoftware/veldt/binning"
)

// Tile represents an interface for generating tile data.
type Tile interface {
	Create(string, *binning.TileCoord, Query) ([]byte, error)
	Parse(map[string]interface{}) error
}

// TileCtor represents a function that instantiates and returns a new tile
// data type.
type TileCtor func() (Tile, error)
