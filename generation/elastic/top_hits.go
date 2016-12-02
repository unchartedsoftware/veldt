package elastic

import (
	"encoding/json"
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/tile"
)

type TopHits struct {
	tile.TopHits
}

func (t *TopHits) GetAggs() map[string]elastic.Aggregation {
	agg := elastic.NewTopHitsAggregation().Size(t.HitsCount)
	// sort
	if t.SortOrder == "desc" {
		agg.Sort(t.SortField, false)
	} else {
		agg.Sort(t.SortField, true)
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

func (t *TopHits) GetTopHits(aggs *elastic.Aggregations) ([]map[string]interface{}, error) {
	topHits, ok := aggs.TopHits("top-hits")
	if !ok {
		return nil, fmt.Errorf("top-hits aggregation `top-hits` was not found")
	}
	hits := make([]map[string]interface{}, len(topHits.Hits.Hits))
	for index, hit := range topHits.Hits.Hits {
		err := json.Unmarshal(*hit.Source, &hits[index])
		if err != nil {
			return nil, err
		}
	}
	return hits, nil
}
