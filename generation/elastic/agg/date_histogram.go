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
		query.Gte(castTime(p.From))
	}
	if p.To != nil {
		query.Lte(castTime(p.To))
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
		agg.ExtendedBoundsMin(castTime(p.From))
		agg.Offset(castTimeToString(p.From))
	}
	if p.To != nil {
		agg.ExtendedBoundsMax(castTime(p.To))
	}
	return agg
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
