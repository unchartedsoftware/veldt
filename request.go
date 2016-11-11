package tile

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"

	"github.com/unchartedsoftware/prism/binning"
)

// TileRequest represents a tile data generation request.
type TileRequest struct {
	URI      string
	Coord    binning.TileCoord
	Query    prism.Query
	Tile     prism.Tile
}

// GetHash returns a unique hash for the request.
func (r *Request) GetHash() string {
	spew.Config.SortKeys = true
	return fmt.Sprintf("%s:%s", "tile", spew.Dump(r))
}

// MetaRequest represents a meta data generation request.
type MetaRequest struct {
	URI   string
	Meta  prism.Meta
}

// GetHash returns a unique hash for the request.
func (r *Request) GetHash() string {
	spew.Config.SortKeys = true
	return fmt.Sprintf("%s:%s", "meta", spew.Dump(r))
}
