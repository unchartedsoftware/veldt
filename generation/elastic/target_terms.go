package elastic

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/veldt/tile"
)

// TargetTerms represents an elasticsearch implementation of the target terms
// tile.
type TargetTerms struct {
	tile.TargetTerms
}

// GetQuery returns the appropriate elasticsearch query for the tile.
func (t *TargetTerms) GetQuery() elastic.Query {
	terms := make([]interface{}, len(t.Terms))
	for i, term := range t.Terms {
		terms[i] = term
	}
	return elastic.NewTermsQuery(t.TermsField, terms...)
}

// GetAggs returns the appropriate elasticsearch aggregation for the tile.
func (t *TargetTerms) GetAggs() map[string]*elastic.FilterAggregation {
	aggs := make(map[string]*elastic.FilterAggregation, len(t.Terms))
	// add all filter aggregations
	for _, term := range t.Terms {
		aggs[term] = elastic.NewFilterAggregation().
			Filter(elastic.NewTermQuery(t.TermsField, term))
	}
	return aggs
}

// GetTerms returns the individual term buckets from the provided aggregation.
func (t *TargetTerms) GetTerms(aggs *elastic.Aggregations) (map[string]*elastic.AggregationSingleBucket, error) {
	res := make(map[string]*elastic.AggregationSingleBucket)
	for _, term := range t.Terms {
		filter, ok := aggs.Filter(term)
		if !ok {
			return nil, fmt.Errorf("filter aggregation '%s' was not found", term)
		}
		res[term] = filter
	}
	return res, nil
}
