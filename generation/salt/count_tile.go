package salt

import (
	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/tile"
	"github.com/unchartedsoftware/veldt/binning"
)

// Count represents a salt implementation of a simple counting tiling scheme
type Count struct {
	tile.Bivariate
	Tile
}

// NewCountTile returns a salt-based data-counting tile
func NewCountTile (host, port string) veldt.TileCtor {
	return func() (veldt.Tile, error) {
		t := &Count{}
		t.host = host
		t.port = port
		return t, nil
	}
}

// Create generates a tile from the provided URI, tile coordinate, and query parameters
func (t *Count) Create (uri string, coord *binning.TileCoord, query veldt.Query) ([]byte, error) {
	return []byte("this is a test"), nil
}

