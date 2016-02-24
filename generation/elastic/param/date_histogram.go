package param

import (
	"fmt"
	"time"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/util/json"
)

const (
	defaultField    = "timestamp"
	defaultInterval = "1d"
	defaultFrom     = 0
)

// DateHistogram represents params for restricting the time range of a tile.
type DateHistogram struct {
	Field    string
	From     uint64
	To       uint64
	Interval string
}

// NewDateHistogram instantiates and returns a new time bucketing parameter object.
func NewDateHistogram(tileReq *tile.Request) (*DateHistogram, error) {
	params := json.GetChildOrEmpty(tileReq.Params, "date_histogram")
	field := json.GetStringDefault(params, "field", defaultField)
	from := uint64(json.GetNumberDefault(params, "from", defaultFrom))
	to := uint64(json.GetNumberDefault(params, "to", float64(time.Now().Unix())))
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
	return elastic.NewRangeQuery(p.Field).
		Gte(p.From).
		Lte(p.To)
}

// GetAggregation returns an elastic query.
func (p *DateHistogram) GetAggregation() *elastic.DateHistogramAggregation {
	return elastic.NewDateHistogramAggregation().
		Field(p.Field).
		MinDocCount(0).
		Interval(p.Interval).
		ExtendedBoundsMin(p.From).
		ExtendedBoundsMax(p.To).
		Offset(fmt.Sprintf("%ds", p.From/1000))
}
