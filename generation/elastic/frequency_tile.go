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

// FrequencyTile represents a tiling generator that produces an array of counts.
type FrequencyTile struct {
	TileGenerator
	Tiling    *param.Tiling
	Time      *agg.DateHistogram
	Query     *query.Bool
	Histogram *agg.Histogram
}

// FrequencyBin represents the result type in the returned array. (what is json type for Timestamp?)
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
		// optional
		histogram, err := agg.NewHistogram(tileReq.Params)
		if param.IsOptionalErr(err) {
			return nil, err
		}
		t := &FrequencyTile{}
		t.Elastic = elastic
		t.Tiling = tiling
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
func (g *FrequencyTile) GetParams() []tile.Param {
	return []tile.Param{
		g.Tiling,
		g.Time,
		g.Query,
		g.Histogram,
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
	agg := g.Time.GetAgg()
	// if histogram param is provided, add histogram agg
	if g.Histogram != nil {
		agg.SubAggregation(histogramAggName, g.Histogram.GetAgg())
	}
	return agg
}

func (g *FrequencyTile) parseResult(res *elastic.SearchResult) ([]byte, error) {
	time, ok := res.Aggregations.DateHistogram(timeAggName)
	if !ok {
		return nil, fmt.Errorf("DateHistogram aggregation '%s' was not found in response for request %s", timeAggName, g.req.String())
	}
	bins := make([]FrequencyBin, len(time.Buckets))
	counts := make([]interface{}, len(time.Buckets))
	for i, bucket := range time.Buckets {
		if g.Histogram != nil {
			histogram, ok := bucket.Aggregations.Histogram(histogramAggName)
			if !ok {
				return nil, fmt.Errorf("Histogram aggregation '%s' was not found in response for request %s",
					histogramAggName,
					g.req.String())
			}
			// TODO What is this case for?
			counts[i] = g.Histogram.GetBucketMap(histogram)
		} else {
			counts[i] = bucket.DocCount
			bin := FrequencyBin{Count: bucket.DocCount, Timestamp: bucket.Key}
			bins[i] = bin
		}
	}
	// marshal results map
	return json.Marshal(bins)
}

// GetTile returns the marshalled tile data.
func (g *FrequencyTile) GetTile() ([]byte, error) {
	// temp debug query
	querySource, err := g.getQuery().Source()
	marshalledQuery, err := json.Marshal(querySource)
	fmt.Println("=== QUERY: " + string(marshalledQuery[:]) + "\n")

	aggSource, err := g.getAgg().Source()
	marshalledAgg, err := json.Marshal(aggSource)
	fmt.Println("=== AGG: " + string(marshalledAgg[:]) + "\n")

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
