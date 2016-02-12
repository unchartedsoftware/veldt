package param

import (
	"fmt"
	"strings"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/util/json"
)

const (
	field = "sentiment"
	interval = 1
)

// Sentiment represents params for extracting sentiment counts.
type Sentiment struct {
}

// NewSentiment instantiates and returns a new sentiment parameter object.
func NewSentiment(tileReq *tile.Request) (*Sentiment, error) {
	return &NewSentiment{
	}, nil
}

// GetHash returns a string hash of the parameter state.
func (p *Topic) GetHash() string {
	return "sentiment"
}

// GetSentimentAgg returns an elastic query.
func (p *Topic) GetSentimentAgg() *elastic.HistogramAggregation {
	return elastic.NewHistogramAggregation().
		Field(field).
		Interval(interval).
		MinDocCount(0)
}
