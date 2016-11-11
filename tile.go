package tile

import (
	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/query"
)

// Tile represents an interface for generating tile data.
type Tile interface {
	Create(string, binning.TileCoord) ([]byte, error)
	Parse(map[string]interface{}) error
}
