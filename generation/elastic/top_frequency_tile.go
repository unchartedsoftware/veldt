package elastic

import (
	"encoding/json"
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/elastic/param"
	"github.com/unchartedsoftware/prism/generation/elastic/throttle"
	"github.com/unchartedsoftware/prism/generation/tile"
)

// TopFrequencyTile represents a tiling generator that produces top term
// frequency counts.
type TopFrequencyTile struct {
	TileGenerator
	Tiling    *param.Tiling
	TopTerms  *param.TopTerms
	Terms     *param.TermsFilter
	Prefixes  *param.PrefixFilter
	Range     *param.Range
	Time      *param.DateHistogram
	Histogram *param.Histogram
}

// NewTopFrequencyTile instantiates and returns a pointer to a new generator.
func NewTopFrequencyTile(host, port string) tile.GeneratorConstructor {
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
		time, err := param.NewDateHistogram(tileReq)
		if err != nil {
			return nil, err
		}
		terms, _ := param.NewTermsFilter(tileReq)
		prefixes, _ := param.NewPrefixFilter(tileReq)
		rang, _ := param.NewRange(tileReq)
		histogram, _ := param.NewHistogram(tileReq)
		t := &TopFrequencyTile{}
		t.Tiling = tiling
		t.TopTerms = topTerms
		t.Terms = terms
		t.Prefixes = prefixes
		t.Time = time
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
func (g *TopFrequencyTile) GetParams() []tile.Param {
	return []tile.Param{
		g.Tiling,
		g.Terms,
		g.TopTerms,
		g.Prefixes,
		g.Range,
		g.Time,
		g.Histogram,
	}
}

func (g *TopFrequencyTile) getQuery() elastic.Query {
	// optional filters
	filters := elastic.NewBoolQuery()
	// if range param is provided, add range queries
	if g.Range != nil {
		for _, query := range g.Range.GetQueries() {
			filters.Must(query)
		}
	}
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
	return elastic.NewBoolQuery().
		Must(g.Tiling.GetXQuery()).
		Must(g.Tiling.GetYQuery()).
		Must(g.Time.GetQuery()).
		Must(filters)
}

func (g *TopFrequencyTile) getAgg() elastic.Aggregation {
	// get top terms agg
	agg := g.TopTerms.GetAggregation()
	// get date histogram agg
	timeAgg := g.Time.GetAggregation()
	// if histogram param is provided, add histogram agg
	if g.Histogram != nil {
		timeAgg.SubAggregation(histogramAggName, g.Histogram.GetAggregation())
	}
	// add date histogram agg
	agg.SubAggregation(timeAggName, timeAgg)
	return agg
}

func (g *TopFrequencyTile) parseResult(res *elastic.SearchResult) ([]byte, error) {
	// build map of topics and frequency arrays
	frequencies := make(map[string][]interface{})
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
		time, ok := bucket.Aggregations.DateHistogram(timeAggName)
		if !ok {
			return nil, fmt.Errorf("DateHistogram aggregation '%s' was not found in response for request %s", timeAggName, g.req.String())
		}
		counts := make([]interface{}, len(time.Buckets))
		for i, bucket := range time.Buckets {
			if g.Histogram != nil {
				histogram, ok := bucket.Aggregations.Histogram(histogramAggName)
				if !ok {
					return nil, fmt.Errorf("Histogram aggregation '%s' was not found in response for request %s",
						histogramAggName,
						g.req.String())
				}
				counts[i] = g.Histogram.GetBucketMap(histogram)
			} else {
				counts[i] = bucket.DocCount
			}
		}
		// add counts to frequencies map
		frequencies[term] = counts
	}
	// marshal results map
	return json.Marshal(frequencies)
}

// GetTile returns the marshalled tile data.
func (g *TopFrequencyTile) GetTile() ([]byte, error) {
	// build query
	query := g.client.
		Search(g.req.Index).
		Size(0).
		Query(g.getQuery()).
		Aggregation(termsAggName, g.getAgg())
	// send query through equalizer
	res, err := throttle.Send(query)
	if err != nil {
		return nil, err
	}
	// parse and return results
	return g.parseResult(res);
}
