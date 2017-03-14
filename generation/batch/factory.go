package batch


import (
	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
)

// TileResponse is the structure returned by a tile factory containing either
// a finished tile or a error.  These two fields should be mutually exclusive.
type TileResponse struct {
	tile []byte
	err  error
}

// TileRequest contains all the information a tile factory needs to fulfill a
// request for a single tile
type TileRequest struct {
	// The parameters passed to our tile request for parsing
	parameters    map[string]interface{}
	// The URI to which our tile request was made 
	uri           string
	// The coordinates of the requested tile
	coordinates  *binning.TileCoord
	// The filter to apply to the data for our tile request
	query         veldt.Query
	// A channel on which the tile should be returned to us by the tile factory
	resultChannel chan TileResponse
}


// TileFactory is minimal set of functions required to construct tiles from a
// set of tile requests.
type TileFactory interface {
	Create(requests []*TileRequest)
}
