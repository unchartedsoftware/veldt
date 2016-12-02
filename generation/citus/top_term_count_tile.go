package citus

import (
	"encoding/json"

	"github.com/unchartedsoftware/prism"
	"github.com/unchartedsoftware/prism/binning"
)

type TopTermCount struct {
	Bivariate
	TopTerms
	Tile
}

func NewTopTermCountTile(host, port string) prism.TileCtor {
	return func() (prism.Tile, error) {
		t := &TopTermCount{}
		t.Host = host
		t.Port = port
		return t, nil
	}
}

func (t *TopTermCount) Parse(params map[string]interface{}) error {
	err := t.Bivariate.Parse(params)
	if err != nil {
		return nil
	}
	return t.TopTerms.Parse(params)
}

func (t *TopTermCount) Create(uri string, coord *binning.TileCoord, query prism.Query) ([]byte, error) {
	// get client
	client, err := NewClient(t.Host, t.Port)
	if err != nil {
		return nil, err
	}

	// create root query
	citusQuery, err := t.CreateQuery(query)
	if err != nil {
		return nil, err
	}
	citusQuery.From(uri)

	// add tiling query
	citusQuery = t.Bivariate.AddQuery(coord, citusQuery)

	// get agg
	citusQuery = t.TopTerms.AddAggs(citusQuery)

	// send query
	res, err := client.Query(citusQuery.GetQuery(false), citusQuery.QueryArgs...)
	if err != nil {
		return nil, err
	}

	// marshal results
	counts, err := t.TopTerms.GetTerms(res)
	return json.Marshal(counts)
}
