package batch


import (
	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
)

// TileResponse is the structure returned by a tile factory containing either
// a finished tile or a error.  These two fields should be mutually exclusive.
type TileResponse struct {
	// Tile is the tile created from the request, if there was no error
	Tile []byte
	// Err is the error thrown as a result of trying to fulfill the request,
	// if there was one
	Err  error
}

// TileRequest contains all the information a tile factory needs to fulfill a
// request for a single tile
type TileRequest struct {
	// The parameters passed to our tile request for parsing
	Parameters    map[string]interface{}
	// The URI to which our tile request was made 
	URI           string
	// The coordinates of the requested tile
	Coordinates  *binning.TileCoord
	// The filter to apply to the data for our tile request
	Query         veldt.Query
	// A channel on which the tile should be returned to us by the tile factory
	ResultChannel chan TileResponse
}


// TileFactory represents an interface for generating data for multiple tiles
// simultaneously
type TileFactory interface {
	// CreateTiles creates tiles for the given tile requests.  Tiles or errors
	// should be returned individually on the channels in each TileRequest, and
	// must be returned to every request listed.
	CreateTiles(requests []*TileRequest)
}
// TileFactoryCtor represents a function that instantiates and returns a new
// tile factory type
type TileFactoryCtor func () (TileFactory, error)
