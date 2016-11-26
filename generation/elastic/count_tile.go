package elastic

import (
	"fmt"

	"github.com/unchartedsoftware/prism"
	"github.com/unchartedsoftware/prism/binning"
)

type Count struct {
	Bivariate
	Tile
}

func NewCountTile(host, port string) prism.TileCtor {
	return func() (prism.Tile, error) {
		t := &Count{}
		t.Host = host
		t.Port = port
		return t, nil
	}
}

func (t *Count) Create(uri string, coord *binning.TileCoord, query prism.Query) ([]byte, error) {
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

	return []byte(fmt.Sprintf("{\"count\":%d}\n", res.Hits.TotalHits)), nil
}
