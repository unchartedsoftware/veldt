package prism

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"

	"github.com/unchartedsoftware/prism/binning"
)

// TileRequest represents a tile data generation request.
type TileRequest struct {
	URI      string
	Coord    binning.TileCoord
	Query    Query
	Tile     Tile
}

// GetHash returns a unique hash for the request.
func (r *TileRequest) GetHash() string {
	spew.Config.SortKeys = true
	return fmt.Sprintf("%s:%s", "tile", spew.Sdump(r))
}

// MetaRequest represents a meta data generation request.
type MetaRequest struct {
	URI   string
	Meta  Meta
}

// GetHash returns a unique hash for the request.
func (r *MetaRequest) GetHash() string {
	spew.Config.SortKeys = true
	return fmt.Sprintf("%s:%s", "meta", spew.Sdump(r))
}
