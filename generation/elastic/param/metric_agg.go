package param

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/util/json"
)

// MetricAgg represents params for binning the data within the tile.
type MetricAgg struct {
	Field string
	Type  string
}

// NewMetricAgg instantiates and returns a new metric aggregation parameter.
func NewMetricAgg(tileReq *tile.Request) (*MetricAgg, error) {
	params, ok := json.GetChild(tileReq.Params, "metric_agg")
	if !ok {
		return nil, ErrMissing
	}
	field, ok := json.GetString(params, "field")
	if !ok {
		return nil, fmt.Errorf("MetricAgg `field` parameter missing from tiling request %s", tileReq.String())
	}
	typ, ok := json.GetString(params, "type")
	if !ok {
		return nil, fmt.Errorf("MetricAgg `type` parameter missing from tiling request %s", tileReq.String())
	}
	return &MetricAgg{
		Field: field,
		Type:  typ,
	}, nil
}

// GetHash returns a string hash of the parameter state.
func (p *MetricAgg) GetHash() string {
	return fmt.Sprintf("%s:%s",
		p.Type,
		p.Field)
}

// GetAgg returns an elastic aggregation.
func (p *MetricAgg) GetAgg() elastic.Aggregation {
	switch p.Type {
	case "min":
		return elastic.NewMinAggregation().
			Field(p.Field)
	case "max":
		return elastic.NewMaxAggregation().
			Field(p.Field)
	case "avg":
		return elastic.NewAvgAggregation().
			Field(p.Field)
	default:
		return elastic.NewSumAggregation().
			Field(p.Field)
	}
}

// GetAggValue extracts the value metric based on the type of operation
// specified.
func (p *MetricAgg) GetAggValue(aggName string, aggs *elastic.AggregationBucketHistogramItem) (float64, bool) {
	var metric *elastic.AggregationValueMetric
	var ok bool
	switch p.Type {
	case "min":
		metric, ok = aggs.Min(aggName)
	case "max":
		metric, ok = aggs.Max(aggName)
	case "avg":
		metric, ok = aggs.Avg(aggName)
	default:
		metric, ok = aggs.Sum(aggName)
	}
	if !ok {
		return 0, false
	}
	if metric.Value != nil {
		return *metric.Value, true
	}
	return 0, true
}
