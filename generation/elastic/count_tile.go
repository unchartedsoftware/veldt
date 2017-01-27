package elastic

import (
	"fmt"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
)

// Count represents an elasticsearch implementation of the count tile.
type Count struct {
	Bivariate
	Tile
}

// NewCountTile instantiates and returns a new tile struct.
func NewCountTile(host, port string) veldt.TileCtor {
	return func() (veldt.Tile, error) {
		t := &Count{}
		t.Host = host
		t.Port = port
		return t, nil
	}
}

// Create generates a tile from the provided URI, tile coordinate and query
// parameters.
func (t *Count) Create(uri string, coord *binning.TileCoord, query veldt.Query) ([]byte, error) {
	// get client
	client, err := NewClient(t.Host, t.Port)
	if err != nil {
		return nil, err
	}
	// create search service
	search := client.Search().
		Index(uri).
		Size(0)

	// create root query
	q, err := t.CreateQuery(query)
	if err != nil {
		return nil, err
	}
	// add tiling query
	q.Must(t.Bivariate.GetQuery(coord))
	// set the query
	search.Query(q)

	// send query
	res, err := search.Do()
	if err != nil {
		return nil, err
	}

	return []byte(fmt.Sprintf(`{"count":%d}`, res.Hits.TotalHits)), nil
}
