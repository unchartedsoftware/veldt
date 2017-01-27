package elastic

import (
	"encoding/json"
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/veldt/tile"
)

// TopHits represents an elasticsearch implementation of the top hits tile.
type TopHits struct {
	tile.TopHits
}

// GetAggs returns the appropriate elasticsearch aggregation for the tile.
func (t *TopHits) GetAggs() map[string]elastic.Aggregation {
	agg := elastic.NewTopHitsAggregation().Size(t.HitsCount)
	// sort
	if t.SortField != "" {
		if t.SortOrder == "desc" {
			agg.Sort(t.SortField, false)
		} else {
			agg.Sort(t.SortField, true)
		}
	}
	// add includes
	if t.IncludeFields != nil {
		agg.FetchSourceContext(
			elastic.NewFetchSourceContext(true).
				Include(t.IncludeFields...))
	}
	return map[string]elastic.Aggregation{
		"top-hits": agg,
	}
}

// GetTopHits returns the individual hits from the provided aggregation.
func (t *TopHits) GetTopHits(aggs *elastic.Aggregations) ([]map[string]interface{}, error) {
	topHits, ok := aggs.TopHits("top-hits")
	if !ok {
		return nil, fmt.Errorf("top-hits aggregation `top-hits` was not found")
	}
	hits := make([]map[string]interface{}, len(topHits.Hits.Hits))
	for index, hit := range topHits.Hits.Hits {
		var src map[string]interface{}
		err := json.Unmarshal(*hit.Source, &src)
		if err != nil {
			return nil, err
		}
		// flatten the source paths
		flattened := make(map[string]interface{})
		t.flattenSource(flattened, src, "")
		hits[index] = flattened
	}
	return hits, nil
}

func (t *TopHits) flattenSource(res map[string]interface{}, node map[string]interface{}, path string) {
	for key, val := range node {
		subpath := key
		if path != "" {
			subpath = path + "." + key
		}
		sub, ok := val.(map[string]interface{})
		if ok {
			t.flattenSource(res, sub, subpath)
			continue
		}
		res[subpath] = val
	}
}
