package batch

import (
	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
)

// TileRequest contains all the information a tile factory needs to fulfill a
// request for a single tile
type TileRequest struct {
	// The parameters passed to our tile request for parsing
	Params map[string]interface{}
	// The URI to which our tile request was made
	URI string
	// The coordinates of the requested tile
	Coord *binning.TileCoord
	// The filter to apply to the data for our tile request
	Query veldt.Query
	// A channel on which the tile should be returned to us by the tile factory
	ResultChannel chan TileResponse
}
