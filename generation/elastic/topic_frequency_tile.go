package elastic

import (
	"encoding/json"
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/elastic/agg"
	"github.com/unchartedsoftware/prism/generation/elastic/param"
	"github.com/unchartedsoftware/prism/generation/elastic/query"
	"github.com/unchartedsoftware/prism/generation/tile"
)

const (
	timeAggName = "time"
)

// TopicFrequencyTile represents a tiling generator that produces term
// frequency counts.
type TopicFrequencyTile struct {
	TileGenerator
	Tiling    *param.Tiling
	Terms     *agg.Terms
	Time      *agg.DateHistogram
	Query     *query.Bool
	Histogram *agg.Histogram
}

// NewTopicFrequencyTile instantiates and returns a pointer to a new generator.
func NewTopicFrequencyTile(host, port string) tile.GeneratorConstructor {
	return func(tileReq *tile.Request) (tile.Generator, error) {
		client, err := NewClient(host, port)
		if err != nil {
			return nil, err
		}
		// required
		tiling, err := param.NewTiling(tileReq)
		if err != nil {
			return nil, err
		}
		terms, err := agg.NewTerms(tileReq.Params)
		if err != nil {
			return nil, err
		}
		time, err := agg.NewDateHistogram(tileReq.Params)
		if err != nil {
			return nil, err
		}
		query, err := query.NewBool(tileReq.Params)
		if err != nil {
			return nil, err
		}
		// optional
		histogram, err := agg.NewHistogram(tileReq.Params)
		if param.IsOptionalErr(err) {
			return nil, err
		}
		t := &TopicFrequencyTile{}
		t.Tiling = tiling
		t.Terms = terms
		t.Time = time
		t.Query = query
		t.Histogram = histogram
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
		g.Time,
		g.Query,
		g.Histogram,
	}
}

func (g *TopicFrequencyTile) getQuery() elastic.Query {
	return elastic.NewBoolQuery().
		Must(g.Tiling.GetXQuery()).
		Must(g.Tiling.GetYQuery()).
		Must(g.Time.GetQuery()).
		Must(g.Query.GetQuery())
}

func (g *TopicFrequencyTile) addAggs(query *elastic.SearchService) *elastic.SearchService {
	// date histogram aggregation
	timeAgg := g.Time.GetAggregation()
	// if histogram param is provided, add histogram agg
	if g.Histogram != nil {
		timeAgg.SubAggregation(histogramAggName, g.Histogram.GetAggregation())
	}
	// add all filter aggregations
	for term, termAgg := range g.Terms.GetAggregations() {
		query.Aggregation(term, termAgg.SubAggregation(timeAggName, timeAgg))
	}
	return query
}

func (g *TopicFrequencyTile) parseResult(res *elastic.SearchResult) ([]byte, error) {
	// build map of topics and frequency arrays
	frequencies := make(map[string][]interface{})
	for _, term := range g.Terms.Terms {
		filter, ok := res.Aggregations.Filter(term)
		if !ok {
			return nil, fmt.Errorf("Filter aggregation '%s' was not found in response for request %s",
				term,
				g.req.String())
		}
		timeAgg, ok := filter.DateHistogram(timeAggName)
		if !ok {
			return nil, fmt.Errorf("DateHistogram aggregation '%s' was not found in response for request %s",
				timeAggName,
				g.req.String())
		}
		counts := make([]interface{}, len(timeAgg.Buckets))
		topicExists := false
		for i, bucket := range timeAgg.Buckets {
			if g.Histogram != nil {
				histogramAgg, ok := bucket.Aggregations.Histogram(histogramAggName)
				if !ok {
					return nil, fmt.Errorf("Histogram aggregation '%s' was not found in response for request %s",
						histogramAggName,
						g.req.String())
				}
				counts[i] = g.Histogram.GetBucketMap(histogramAgg)
			} else {
				counts[i] = bucket.DocCount
			}
			if bucket.DocCount > 0 {
				topicExists = true
			}
		}
		// only add topics if they have at least one count
		if topicExists {
			frequencies[term] = counts
		}
	}
	// marshal results map
	return json.Marshal(frequencies)
}

// GetTile returns the marshalled tile data.
func (g *TopicFrequencyTile) GetTile() ([]byte, error) {
	// build query
	query := g.client.
		Search(g.req.Index).
		Size(0).
		Query(g.getQuery())
	// add all filter aggregations
	query = g.addAggs(query)
	// send query through equalizer
	res, err := query.Do()
	if err != nil {
		return nil, err
	}
	// parse and return results
	return g.parseResult(res)
}
