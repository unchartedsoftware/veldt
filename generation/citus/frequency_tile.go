package citus

import (
	"encoding/json"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
)

// FrequencyTile represents a citus implementation of the frequency tile.
type FrequencyTile struct {
	Bivariate
	Frequency
	Tile
}

// NewFrequencyTile instantiates and returns a new tile struct.
func NewFrequencyTile(host, port string) veldt.TileCtor {
	return func() (veldt.Tile, error) {
		t := &FrequencyTile{}
		t.Host = host
		t.Port = port
		return t, nil
	}
}

// Parse parses the provided JSON object and populates the tiles attributes.
func (t *FrequencyTile) Parse(params map[string]interface{}) error {
	err := t.Bivariate.Parse(params)
	if err != nil {
		return err
	}
	return t.Frequency.Parse(params)
}

// Create generates a tile from the provided URI, tile coordinate and query
// parameters.
func (t *FrequencyTile) Create(uri string, coord *binning.TileCoord, query veldt.Query) ([]byte, error) {
	// Initialize the tile processing.
	client, citusQuery, err := t.InitializeTile(uri, query)

	// add tiling query
	citusQuery = t.Bivariate.AddQuery(coord, citusQuery)

	// add frequency query
	citusQuery = t.Frequency.AddQuery(citusQuery)

	// add aggs
	citusQuery = t.Frequency.AddAggs(citusQuery)

	// send query
	res, err := client.Query(citusQuery.GetQuery(false), citusQuery.QueryArgs...)
	if err != nil {
		return nil, err
	}

	// get buckets
	frequency, err := t.Frequency.GetBuckets(res)
	if err != nil {
		return nil, err
	}

	buckets := t.Frequency.encodeResult(frequency)

	// marshal results
	return json.Marshal(buckets)
}
