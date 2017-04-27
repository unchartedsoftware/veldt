package veldt

import (
	"strings"

	"github.com/davecgh/go-spew/spew"

	"github.com/unchartedsoftware/veldt/binning"
)

var (
	spewer *spew.ConfigState
)

func init() {
	spewer = &spew.ConfigState{
		Indent:                  "",
		MaxDepth:                0,
		DisableMethods:          true,
		DisablePointerMethods:   true,
		DisablePointerAddresses: true,
		DisableCapacities:       true,
		ContinueOnMethod:        false,
		SortKeys:                true,
		SpewKeys:                true,
	}
}

// Request represents a basic request interface.
type Request interface {
	Create() ([]byte, error)
	GetHash() string
}

// TileRequest represents a tile data generation request.
type TileRequest struct {
	URI   string
	Coord *binning.TileCoord
	Query Query
	Tile  Tile
}

// Create generates and returns the tile for the request.
func (r *TileRequest) Create() ([]byte, error) {
	return r.Tile.Create(r.URI, r.Coord, r.Query)
}

// GetHash returns a unique hash for the request.
func (r *TileRequest) GetHash() string {
	return strings.Join(strings.Fields(spewer.Sdump(r)), "")
}

// MetaRequest represents a meta data generation request.
type MetaRequest struct {
	URI  string
	Meta Meta
}

// Create generates and returns the meta data for the request.
func (r *MetaRequest) Create() ([]byte, error) {
	return r.Meta.Create(r.URI)
}

// GetHash returns a unique hash for the request.
func (r *MetaRequest) GetHash() string {
	return strings.Join(strings.Fields(spewer.Sdump(r)), "")
}
