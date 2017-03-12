package veldt

import (
	"github.com/unchartedsoftware/veldt/binning"
)

// Tile represents an interface for generating tile data.
type Tile interface {
	// Parse parses a tile request for future creation
	// parameter 1 (string): the name under which this tile constructor is registered
	//             in the pipeline
	// parameter 2 (map[string]interface{}) Any parameters specifying the how this tile
	//             is to be created
	Parse(string, map[string]interface{}) error
	// Create creates a tile.
	// parameter 1 (string): A dataset ID (typically called uri)
	// parameter 2 (*binning.TileCoord): the coordinates of the requested tile
	// parameter 3 (Query): A query to specify which data should be included in the
	//             tile - essentially a filter
	Create(string, *binning.TileCoord, Query) ([]byte, error)
}

// TileCtor represents a function that instantiates and returns a new tile
// data type.
type TileCtor func() (Tile, error)
