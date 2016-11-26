package elastic

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/tile"
)

// TopTerms represents a tiling generator that produces heatmaps.
type TopTerms struct {
	tile.TopTerms
}

func (t *TopTerms) GetAggs() map[string]elastic.Aggregation {
	agg := elastic.NewTermsAggregation().
		Field(t.TermsField).
		Size(int(t.TermsCount))
	return map[string]elastic.Aggregation{
		"top-terms": agg,
	}
}

// GetBins parses the resulting histograms into bins.
func (t *TopTerms) GetTerms(aggs *elastic.Aggregations) (map[string]*elastic.AggregationBucketKeyItem, error) {
	// build map of topics and counts
	counts := make(map[string]*elastic.AggregationBucketKeyItem)
	terms, ok := aggs.Terms("top-terms")
	if !ok {
		return nil, fmt.Errorf("terms aggregation `top-term` was not found")
	}
	for _, bucket := range terms.Buckets {
		term, ok := bucket.Key.(string)
		if !ok {
			return nil, fmt.Errorf("terms aggregation key was not of type `string`")
		}
		counts[term] = bucket
	}
	return counts, nil
}
