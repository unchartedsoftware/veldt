package citus

import (
	"fmt"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
	"github.com/unchartedsoftware/veldt/util/json"
)

// TermsFrequencyTile represents a citus implementation of the terms frequency tile.
type TermsFrequencyTile struct {
	Bivariate
	TermsFrequency
	Frequency
	Tile
}

// NewTermsFrequencyTile instantiates and returns a new tile struct.
func NewTermsFrequencyTile(cfg *Config) veldt.TileCtor {
	return func() (veldt.Tile, error) {
		t := &TermsFrequencyTile{}
		t.Config = cfg
		return t, nil
	}
}

// Parse parses the provided JSON object and populates the tiles attributes.
func (t *TermsFrequencyTile) Parse(params map[string]interface{}) error {
	err := t.Bivariate.Parse(params)
	if err != nil {
		return err
	}
	return t.TermsFrequency.Parse(params)
}

// Create generates a tile from the provided URI, tile coordinate and query parameters.
func (t *TermsFrequencyTile) Create(uri string, coord *binning.TileCoord, query veldt.Query) ([]byte, error) {
	// Initialize the tile processing.
	client, citusQuery, err := t.InitializeTile(uri, query)
	if err != nil {
		return nil, err
	}

	// add tiling query
	citusQuery = t.Bivariate.AddQuery(coord, citusQuery)

	// get aggs
	citusQuery.Select(t.Frequency.FrequencyField)
	citusQuery = t.TermsFrequency.AddAggs(citusQuery)
	citusQuery = t.Frequency.AddAggs(citusQuery)

	// send query
	res, err := client.Query(citusQuery.GetQuery(false), citusQuery.QueryArgs...)
	if err != nil {
		return nil, err
	}

	// parse results. Every row should have the frequency buckets + the term.
	// Probably best to add a sort on the query to group the terms together.
	// Can also determine the buckets for the frequency once and then just read the values.
	// Results are stored in a map -> frequency bucket.
	rawResults := make(map[string]map[int64]float64)
	for res.Next() {
		var term string
		var count uint32
		var bucket int64
		var frequency int
		err := res.Scan(&term, &count, &bucket, &frequency)
		if err != nil {
			return nil, fmt.Errorf("Error parsing top terms: %v", err)
		}
		// TODO: May need to do some checking to see if things exist already.
		rawResults[term][bucket] = float64(frequency)
	}

	// Build frequency buckets & encode.
	result := make(map[string][]map[string]interface{})
	for term, frequency := range rawResults {
		// get buckets
		buckets, err := t.Frequency.CreateBuckets(frequency)
		if err != nil {
			return nil, err
		}
		// add frequency
		frequency := t.Frequency.encodeResult(buckets)
		result[term] = frequency
	}
	// marshal results
	return json.Marshal(result)
}
