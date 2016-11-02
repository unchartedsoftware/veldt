package elastic

import (
	"encoding/json"
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/elastic/agg"
	"github.com/unchartedsoftware/prism/generation/elastic/param"
	"github.com/unchartedsoftware/prism/generation/elastic/query"
	"github.com/unchartedsoftware/prism/tile"
)

// FrequencyTile represents a tiling generator that produces an array of counts.
type FrequencyTile struct {
	TileGenerator
	Tiling *param.Tiling
	Time   *agg.DateHistogram
	Query  *query.Bool
}

// FrequencyBin represents the result type in the returned array.
type FrequencyBin struct {
	Count     int64 `json:"count"`
	Timestamp int64 `json:"timestamp"`
}

// NewFrequencyTile instantiates and returns a pointer to a new generator.
func NewFrequencyTile(host, port string) tile.GeneratorConstructor {
	return func(tileReq *tile.Request) (tile.Generator, error) {
		client, err := NewClient(host, port)
		if err != nil {
			return nil, err
		}
		elastic, err := param.NewElastic(tileReq)
		if err != nil {
			return nil, err
		}
		// required
		tiling, err := param.NewTiling(tileReq)
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
		t := &FrequencyTile{}
		t.Elastic = elastic
		t.Tiling = tiling
		t.Time = time
		t.Query = query
		t.req = tileReq
		t.host = host
		t.port = port
		t.client = client
		return t, nil
	}
}

// GetParams returns a slice of tiling parameters.
func (g *FrequencyTile) GetParams() []tile.Param {
	return []tile.Param{
		g.Tiling,
		g.Time,
		g.Query,
	}
}

func (g *FrequencyTile) getQuery() elastic.Query {
	return elastic.NewBoolQuery().
		Must(g.Tiling.GetXQuery()).
		Must(g.Tiling.GetYQuery()).
		Must(g.Time.GetQuery()).
		Must(g.Query.GetQuery())
}

func (g *FrequencyTile) getAgg() elastic.Aggregation {
	// get date histogram agg
	return g.Time.GetAgg()
}

func (g *FrequencyTile) parseResult(res *elastic.SearchResult) ([]byte, error) {
	time, ok := res.Aggregations.DateHistogram(timeAggName)
	if !ok {
		return nil, fmt.Errorf("DateHistogram aggregation '%s' was not found in response for request %s",
			timeAggName,
			g.req.String())
	}
	bins := make([]FrequencyBin, len(time.Buckets))
	for i, bucket := range time.Buckets {
		bins[i] = FrequencyBin{
			Count:     bucket.DocCount,
			Timestamp: bucket.Key,
		}
	}
	// marshal results map
	return json.Marshal(bins)
}

// GetTile returns the marshalled tile data.
func (g *FrequencyTile) GetTile() ([]byte, error) {
	// send query
	res, err := g.Elastic.GetSearchService(g.client).
		Index(g.req.URI).
		Size(0).
		Query(g.getQuery()).
		Aggregation(timeAggName, g.getAgg()).
		Do()
	if err != nil {
		return nil, err
	}
	// parse and return results
	return g.parseResult(res)
}
