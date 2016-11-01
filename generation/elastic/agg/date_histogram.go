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
	From     interface{}
	To       interface{}
	Interval string
}

// NewDateHistogram instantiates and returns a new time bucketing parameter object.
func NewDateHistogram(params map[string]interface{}) (*DateHistogram, error) {
	params = json.GetChildOrEmpty(params, "date_histogram")
	field := json.GetStringDefault(params, defaultField, "field")
	from, _ := json.Get(params, "from")
	to, _ := json.Get(params, "to")
	interval := json.GetStringDefault(params, defaultInterval, "interval")
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
	if p.From != nil {
		query.Gte(p.From)
	}
	if p.To != nil {
		query.Lte(p.To)
	}
	return query
}

// GetAgg returns an elastic query.
func (p *DateHistogram) GetAgg() *elastic.DateHistogramAggregation {
	agg := elastic.NewDateHistogramAggregation().
		Field(p.Field).
		Interval(p.Interval).
		MinDocCount(0)
	if p.From != nil {
		agg.ExtendedBoundsMin(p.From)
		num, isNum := p.From.(float64)
		if isNum {
			// assume milliseconds
			agg.Offset(fmt.Sprintf("%ds", int64(num)/1000))
		} else {
			str, isStr := p.From.(string)
			if isStr {
				agg.Offset(str)
			}
		}
	}
	if p.To != nil {
		agg.ExtendedBoundsMax(p.To)
	}
	return agg
}
