package elastic

import (
	"encoding/json"
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/elastic/param"
	"github.com/unchartedsoftware/prism/generation/elastic/throttle"
	"github.com/unchartedsoftware/prism/generation/tile"
)

// TopicCountTile represents a tiling generator that produces terms counts.
type TopicCountTile struct {
	TileGenerator
	Tiling    *param.Tiling
	Terms     *param.TermsAgg
	Range     *param.Range
	Histogram *param.Histogram
}

// NewTopicCountTile instantiates and returns a pointer to a new generator.
func NewTopicCountTile(host, port string) tile.GeneratorConstructor {
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
		rang, _ := param.NewRange(tileReq)
		histogram, _ := param.NewHistogram(tileReq)
		t := &TopicCountTile{}
		t.Tiling = tiling
		t.Terms = terms
		t.Range = rang
		t.Histogram = histogram
		t.req = tileReq
		t.host = host
		t.port = port
		t.client = client
		return t, nil
	}
}

// GetParams returns a slice of tiling parameters.
func (g *TopicCountTile) GetParams() []tile.Param {
	return []tile.Param{
		g.Tiling,
		g.Terms,
		g.Range,
		g.Histogram,
	}
}

// GetTile returns the marshalled tile data.
func (g *TopicCountTile) GetTile() ([]byte, error) {
	tiling := g.Tiling
	terms := g.Terms
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
	// if terms param is provided, add terms query
	if g.Terms != nil {
		boolQuery.Should(g.Terms.GetQuery())
	}
	// build query
	query := client.
		Search(tileReq.Index).
		Size(0).
		Query(boolQuery)
	// add all filter aggregations
	termsAggs := g.Terms.GetAggregations()
	for term, termAgg := range termsAggs {
		// if histogram param is provided, add histogram agg
		if g.Histogram != nil {
			termAgg.SubAggregation(histogramAggName, g.Histogram.GetAggregation())
		}
		query.Aggregation(term, termAgg)
	}
	// send query through equalizer
	result, err := throttle.Send(query)
	if err != nil {
		return nil, err
	}
	// build map of topics and counts
	termCounts := make(map[string]interface{})
	for _, term := range terms.Terms {
		filter, ok := result.Aggregations.Filter(term)
		if !ok {
			return nil, fmt.Errorf("Filter aggregation '%s' was not found in response for request %s", term, tileReq.String())
		}
		if filter.DocCount > 0 {
			if g.Histogram != nil {
				histogramAgg, ok := filter.Aggregations.Histogram(histogramAggName)
				if !ok {
					return nil, fmt.Errorf("Histogram aggregation '%s' was not found in response for request %s",
						histogramAggName,
						tileReq.String())
				}
				termCounts[term] = g.Histogram.GetBucketMap(histogramAgg)
			} else {
				termCounts[term] = filter.DocCount
			}
		}
	}
	// marshal results map
	return json.Marshal(termCounts)
}
