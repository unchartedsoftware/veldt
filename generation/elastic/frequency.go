package elastic

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/tile"
)

// Frequency represents a tiling generator that produces heatmaps.
type Frequency struct {
	tile.Frequency
}

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

func (f *Frequency) GetBuckets(aggs *elastic.Aggregations) ([]*elastic.AggregationBucketHistogramItem, error) {
	frequency, ok := aggs.DateHistogram("frequency")
	if !ok {
		return nil, fmt.Errorf("date histogram aggregation `frequency` was not found")
	}
	return frequency.Buckets, nil
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
