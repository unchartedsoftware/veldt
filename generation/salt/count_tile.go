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
func NewCountTile (rmqConfig Configuration) veldt.TileCtor {
	return func() (veldt.Tile, error) {
		t := &Count{}
		t.rmqConfig = rmqConfig
		return t, nil
	}
}

// Create generates a tile from the provided URI, tile coordinate, and query parameters
func (t *Count) Create (uri string, coord *binning.TileCoord, query veldt.Query) ([]byte, error) {
	connection, err := NewConnection(t.rmqConfig)
	if err != nil {
		return nil, err
	}

	return connection.Query(t.rmqConfig.serverQueue, []byte("this is a test"))
}

