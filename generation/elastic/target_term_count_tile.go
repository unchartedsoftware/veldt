package elastic

import (
	"encoding/json"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
)

// TargetTermCountTile represents an elasticsearch implementation of the
// target term count tile.
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
	aggs := t.TargetTerms.GetAggs()
	for term, agg := range aggs {
		// set the aggregation
		search.Aggregation(term, agg)
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
	counts := make(map[string]uint32)
	for term, bucket := range terms {
		counts[term] = uint32(bucket.DocCount)
	}
	// marshal results
	return json.Marshal(counts)
}
