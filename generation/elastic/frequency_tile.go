package elastic

import (
	"encoding/json"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
)

// FrequencyTile represents an elasticsearch implementation of the frequency
// tile.
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
	// add frequency query
	q.Must(t.Frequency.GetQuery())
	// set the query
	search.Query(q)

	// get agg
	aggs := t.Frequency.GetAggs()
	// set the aggregation
	search.Aggregation("frequency", aggs["frequency"])

	// send query
	res, err := search.Do()
	if err != nil {
		return nil, err
	}

	// get buckets
	frequency, err := t.Frequency.GetBuckets(&res.Aggregations)
	if err != nil {
		return nil, err
	}

	buckets := make([]map[string]interface{}, len(frequency))
	for i, bucket := range frequency {
		buckets[i] = map[string]interface{}{
			"timestamp": bucket.Key,
			"count":     bucket.DocCount,
		}
	}
	// marshal results
	return json.Marshal(buckets)
}
