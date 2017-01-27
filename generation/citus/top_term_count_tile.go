package citus

import (
	"encoding/json"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
)

// TopTermCountTile represents a citus implementation of the top term count
// tile.
type TopTermCountTile struct {
	Bivariate
	TopTerms
	Tile
}

// NewTopTermCountTile instantiates and returns a new tile struct.
func NewTopTermCountTile(host, port string) veldt.TileCtor {
	return func() (veldt.Tile, error) {
		t := &TopTermCountTile{}
		t.Host = host
		t.Port = port
		return t, nil
	}
}

// Parse parses the provided JSON object and populates the tiles attributes.
func (t *TopTermCountTile) Parse(params map[string]interface{}) error {
	err := t.Bivariate.Parse(params)
	if err != nil {
		return err
	}
	return t.TopTerms.Parse(params)
}

// Create generates a tile from the provided URI, tile coordinate and query
// parameters.
func (t *TopTermCountTile) Create(uri string, coord *binning.TileCoord, query veldt.Query) ([]byte, error) {
	// Initialize the tile processing.
	client, citusQuery, err := t.InitializeTile(uri, query)

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
	if err != nil {
		return nil, err
	}
	return json.Marshal(counts)
}
