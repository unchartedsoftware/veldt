package elastic

import (
	"encoding/json"
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/elastic/param"
	"github.com/unchartedsoftware/prism/generation/tile"
)

// TopicCountTile represents a tiling generator that produces topic counts.
type TopicCountTile struct {
	Tiling    *param.Tiling
	Topic     *param.Topic
	TimeRange *param.TimeRange
}

// NewTopicCountTile instantiates and returns a pointer to a new generator.
func NewTopicCountTile(tileReq *tile.Request) (tile.Generator, error) {
	tiling, err := param.NewTiling(tileReq)
	if err != nil {
		return nil, err
	}
	topic, err := param.NewTopic(tileReq)
	if err != nil {
		return nil, err
	}
	time, _ := param.NewTimeRange(tileReq)
	return &TopicCountTile{
		Tiling:    tiling,
		Topic:     topic,
		TimeRange: time,
	}, nil
}

// GetParams returns a slice of tiling parameters.
func (g *TopicCountTile) GetParams() []tile.Param {
	return []tile.Param{
		g.Tiling,
		g.Topic,
		g.TimeRange,
	}
}

// GetTile returns the marshalled tile data.
func (g *TopicCountTile) GetTile(tileReq *tile.Request) ([]byte, error) {
	tiling := g.Tiling
	timeRange := g.TimeRange
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
	topicAggs := topic.GetTopicAggregations()
	for topic, topicAgg := range topicAggs {
		query.Aggregation(topic, topicAgg)
	}
	// send query
	result, err := query.Do()
	if err != nil {
		return nil, err
	}
	// build map of topics and counts
	topicCounts := make(map[string]int64)
	for _, topic := range topic.Topics {
		filter, ok := result.Aggregations.Filter(topic)
		if !ok {
			return nil, fmt.Errorf("Filter aggregation '%s' was not found in response for request %s", topic, tileReq.String())
		}
		if filter.DocCount > 0 {
			topicCounts[topic] = filter.DocCount
		}
	}
	// marshal results map
	return json.Marshal(topicCounts)
}
