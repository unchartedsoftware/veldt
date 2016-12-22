package citus

import (
	"encoding/json"

	"github.com/unchartedsoftware/prism"
	"github.com/unchartedsoftware/prism/binning"
)

type TopTermCountTile struct {
	Bivariate
	TopTerms
	Tile
}

func NewTopTermCountTile(host, port string) prism.TileCtor {
	return func() (prism.Tile, error) {
		t := &TopTermCountTile{}
		t.Host = host
		t.Port = port
		return t, nil
	}
}

func (t *TopTermCountTile) Parse(params map[string]interface{}) error {
	err := t.Bivariate.Parse(params)
	if err != nil {
		return err
	}
	return t.TopTerms.Parse(params)
}

func (t *TopTermCountTile) Create(uri string, coord *binning.TileCoord, query prism.Query) ([]byte, error) {
	// Initialize the tile processing.
	client, citusQuery, err := t.InitliazeTile(uri, query)

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
