package elastic

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/veldt/tile"
)

// Frequency represents an elasticsearch implementation of the frequency
// tile.
type Frequency struct {
	tile.Frequency
}

// GetQuery returns the appropriate elasticsearch query for the tile.
func (f *Frequency) GetQuery() *elastic.RangeQuery {
	query := elastic.NewRangeQuery(f.FrequencyField)
	if f.GTE != nil {
		query.Gte(castTime(f.GTE))
	}
	if f.GT != nil {
		query.Gt(castTime(f.GT))
	}
	if f.LTE != nil {
		query.Lte(castTime(f.LTE))
	}
	if f.LT != nil {
		query.Lt(castTime(f.LT))
	}
	return query
}

// GetAggs returns the appropriate elasticsearch aggregation for the tile.
func (f *Frequency) GetAggs() map[string]elastic.Aggregation {
	agg := elastic.NewDateHistogramAggregation().
		Field(f.FrequencyField).
		Interval(f.Interval).
		MinDocCount(0)
	if f.GTE != nil {
		agg.ExtendedBoundsMin(castTime(f.GTE))
		agg.Offset(castTimeToString(f.GTE))
	}
	if f.GT != nil {
		agg.ExtendedBoundsMin(castTime(f.GT))
		agg.Offset(castTimeToString(f.GT))
	}
	if f.LTE != nil {
		agg.ExtendedBoundsMax(castTime(f.LTE))
	}
	if f.LT != nil {
		agg.ExtendedBoundsMax(castTime(f.LT))
	}
	return map[string]elastic.Aggregation{
		"frequency": agg,
	}
}

// GetBuckets returns the individual frequency buckets from an elasticsearch
// aggregation.
func (f *Frequency) GetBuckets(aggs *elastic.Aggregations) ([]*elastic.AggregationBucketHistogramItem, error) {
	frequency, ok := aggs.DateHistogram("frequency")
	if !ok {
		return nil, fmt.Errorf("date histogram aggregation `frequency` was not found")
	}
	return frequency.Buckets, nil
}

func castTimeToString(val interface{}) string {
	num, isNum := val.(float64)
	if isNum {
		// assume milliseconds
		return fmt.Sprintf("%dms\n", int64(num))
	}
	str, isStr := val.(string)
	if isStr {
		return str
	}
	return ""
}

func castTime(val interface{}) interface{} {
	num, isNum := val.(float64)
	if isNum {
		return int64(num)
	}
	str, isStr := val.(string)
	if isStr {
		return str
	}
	return val
}
