package elastic

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/tile"
)

type Frequency struct {
	tile.Frequency
}

func (f *Frequency) GetQuery() *elastic.RangeQuery {
	query := elastic.NewRangeQuery(f.FrequencyField)
	if f.GTE != nil {
		query.Gte(f.CastTime(f.GTE))
	}
	if f.GT != nil {
		query.Gt(f.CastTime(f.GT))
	}
	if f.LTE != nil {
		query.Lte(f.CastTime(f.LTE))
	}
	if f.LT != nil {
		query.Lt(f.CastTime(f.LT))
	}
	return query
}

func (f *Frequency) GetAggs() map[string]elastic.Aggregation {
	agg := elastic.NewDateHistogramAggregation().
		Field(f.FrequencyField).
		Interval(f.Interval).
		MinDocCount(0)
	if f.GTE != nil {
		agg.ExtendedBoundsMin(f.CastTime(f.GTE))
		agg.Offset(f.CastTimeToString(f.GTE))
	}
	if f.GT != nil {
		agg.ExtendedBoundsMin(f.CastTime(f.GT))
		agg.Offset(f.CastTimeToString(f.GT))
	}
	if f.LTE != nil {
		agg.ExtendedBoundsMax(f.CastTime(f.LTE))
	}
	if f.LT != nil {
		agg.ExtendedBoundsMax(f.CastTime(f.LT))
	}
	return map[string]elastic.Aggregation{
		"frequency": agg,
	}
}

func (f *Frequency) GetBuckets(aggs *elastic.Aggregations) ([]*elastic.AggregationBucketHistogramItem, error) {
	frequency, ok := aggs.DateHistogram("frequency")
	if !ok {
		return nil, fmt.Errorf("date histogram aggregation `frequency` was not found")
	}
	return frequency.Buckets, nil
}
