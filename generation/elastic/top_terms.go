package elastic

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/veldt/tile"
)

// TopTerms represents an elasticsearch implementation of the top terms tile.
type TopTerms struct {
	tile.TopTerms
}

// GetAggs returns the appropriate elasticsearch aggregation for the tile.
func (t *TopTerms) GetAggs() map[string]*elastic.TermsAggregation {
	agg := elastic.NewTermsAggregation().
		Field(t.TermsField).
		Size(int(t.TermsCount))
	return map[string]*elastic.TermsAggregation{
		"top-terms": agg,
	}
}

// GetTerms returns the individual term buckets from the provided aggregation.
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
