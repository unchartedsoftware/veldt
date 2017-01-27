package veldt

import (
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"

	"github.com/unchartedsoftware/veldt/binning"
)

// TileRequest represents a tile data generation request.
type TileRequest struct {
	URI   string
	Coord *binning.TileCoord
	Query Query
	Tile  Tile
}

// GetHash returns a unique hash for the request.
func (r *TileRequest) GetHash() string {
	spew.Config.SortKeys = true
	spew.Config.SpewKeys = true
	spew.Config.DisablePointerAddresses = true
	spew.Config.DisableCapacities = true
	spew.Config.DisablePointerMethods = true
	spew.Config.DisableMethods = true
	hash := fmt.Sprintf("%s:%s", "tile", spew.Sdump(r))
	return strings.Join(strings.Fields(hash), ":")
}

// MetaRequest represents a meta data generation request.
type MetaRequest struct {
	URI  string
	Meta Meta
}

// GetHash returns a unique hash for the request.
func (r *MetaRequest) GetHash() string {
	spew.Config.SortKeys = true
	spew.Config.SpewKeys = true
	spew.Config.DisablePointerAddresses = true
	spew.Config.DisableCapacities = true
	spew.Config.DisablePointerMethods = true
	spew.Config.DisableMethods = true
	hash := fmt.Sprintf("%s:%s", "meta", spew.Sdump(r))
	return strings.Join(strings.Fields(hash), ":")
}
