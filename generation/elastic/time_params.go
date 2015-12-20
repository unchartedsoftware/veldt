package elastic

import (
	"fmt"
	"time"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/util/json"
)

// TimeParams represents params for restricting the time range of a tile.
type TimeParams struct {
	Time     string
	From   	 int64
	To       int64
	Interval string
}

// NewTimeParams parses the params map returns a pointer to the param struct.
func NewTimeParams(tileReq *tile.Request) *TimeParams {
	params := tileReq.Params
	if json.Exists(params, "from") || json.Exists(params, "to") {
		return &TimeParams{
			Time: json.GetStringDefault(params, "time", "timestamp"),
			From: int64(json.GetNumberDefault(params, "from", 0)),
			To: int64(json.GetNumberDefault(params, "to", float64(time.Now().Unix()))),
			Interval: json.GetStringDefault(params, "interval", "1d"),
		}
	}
	return nil
}

// GetHash returns a string hash of the parameter state.
func (p *TimeParams) GetHash() string {
	return fmt.Sprintf("%s:%d:%d:%s",
		p.Time,
		p.From,
		p.To,
		p.Interval)
}

// GetTimeQuery returns an elastic query.
func (p *TimeParams) GetTimeQuery() *elastic.RangeQuery {
	return elastic.NewRangeQuery(p.Time).
		Gte(p.From).
		Lte(p.To)
}

// GetTimeAggregation returns an elastic query.
func (p *TimeParams) GetTimeAggregation() *elastic.DateHistogramAggregation {
	return elastic.NewDateHistogramAggregation().
		Field(p.Time).
		MinDocCount(0).
		Interval(p.Interval).
		ExtendedBoundsMin(p.From).
		ExtendedBoundsMax(p.To).
		PreOffset(-p.From).
		PreOffset(p.From)
}
