package agg

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/util/json"
)

const (
	defaultField    = "timestamp"
	defaultInterval = "1d"
)

// DateHistogram represents params for restricting the time range of a tile.
type DateHistogram struct {
	Field    string
	From     int64
	To       int64
	Interval string
}

// NewDateHistogram instantiates and returns a new time bucketing parameter object.
func NewDateHistogram(params map[string]interface{}) (*DateHistogram, error) {
	params = json.GetChildOrEmpty(params, "date_histogram")
	field := json.GetStringDefault(params, "field", defaultField)
	from := int64(json.GetNumberDefault(params, "from", -1))
	to := int64(json.GetNumberDefault(params, "to", -1))
	interval := json.GetStringDefault(params, "interval", defaultInterval)
	return &DateHistogram{
		Field:    field,
		From:     from,
		To:       to,
		Interval: interval,
	}, nil
}

// GetHash returns a string hash of the parameter state.
func (p *DateHistogram) GetHash() string {
	return fmt.Sprintf("%s:%d:%d:%s",
		p.Field,
		p.From,
		p.To,
		p.Interval)
}

// GetQuery returns an elastic query.
func (p *DateHistogram) GetQuery() *elastic.RangeQuery {
	query := elastic.NewRangeQuery(p.Field)
	if p.From != -1 {
		query.Gte(p.From)
	}
	if p.To != -1 {
		query.Lte(p.To)
	}
	return query
}

// GetAggregation returns an elastic query.
func (p *DateHistogram) GetAggregation() *elastic.DateHistogramAggregation {
	agg := elastic.NewDateHistogramAggregation().
		Field(p.Field).
		Interval(p.Interval).
		MinDocCount(0)
	if p.From != -1 {
		agg.ExtendedBoundsMin(p.From)
		agg.Offset(fmt.Sprintf("%ds", p.From/1000))
	}
	if p.To != -1 {
		agg.ExtendedBoundsMax(p.To)
	}
	return agg
}
