package elastic

import (
	"encoding/json"
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/elastic/param"
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
	TopTerms  *param.TopTerms
	Terms     *param.TermsFilter
	Prefixes  *param.PrefixFilter
	Range     *param.Range
	Histogram *param.Histogram
}

// NewTopCountTile instantiates and returns a pointer to a new generator.
func NewTopCountTile(host, port string) tile.GeneratorConstructor {
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
		terms, err := param.NewTermsFilter(tileReq)
		if param.IsOptionalErr(err) {
			return nil, err
		}
		prefixes, err := param.NewPrefixFilter(tileReq)
		if param.IsOptionalErr(err) {
			return nil, err
		}
		rang, err := param.NewRange(tileReq)
		if param.IsOptionalErr(err) {
			return nil, err
		}
		histogram, err := param.NewHistogram(tileReq)
		if param.IsOptionalErr(err) {
			return nil, err
		}
		t := &TopCountTile{}
		t.Tiling = tiling
		t.TopTerms = topTerms
		t.Range = rang
		t.Terms = terms
		t.Prefixes = prefixes
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
		g.Prefixes,
		g.Terms,
		g.Range,
		g.Histogram,
	}
}

func (g *TopCountTile) getQuery() elastic.Query {
	// optional filters
	filters := elastic.NewBoolQuery()
	// if range param is provided, add range queries
	if g.Range != nil {
		for _, query := range g.Range.GetQueries() {
			filters.Must(query)
		}
	}
	// the following filters need to be wrapped in a `must` otherwise the
	// above `must` query will override them.
	if g.Terms != nil || g.Prefixes != nil {
		// create sub-filter
		subfilters := elastic.NewBoolQuery()
		// if terms param is provided, add terms queries
		if g.Terms != nil {
			for _, query := range g.Terms.GetQueries() {
				filters.Should(query)
			}
		}
		// if prefixes param is provided, add prefix queries
		if g.Prefixes != nil {
			for _, query := range g.Prefixes.GetQueries() {
				filters.Should(query)
			}
		}
		// add sub-filter to the parent filter
		filters.Must(subfilters)
	}
	return elastic.NewBoolQuery().
		Must(g.Tiling.GetXQuery()).
		Must(g.Tiling.GetYQuery()).
		Must(filters)
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
