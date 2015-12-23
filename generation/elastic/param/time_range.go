package param

import (
	"fmt"
	"time"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/util/json"
)

// TimeRange represents params for restricting the time range of a tile.
type TimeRange struct {
	Time string
	From int64
	To   int64
}

// NewTimeRange instantiates and returns a new time range parameter object.
func NewTimeRange(tileReq *tile.Request) (*TimeRange, error) {
	params := tileReq.Params
	if json.Exists(params, "from") || json.Exists(params, "to") {
		return &TimeRange{
			Time: json.GetStringDefault(params, "time", "timestamp"),
			From: int64(json.GetNumberDefault(params, "from", 0)),
			To:   int64(json.GetNumberDefault(params, "to", float64(time.Now().Unix()))),
		}, nil
	}
	return nil, fmt.Errorf("Time parameters missing from tiling request %s", tileReq.String())
}

// GetHash returns a string hash of the parameter state.
func (p *TimeRange) GetHash() string {
	return fmt.Sprintf("%s:%d:%d",
		p.Time,
		p.From,
		p.To)
}

// GetTimeQuery returns an elastic query.
func (p *TimeRange) GetTimeQuery() *elastic.RangeQuery {
	return elastic.NewRangeQuery(p.Time).
		Gte(p.From).
		Lte(p.To)
}
