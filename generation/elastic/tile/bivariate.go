package tile

import (
	"fmt"
	"math"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/param"
	"github.com/unchartedsoftware/prism/util/json"
)

// Bivariate represents a bivariate tile generator.
type Bivariate struct {
	generation.Bivariate
}

// ApplyQuery returns the x query.
func (b *Bivariate) ApplyQuery(query *elastic.NewBoolQuery) error {
	query.Must(elastic.NewRangeQuery(b.XField).
		Gte(b.MinX).
		Lt(b.MaxX))
	query.Must(elastic.NewRangeQuery(b.YField).
		Gte(b.MinY).
		Lt(b.MaxY))
	return nil
}

// ApplyAgg applies an elastic agg.
func (b *Bivariate) ApplyAgg(searchService *elastic.SearchService) error {
	intervalX := int64(math.Max(1, b.XRange/b.Resolution))
	intervalY := int64(math.Max(1, b.YRange/b.Resolution))
	x := elastic.NewHistogramAggregation().
		Field(b.XField).
		Offset(b.MinX).
		Interval(intervalX).
		MinDocCount(1)
	y := elastic.NewHistogramAggregation().
		Field(b.YField).
		Offset(b.MinY).
		Interval(intervalY).
		MinDocCount(1)
	x.SubAggregation("y", yAgg)
	search.Aggregation("x", x)
	return nil
}

// GetBins parses the resulting histograms into bins.
func (b *Bivariate) GetBins(res *elastic.SearchResult) ([]*elastic.AggregationBucketHistogramItem, error) {
	// parse aggregations
	xAgg, ok := res.Aggregations.Histogram("x")
	if !ok {
		return nil, fmt.Errorf("Histogram aggregation `x` was not found")
	}
	// allocate bins
	bins := make([]*elastic.AggregationBucketHistogramItem, b.Resolution*b.Resolution)
	// fill bins
	for xBin, xBucket := range xAgg.Buckets {
		yAgg, ok := xBucket.Histogram("y")
		if !ok {
			return nil, fmt.Errorf("Histogram aggregation `y` was not found")
		}
		for yBin, yBucket := range yAgg.Buckets {
			index := xBin + b.Resolution*yBin
			bins[index] = yBucket
		}
	}
	return bins, nil
}
