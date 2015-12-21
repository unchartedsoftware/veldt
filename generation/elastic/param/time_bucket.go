package param

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/util/json"
)

// TimeBucket represents params for restricting the time range of a tile.
type TimeBucket struct {
	TimeRange *TimeRange
	Interval  string
}

// NewTimeBucket instantiates and returns a new time bucketing parameter object.
func NewTimeBucket(tileReq *tile.Request) (*TimeBucket, error) {
	timeRange, err := NewTimeRange(tileReq)
	if err != nil {
		return nil, err
	}
	params := tileReq.Params
	return &TimeBucket{
		TimeRange: timeRange,
		Interval:  json.GetStringDefault(params, "interval", "1d"),
	}, nil
}

// GetHash returns a string hash of the parameter state.
func (p *TimeBucket) GetHash() string {
	return fmt.Sprintf("%s:%s",
		p.TimeRange.GetHash(),
		p.Interval)
}

// GetTimeAggregation returns an elastic query.
func (p *TimeBucket) GetTimeAggregation() *elastic.DateHistogramAggregation {
	return elastic.NewDateHistogramAggregation().
		Field(p.TimeRange.Time).
		MinDocCount(0).
		Interval(p.Interval).
		ExtendedBoundsMin(p.TimeRange.From).
		ExtendedBoundsMax(p.TimeRange.To).
		PreOffset(-p.TimeRange.From).
		PreOffset(p.TimeRange.From)
}
