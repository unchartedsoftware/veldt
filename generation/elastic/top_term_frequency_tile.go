package elastic

import (
	"encoding/json"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
)

// TopTermFrequencyTile represents an elasticsearch implementation of the
// top term frequency tile.
type TopTermFrequencyTile struct {
	Bivariate
	TopTerms
	Frequency
	Tile
}

// NewTopTermFrequencyTile instantiates and returns a new tile struct.
func NewTopTermFrequencyTile(host, port string) veldt.TileCtor {
	return func() (veldt.Tile, error) {
		t := &TopTermFrequencyTile{}
		t.Host = host
		t.Port = port
		return t, nil
	}
}

// Parse parses the provided JSON object and populates the tiles attributes.
func (t *TopTermFrequencyTile) Parse(params map[string]interface{}) error {
	err := t.Bivariate.Parse(params)
	if err != nil {
		return err
	}
	return t.TopTerms.Parse(params)
}

// Create generates a tile from the provided URI, tile coordinate and query
// parameters.
func (t *TopTermFrequencyTile) Create(uri string, coord *binning.TileCoord, query veldt.Query) ([]byte, error) {
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
	// get agg
	topTermsAggs := t.TopTerms.GetAggs()
	freqAggs := t.Frequency.GetAggs()
	agg := topTermsAggs["top-terms"].
		SubAggregation("frequency", freqAggs["frequency"])
	// set the aggregation
	search.Aggregation("top-terms", agg)
	// send query
	res, err := search.Do()
	if err != nil {
		return nil, err
	}
	// get terms
	terms, err := t.TopTerms.GetTerms(&res.Aggregations)
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
