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
	termsAggName = "topterms"
)

// TopTermsTile represents a tiling generator that produces topic counts.
type TopTermsTile struct {
	TileGenerator
	Tiling    *param.Tiling
	TopTerms  *param.TopTerms
	TimeRange *param.TimeRange
}

// NewTopTermsTile instantiates and returns a pointer to a new generator.
func NewTopTermsTile(host, port string) tile.GeneratorConstructor {
	return func(tileReq *tile.Request) (tile.Generator, error) {
		client, err := NewClient(host, port)
		if err != nil {
			return nil, err
		}
		tiling, err := param.NewTiling(tileReq)
		if err != nil {
			return nil, err
		}
		topTerms, err := param.NewTopTerms(tileReq)
		if err != nil {
			return nil, err
		}
		time, _ := param.NewTimeRange(tileReq)
		t := &TopTermsTile{}
		t.Tiling = tiling
		t.TopTerms = topTerms
		t.TimeRange = time
		t.req = tileReq
		t.host = host
		t.port = port
		t.client = client
		return t, nil
	}
}

// GetParams returns a slice of tiling parameters.
func (g *TopTermsTile) GetParams() []tile.Param {
	return []tile.Param{
		g.Tiling,
		g.TopTerms,
		g.TimeRange,
	}
}

// GetTile returns the marshalled tile data.
func (g *TopTermsTile) GetTile() ([]byte, error) {
	tiling := g.Tiling
	timeRange := g.TimeRange
	topTerms := g.TopTerms
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
		Query(boolQuery).
		Aggregation(termsAggName, topTerms.GetTermsAggregation())
	// send query through equalizer
	result, err := throttle.Send(query)
	if err != nil {
		return nil, err
	}
	// build map of topics and counts
	topTermCounts := make(map[string]int64)
	termsRes, ok := result.Aggregations.Terms(termsAggName)
	if !ok {
		return nil, fmt.Errorf("Terms aggregation '%s' was not found in response for request %s",
			termsAggName,
			tileReq.String())
	}
	for _, bucket := range termsRes.Buckets {
		term, ok := bucket.Key.(string)
		if !ok {
			return nil, fmt.Errorf("Terms aggregation key was not of type `string` '%s' in response for request %s",
				termsAggName,
				tileReq.String())
		}
		count := bucket.DocCount
		if count > 0 {
			topTermCounts[term] = count
		}
	}
	// marshal results map
	return json.Marshal(topTermCounts)
}
