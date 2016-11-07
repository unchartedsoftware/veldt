package tile

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"

	"github.com/unchartedsoftware/prism/binning"
)

func init() {
	spew.Config.SortKeys = true
}

// Request represents the tile type and tile coord.
type Request struct {
	Type string
	Tile string
	Param  param.Params
	URI    string
	Store  string
	Coord  *binning.TileCoord
	Query  query.Query
}

// GetHash returns a unique hash for the request.
func (r *Request) GetHash() string {
	return fmt.Sprintf("%s:%s", "tile", spew.Dump(r))
}
