package citus

import (
	"encoding/json"
	"fmt"

	"github.com/unchartedsoftware/prism"
	"github.com/unchartedsoftware/prism/binning"
)

type TargetTermFrequencyTile struct {
	Bivariate
	TargetTerms
	Frequency
	Tile
}

func NewTargetTermFrequencyTile(host, port string) prism.TileCtor {
	return func() (prism.Tile, error) {
		t := &TargetTermFrequencyTile{}
		t.Host = host
		t.Port = port
		return t, nil
	}
}

func (t *TargetTermFrequencyTile) Parse(params map[string]interface{}) error {
	err := t.Bivariate.Parse(params)
	if err != nil {
		return err
	}
	return t.TargetTerms.Parse(params)
}

func (t *TargetTermFrequencyTile) Create(uri string, coord *binning.TileCoord, query prism.Query) ([]byte, error) {
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

	// get aggs
	//Need to add the frequency field to the select list before calling the target terms.
	//Otherwise it will be impossible to build the aggregate since target terms nests the query.
	citusQuery.Select(t.Frequency.FrequencyField)
	citusQuery = t.TargetTerms.AddAggs(citusQuery)
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
		var term_count uint32
		var bucket int64
		var frequency int
		err := res.Scan(&term, &term_count, &bucket, &frequency)
		if err != nil {
			return nil, fmt.Errorf("Error parsing top terms: %v", err)
		}
		//TODO: May need to do some checking to see if things exist already.
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
		frequency := make([]map[string]interface{}, len(buckets))
		for i, bucket := range buckets {
			frequency[i] = map[string]interface{}{
				"timestamp": bucket.Bucket,
				"count":     bucket.Value,
			}
		}
		result[term] = frequency
	}
	// marshal results
	return json.Marshal(result)
}
