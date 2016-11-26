package elastic

import (
	"encoding/json"

	"github.com/unchartedsoftware/prism"
	"github.com/unchartedsoftware/prism/binning"
)

type TopTermCount struct {
	Bivariate
	TopTerms
	Tile
}

func NewTopTermCountTile(host, port string) prism.TileCtor {
	return func() (prism.Tile, error) {
		t := &TopTermCount{}
		t.Host = host
		t.Port = port
		return t, nil
	}
}

func (t *TopTermCount) Parse(params map[string]interface{}) error {
	err := t.Bivariate.Parse(params)
	if err != nil {
		return nil
	}
	return t.TopTerms.Parse(params)
}

func (t *TopTermCount) Create(uri string, coord *binning.TileCoord, query prism.Query) ([]byte, error) {
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
	aggs := t.TopTerms.GetAggs()
	// set the aggregation
	search.Aggregation("top-terms", aggs["top-terms"])

	// send query
	res, err := search.Do()
	if err != nil {
		return nil, err
	}

	// get bins
	terms, err := t.TopTerms.GetTerms(&res.Aggregations)
	if err != nil {
		return nil, err
	}

	counts := make(map[string]uint32)
	for term, bucket := range terms {
		counts[term] = uint32(bucket.DocCount)
	}
	// marshal results
	return json.Marshal(counts)
}
