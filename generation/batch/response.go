package batch

// TileResponse is the structure returned by a tile factory containing either
// a finished tile or a error.  These two fields should be mutually exclusive.
type TileResponse struct {
	// Tile is the tile created from the request, if there was no error
	Tile []byte
	// Err is the error thrown as a result of trying to fulfill the request,
	// if there was one
	Err error
}
