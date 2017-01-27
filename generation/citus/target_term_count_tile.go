package citus

import (
	"encoding/json"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
)

// TargetTermCountTile represents a citus implementation of the target term
// count tile.
type TargetTermCountTile struct {
	Bivariate
	TargetTerms
	Tile
}

// NewTargetTermCountTile instantiates and returns a new tile struct.
func NewTargetTermCountTile(host, port string) veldt.TileCtor {
	return func() (veldt.Tile, error) {
		t := &TargetTermCountTile{}
		t.Host = host
		t.Port = port
		return t, nil
	}
}

// Parse parses the provided JSON object and populates the tiles attributes.
func (t *TargetTermCountTile) Parse(params map[string]interface{}) error {
	err := t.Bivariate.Parse(params)
	if err != nil {
		return err
	}
	return t.TargetTerms.Parse(params)
}

// Create generates a tile from the provided URI, tile coordinate and query
// parameters.
func (t *TargetTermCountTile) Create(uri string, coord *binning.TileCoord, query veldt.Query) ([]byte, error) {
	// Initialize the tile processing.
	client, citusQuery, err := t.InitializeTile(uri, query)

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
