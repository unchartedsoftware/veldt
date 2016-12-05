package citus

import (
	"encoding/json"

	"github.com/unchartedsoftware/prism"
	"github.com/unchartedsoftware/prism/binning"
)

type FrequencyTile struct {
	Bivariate
	Frequency
	Tile
}

func NewFrequencyTile(host, port string) prism.TileCtor {
	return func() (prism.Tile, error) {
		t := &FrequencyTile{}
		t.Host = host
		t.Port = port
		return t, nil
	}
}

func (t *FrequencyTile) Parse(params map[string]interface{}) error {
	err := t.Bivariate.Parse(params)
	if err != nil {
		return nil
	}
	return t.Frequency.Parse(params)
}

func (t *FrequencyTile) Create(uri string, coord *binning.TileCoord, query prism.Query) ([]byte, error) {
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

	buckets := make([]map[string]interface{}, len(frequency))
	for i, bucket := range frequency {
		buckets[i] = map[string]interface{}{
			"timestamp": bucket.Bucket,
			"count":     bucket.Value,
		}
	}
	// marshal results
	return json.Marshal(buckets)
}
