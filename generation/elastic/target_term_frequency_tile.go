package elastic

import (
	"encoding/json"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
)

// TargetTermFrequencyTile represents an elasticsearch implementation of the
// target term frequency tile.
type TargetTermFrequencyTile struct {
	Bivariate
	TargetTerms
	Frequency
	Tile
}

// NewTargetTermFrequencyTile instantiates and returns a new tile struct.
func NewTargetTermFrequencyTile(host, port string) veldt.TileCtor {
	return func() (veldt.Tile, error) {
		t := &TargetTermFrequencyTile{}
		t.Host = host
		t.Port = port
		return t, nil
	}
}

// Parse parses the provided JSON object and populates the tiles attributes.
func (t *TargetTermFrequencyTile) Parse(params map[string]interface{}) error {
	err := t.Bivariate.Parse(params)
	if err != nil {
		return err
	}
	return t.TargetTerms.Parse(params)
}

// Create generates a tile from the provided URI, tile coordinate and query
// parameters.
func (t *TargetTermFrequencyTile) Create(uri string, coord *binning.TileCoord, query veldt.Query) ([]byte, error) {
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
