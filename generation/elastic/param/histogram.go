package param

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/util/json"
)

// Histogram represents params for extracting histogram buckets.
type Histogram struct {
	Field    string
	Interval int64
}

// NewHistogram instantiates and returns a new sentiment parameter object.
func NewHistogram(tileReq *tile.Request) (*Histogram, error) {
	params, ok := json.GetChild(tileReq.Params, "histogram")
	if !ok {
		return nil, ErrMissing
	}
	field, ok := json.GetString(params, "field")
	if !ok {
		return nil, fmt.Errorf("Histogram `field` parameter missing from tiling request %s", tileReq.String())
	}
	interval, ok := json.GetNumber(params, "interval")
	if !ok {
		return nil, fmt.Errorf("Histogram `interval` parameter missing from tiling request %s", tileReq.String())
	}
	return &Histogram{
		Field:    field,
		Interval: int64(interval),
	}, nil
}

// GetHash returns a string hash of the parameter state.
func (p *Histogram) GetHash() string {
	return fmt.Sprintf("%s:%d",
		p.Field,
		p.Interval)
}

// GetAggregation returns an elastic query.
func (p *Histogram) GetAggregation() *elastic.HistogramAggregation {
	return elastic.NewHistogramAggregation().
		Field(p.Field).
		Interval(p.Interval).
		MinDocCount(0)
}

// GetBucketMap parses the histogram buckets into a map of keys to counts.
func (p *Histogram) GetBucketMap(agg *elastic.AggregationBucketHistogramItems) map[string]uint64 {
	m := make(map[string]uint64)
	for _, bucket := range agg.Buckets {
		m[fmt.Sprintf("%d", bucket.Key)] = uint64(bucket.DocCount)
	}
	return m
}
