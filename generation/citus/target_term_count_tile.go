package citus

import (
	"encoding/json"

	"github.com/unchartedsoftware/prism"
	"github.com/unchartedsoftware/prism/binning"
)

type TargetTermCountTile struct {
	Bivariate
	TargetTerms
	Tile
}

func NewTargetTermCountTile(host, port string) prism.TileCtor {
	return func() (prism.Tile, error) {
		t := &TargetTermCountTile{}
		t.Host = host
		t.Port = port
		return t, nil
	}
}

func (t *TargetTermCountTile) Parse(params map[string]interface{}) error {
	err := t.Bivariate.Parse(params)
	if err != nil {
		return err
	}
	return t.TargetTerms.Parse(params)
}

func (t *TargetTermCountTile) Create(uri string, coord *binning.TileCoord, query prism.Query) ([]byte, error) {
	// Initialize the tile processing.
	client, citusQuery, err := t.InitliazeTile(uri, query)

	// add tiling query
	citusQuery = t.Bivariate.AddQuery(coord, citusQuery)

	// get aggs
	citusQuery = t.TargetTerms.AddAggs(citusQuery)

	// send query
	res, err := client.Query(citusQuery.GetQuery(false), citusQuery.QueryArgs...)
	if err != nil {
		return nil, err
	}

	// get terms
	terms, err := t.TargetTerms.GetTerms(res)
	if err != nil {
		return nil, err
	}

	// marshal results
	return json.Marshal(terms)
}
