package tile

import (
	"fmt"
	"math"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/param"
	"github.com/unchartedsoftware/prism/util/json"
)

// Univariate represents a univariate tile generator.
type Univariate struct {
	generation.Univariate
}

// ApplyQuery returns the x query.
func (u *Univariate) ApplyQuery(query *elastic.NewBoolQuery) error {
	query.Must(elastic.NewRangeQuery(u.Field).
		Gte(u.Min).
		Lt(u.Max))
	return nil
}

// ApplyAgg applies an elastic agg.
func (u *Univariate) ApplyAgg(searchService *elastic.SearchService) error {
	interval := int64(math.Max(1, u.Range/u.Resolution))
	agg := elastic.NewHistogramAggregation().
		Field(u.Field).
		Offset(u.Min).
		Interval(interval).
		MinDocCount(1)
	return search.Aggregation("bins", agg)
}

// GetBins parses the resulting histograms into bins.
func (u *Univariate) GetBins(res *elastic.SearchResult) ([]*elastic.AggregationBucketHistogramItem, error) {
	// parse aggregations
	agg, ok := res.Aggregations.Histogram("bins")
	if !ok {
		return nil, fmt.Errorf("Histogram aggregation `bins` was not found")
	}
	// allocate bins
	bins := make([]*elastic.AggregationBucketHistogramItem, u.Resolution)
	// fill bins buffer
	for i, bucket := range agg.Buckets {
		bins[i] = bucket
	}
	return bins, nil
}
