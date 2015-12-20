package elastic

import (
	"encoding/json"
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/elastic/param"
	"github.com/unchartedsoftware/prism/generation/tile"
)

const (
	timeAggName = "time"
)

// TopicFrequencyTile represents a tiling generator that produces topic
// frequency counts.
type TopicFrequencyTile struct {
	Tiling     *param.Tiling
	Topic      *param.Topic
	TimeBucket *param.TimeBucket
}

// NewTopicFrequencyTile instantiates and returns a pointer to a new generator.
func NewTopicFrequencyTile(tileReq *tile.Request) (tile.Generator, error) {
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
	return &TopicFrequencyTile{
		Tiling:     tiling,
		Topic:      topic,
		TimeBucket: time,
	}, nil
}

// GetParams returns a slice of tiling parameters.
func (g *TopicFrequencyTile) GetParams() []tile.Param {
	return []tile.Param{
		g.Tiling,
		g.Topic,
		g.TimeBucket,
	}
}

// GetTile returns the marshalled tile data.
func (g *TopicFrequencyTile) GetTile(tileReq *tile.Request) ([]byte, error) {
	tiling := g.Tiling
	timeBucket := g.TimeBucket
	timeRange := timeBucket.TimeRange
	topic := g.Topic
	// get client
	client, err := getClient(tileReq.Endpoint)
	if err != nil {
		return nil, err
	}
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
		query.Aggregation(topic, topicAgg.SubAggregation(timeAggName, timeAgg))
	}
	// send query
	result, err := query.Do()
	if err != nil {
		return nil, err
	}
	// build map of topics and frequency arrays
	topicFrequencies := make(map[string][]int64)
	for _, topic := range topic.Topics {
		filter, ok := result.Aggregations.Filter(topic)
		if !ok {
			return nil, fmt.Errorf("Filter aggregation '%s' was not found in response", topic)
		}
		timeAgg, ok := filter.DateHistogram(timeAggName)
		if !ok {
			return nil, fmt.Errorf("DateHistogram aggregation '%s' was not found in response", timeAggName)
		}
		topicCounts := make([]int64, len(timeAgg.Buckets))
		for _, bucket := range timeAgg.Buckets {
			topicCounts = append(topicCounts, bucket.DocCount)
		}
		topicFrequencies[topic] = topicCounts
	}
	// marshal results map
	return json.Marshal(topicFrequencies)
}
