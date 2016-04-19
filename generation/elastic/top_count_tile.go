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
	termsAggName     = "topterms"
	histogramAggName = "histogramAgg"
)

// TopCountTile represents a tiling generator that produces top term counts.
type TopCountTile struct {
	TileGenerator
	Tiling    *param.Tiling
	TopTerms  *agg.TopTerms
	Query     *query.Bool
	Histogram *agg.Histogram
}

// NewTopCountTile instantiates and returns a pointer to a new generator.
func NewTopCountTile(host, port string) tile.GeneratorConstructor {
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
		topTerms, err := agg.NewTopTerms(tileReq.Params)
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
		t := &TopCountTile{}
		t.Tiling = tiling
		t.TopTerms = topTerms
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
func (g *TopCountTile) GetParams() []tile.Param {
	return []tile.Param{
		g.Tiling,
		g.TopTerms,
		g.Query,
		g.Histogram,
	}
}

func (g *TopCountTile) getQuery() elastic.Query {
	return elastic.NewBoolQuery().
		Must(g.Tiling.GetXQuery()).
		Must(g.Tiling.GetYQuery()).
		Must(g.Query.GetQuery())
}

func (g *TopCountTile) getAgg() elastic.Aggregation {
	// get top terms agg
	agg := g.TopTerms.GetAggregation()
	// if histogram param is provided, add histogram agg
	if g.Histogram != nil {
		agg.SubAggregation(histogramAggName, g.Histogram.GetAggregation())
	}
	return agg
}

func (g *TopCountTile) parseResult(res *elastic.SearchResult) ([]byte, error) {
	// build map of topics and counts
	counts := make(map[string]interface{})
	terms, ok := res.Aggregations.Terms(termsAggName)
	if !ok {
		return nil, fmt.Errorf("Terms aggregation '%s' was not found in response for request %s",
			termsAggName,
			g.req.String())
	}
	for _, bucket := range terms.Buckets {
		term, ok := bucket.Key.(string)
		if !ok {
			return nil, fmt.Errorf("Terms aggregation key was not of type `string` '%s' in response for request %s",
				termsAggName,
				g.req.String())
		}
		if g.Histogram != nil {
			histogramAgg, ok := bucket.Aggregations.Histogram(histogramAggName)
			if !ok {
				return nil, fmt.Errorf("Histogram aggregation '%s' was not found in response for request %s",
					histogramAggName,
					g.req.String())
			}
			counts[term] = g.Histogram.GetBucketMap(histogramAgg)
		} else {
			counts[term] = bucket.DocCount
		}
	}
	// marshal results map
	return json.Marshal(counts)
}

// GetTile returns the marshalled tile data.
func (g *TopCountTile) GetTile() ([]byte, error) {
	// build query
	query := g.client.
		Search(g.req.Index).
		Size(0).
		Query(g.getQuery()).
		Aggregation(termsAggName, g.getAgg())
	// send query through equalizer
	res, err := query.Do()
	if err != nil {
		return nil, err
	}
	// parse and return results
	return g.parseResult(res)
}
