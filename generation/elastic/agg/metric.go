package agg

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/elastic/param"
	"github.com/unchartedsoftware/prism/util/json"
)

// Metric represents params for binning the data within the tile.
type Metric struct {
	Field string
	Type  string
}

// NewMetric instantiates and returns a new metric aggregation parameter.
func NewMetric(params map[string]interface{}) (*Metric, error) {
	params, ok := json.GetChild(params, "metric")
	if !ok {
		return nil, param.ErrMissing
	}
	field, ok := json.GetString(params, "field")
	if !ok {
		return nil, fmt.Errorf("Metric `field` parameter missing from tiling param %v", params)
	}
	typ, ok := json.GetString(params, "type")
	if !ok {
		return nil, fmt.Errorf("Metric `type` parameter missing from tiling param %v", params)
	}
	return &Metric{
		Field: field,
		Type:  typ,
	}, nil
}

// GetHash returns a string hash of the parameter state.
func (p *Metric) GetHash() string {
	return fmt.Sprintf("%s:%s",
		p.Type,
		p.Field)
}

// GetAgg returns an elastic aggregation.
func (p *Metric) GetAgg() elastic.Aggregation {
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
func (p *Metric) GetAggValue(aggName string, aggs *elastic.AggregationBucketHistogramItem) (float64, bool) {
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
