package elastic

import (
	"encoding/json"
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/elastic/param"
	"github.com/unchartedsoftware/prism/generation/elastic/throttle"
	"github.com/unchartedsoftware/prism/generation/tile"
)

// TopicSentimentFrequencyTile represents a tiling generator that produces topic
// frequency counts.
type TopicSentimentFrequencyTile struct {
	TileGenerator
	Tiling     *param.Tiling
	Topic      *param.Topic
	TimeBucket *param.TimeBucket
	Sentiment  *param.Sentiment
}

// NewTopicSentimentFrequencyTile instantiates and returns a pointer to a new generator.
func NewTopicSentimentFrequencyTile(host, port string) tile.GeneratorConstructor {
	return func(tileReq *tile.Request) (tile.Generator, error) {
		client, err := NewClient(host, port)
		if err != nil {
			return nil, err
		}
		tiling, err := param.NewTiling(tileReq)
		if err != nil {
			return nil, err
		}
		topic, err := param.NewTopic(tileReq)
		if err != nil {
			return nil, err
		}
		time, err := param.NewTimeBucket(tileReq)
		if err != nil {
			return nil, err
		}
		sentiment, err := param.NewSentiment(tileReq)
		if err != nil {
			return nil, err
		}
		t := &TopicSentimentFrequencyTile{}
		t.Tiling = tiling
		t.Topic = topic
		t.TimeBucket = time
		t.Sentiment = sentiment
		t.req = tileReq
		t.host = host
		t.port = port
		t.client = client
		return t, nil
	}
}

// GetParams returns a slice of tiling parameters.
func (g *TopicSentimentFrequencyTile) GetParams() []tile.Param {
	return []tile.Param{
		g.Tiling,
		g.Topic,
		g.TimeBucket,
		g.Sentiment,
	}
}

// GetTile returns the marshalled tile data.
func (g *TopicSentimentFrequencyTile) GetTile() ([]byte, error) {
	tiling := g.Tiling
	timeBucket := g.TimeBucket
	timeRange := timeBucket.TimeRange
	topic := g.Topic
	sentiment := g.Sentiment
	tileReq := g.req
	client := g.client
	// create x and y range queries
	boolQuery := elastic.NewBoolQuery().Must(
		tiling.GetXQuery(),
		tiling.GetYQuery())
	// if time params are provided, add time range query
	if timeRange != nil {
		boolQuery.Must(timeRange.GetTimeQuery())
	}
	// build query
	query := client.
		Search(tileReq.Index).
		Size(0).
		Query(boolQuery)
	// add all filter aggregations
	timeAgg := timeBucket.GetTimeAggregation()
	topicAggs := topic.GetTopicAggregations()
	for topic, topicAgg := range topicAggs {
		query.Aggregation(topic,
			topicAgg.SubAggregation(timeAggName,
				timeAgg.SubAggregation(sentimentAggName, sentiment.GetSentimentAgg())))
	}
	// send query through equalizer
	result, err := throttle.Send(query)
	if err != nil {
		return nil, err
	}
	// build map of topics and frequency arrays
	topicFrequencies := make(map[string][]*param.SentimentCounts)
	for _, topic := range topic.Topics {
		filter, ok := result.Aggregations.Filter(topic)
		if !ok {
			return nil, fmt.Errorf("Filter aggregation '%s' was not found in response for request %s", topic, tileReq.String())
		}
		timeAgg, ok := filter.DateHistogram(timeAggName)
		if !ok {
			return nil, fmt.Errorf("DateHistogram aggregation '%s' was not found in response for request %s", timeAggName, tileReq.String())
		}
		topicCounts := make([]*param.SentimentCounts, len(timeAgg.Buckets))
		topicExists := false
		for i, bucket := range timeAgg.Buckets {
			if bucket.DocCount > 0 {
				topicExists = true
				sentimentAgg, ok := filter.Aggregations.Histogram(sentimentAggName)
				if !ok {
					return nil, fmt.Errorf("Histogram aggregation '%s' was not found in response for request %s", sentimentAggName, tileReq.String())
				}
				// extract sentiment counts
				topicCounts[i] = sentiment.GetSentimentCounts(sentimentAgg)
			}
		}
		// only add topics if they have at least one count
		if topicExists {
			topicFrequencies[topic] = topicCounts
		}
	}
	// marshal results map
	return json.Marshal(topicFrequencies)
}
