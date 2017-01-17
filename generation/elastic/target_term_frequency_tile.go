package elastic

import (
	"encoding/json"

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
	// create search service
	search := client.Search().
		Index(uri).
		Size(0)

	// create root query
	q, err := t.CreateQuery(query)
	if err != nil {
		return nil, err
	}
	// add tiling query
	q.Must(t.Bivariate.GetQuery(coord))
	// set the query
	search.Query(q)
	// get aggs
	termAggs := t.TargetTerms.GetAggs()
	freqAggs := t.Frequency.GetAggs()
	for term, agg := range termAggs {
		// set the aggregation
		search.Aggregation(term, agg.SubAggregation("frequency", freqAggs["frequency"]))
	}
	// send query
	res, err := search.Do()
	if err != nil {
		return nil, err
	}
	// get terms
	terms, err := t.TargetTerms.GetTerms(&res.Aggregations)
	if err != nil {
		return nil, err
	}
	// encode
	result := make(map[string][]map[string]interface{})
	for term, item := range terms {
		// get buckets
		buckets, err := t.Frequency.GetBuckets(&item.Aggregations)
		if err != nil {
			return nil, err
		}
		// add frequency
		frequency := make([]map[string]interface{}, len(buckets))
		for i, bucket := range buckets {
			frequency[i] = map[string]interface{}{
				"timestamp": bucket.Key,
				"count":     bucket.DocCount,
			}
		}
		result[term] = frequency
	}
	// marshal results
	return json.Marshal(result)
}
