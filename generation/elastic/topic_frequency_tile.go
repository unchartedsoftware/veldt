package elastic

import (
	"encoding/json"
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/elastic/param"
	"github.com/unchartedsoftware/prism/generation/elastic/throttle"
	"github.com/unchartedsoftware/prism/generation/tile"
)

const (
	timeAggName = "time"
)

// TopicFrequencyTile represents a tiling generator that produces term
// frequency counts.
type TopicFrequencyTile struct {
	TileGenerator
	Tiling *param.Tiling
	Terms  *param.TermsAgg
	Range  *param.Range
	Time   *param.DateHistogram
}

// NewTopicFrequencyTile instantiates and returns a pointer to a new generator.
func NewTopicFrequencyTile(host, port string) tile.GeneratorConstructor {
	return func(tileReq *tile.Request) (tile.Generator, error) {
		client, err := NewClient(host, port)
		if err != nil {
			return nil, err
		}
		tiling, err := param.NewTiling(tileReq)
		if err != nil {
			return nil, err
		}
		terms, err := param.NewTermsAgg(tileReq)
		if err != nil {
			return nil, err
		}
		time, err := param.NewDateHistogram(tileReq)
		if err != nil {
			return nil, err
		}
		rang, _ := param.NewRange(tileReq)
		t := &TopicFrequencyTile{}
		t.Tiling = tiling
		t.Terms = terms
		t.Time = time
		t.Range = rang
		t.req = tileReq
		t.host = host
		t.port = port
		t.client = client
		return t, nil
	}
}

// GetParams returns a slice of tiling parameters.
func (g *TopicFrequencyTile) GetParams() []tile.Param {
	return []tile.Param{
		g.Tiling,
		g.Terms,
		g.Range,
		g.Time,
	}
}

// GetTile returns the marshalled tile data.
func (g *TopicFrequencyTile) GetTile() ([]byte, error) {
	tiling := g.Tiling
	time := g.Time
	tileReq := g.req
	client := g.client
	// create x and y range queries
	boolQuery := elastic.NewBoolQuery().Must(
		tiling.GetXQuery(),
		tiling.GetYQuery())
	// if range param is provided, add range queries
	if g.Range != nil {
		for _, query := range g.Range.GetQueries() {
			boolQuery.Must(query)
		}
	}
	// if terms param is provided, add terms queries
	if g.Terms != nil {
		boolQuery.Must(g.Terms.GetQuery())
	}
	// add time range query
	boolQuery.Must(time.GetQuery())
	// build query
	query := client.
		Search(tileReq.Index).
		Size(0).
		Query(boolQuery)
	// add all filter aggregations
	timeAgg := time.GetAggregation()
	termAggs := g.Terms.GetAggregations()
	for term, termAgg := range termAggs {
		query.Aggregation(term, termAgg.SubAggregation(timeAggName, timeAgg))
	}
	// send query through equalizer
	result, err := throttle.Send(query)
	if err != nil {
		return nil, err
	}
	// build map of topics and frequency arrays
	termFrequencies := make(map[string][]int64)
	for _, term := range g.Terms.Terms {
		filter, ok := result.Aggregations.Filter(term)
		if !ok {
			return nil, fmt.Errorf("Filter aggregation '%s' was not found in response for request %s", term, tileReq.String())
		}
		timeAgg, ok := filter.DateHistogram(timeAggName)
		if !ok {
			return nil, fmt.Errorf("DateHistogram aggregation '%s' was not found in response for request %s", timeAggName, tileReq.String())
		}
		termCounts := make([]int64, len(timeAgg.Buckets))
		topicExists := false
		for i, bucket := range timeAgg.Buckets {
			termCounts[i] = bucket.DocCount
			if bucket.DocCount > 0 {
				topicExists = true
			}
		}
		// only add topics if they have at least one count
		if topicExists {
			termFrequencies[term] = termCounts
		}
	}
	// marshal results map
	return json.Marshal(termFrequencies)
}
