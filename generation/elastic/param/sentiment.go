package param

import (
	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/tile"
)

const (
	field    = "sentiment"
	interval = 1
)

// SentimentCounts represents counts for each sentiment.
type SentimentCounts struct {
	Positive uint64 `json:"positive"`
	Neutral  uint64 `json:"neutral"`
	Negative uint64 `json:"negative"`
}

// Sentiment represents params for extracting sentiment counts.
type Sentiment struct {
}

// NewSentiment instantiates and returns a new sentiment parameter object.
func NewSentiment(tileReq *tile.Request) (*Sentiment, error) {
	return &Sentiment{}, nil
}

// GetHash returns a string hash of the parameter state.
func (p *Sentiment) GetHash() string {
	return "sentiment"
}

// GetSentimentAgg returns an elastic query.
func (p *Sentiment) GetSentimentAgg() *elastic.HistogramAggregation {
	return elastic.NewHistogramAggregation().
		Field(field).
		Interval(interval).
		MinDocCount(0)
}

// GetSentimentCounts returns the sentiment counts from the histogram aggregation object.
func (p *Sentiment) GetSentimentCounts(sentimentAgg *elastic.AggregationBucketHistogramItems) *SentimentCounts {
	counts := &SentimentCounts{}
	for _, bucket := range sentimentAgg.Buckets {
		switch bucket.Key {
		case 1:
			counts.Positive = uint64(bucket.DocCount)
		case 0:
			counts.Neutral = uint64(bucket.DocCount)
		case -1:
			counts.Negative = uint64(bucket.DocCount)
		}
	}
	return counts
}
