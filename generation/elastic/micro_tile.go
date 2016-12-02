package elastic

import (
	"encoding/json"

	"github.com/unchartedsoftware/prism"
	"github.com/unchartedsoftware/prism/binning"
)

type MicroTile struct {
	Bivariate
	Tile
	TopHits
}

func NewMicroTile(host, port string) prism.TileCtor {
	return func() (prism.Tile, error) {
		m := &MicroTile{}
		m.Host = host
		m.Port = port
		return m, nil
	}
}

func (m *MicroTile) Parse(params map[string]interface{}) error {
	err := m.Bivariate.Parse(params)
	if err != nil {
		return nil
	}
	return m.TopHits.Parse(params)
}

func (m *MicroTile) Create(uri string, coord *binning.TileCoord, query prism.Query) ([]byte, error) {
	// get client
	client, err := NewClient(m.Host, m.Port)
	if err != nil {
		return nil, err
	}
	// create search service
	search := client.Search().
		Index(uri).
		Size(0)

	// create root query
	q, err := m.CreateQuery(query)
	if err != nil {
		return nil, err
	}
	// add tiling query
	q.Must(m.Bivariate.GetQuery(coord))
	// set the query
	search.Query(q)

	// get aggs
	aggs := m.TopHits.GetAggs()
	// set the aggregation
	search.Aggregation("top-hits", aggs["top-hits"])

	// send query
	res, err := search.Do()
	if err != nil {
		return nil, err
	}

	// get top hits
	hits, err := m.TopHits.GetTopHits(&res.Aggregations)
	if err != nil {
		return nil, err
	}

	return json.Marshal(hits)
}
