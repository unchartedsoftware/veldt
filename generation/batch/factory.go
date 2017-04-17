package batch

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
type TileFactoryCtor func() (TileFactory, error)
