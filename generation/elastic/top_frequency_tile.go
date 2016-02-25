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

// GetTile returns the marshalled tile data.
func (g *TopFrequencyTile) GetTile() ([]byte, error) {
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
		for _, query := range g.Terms.GetQueries() {
			boolQuery.Must(query)
		}
	}
	// if prefixes param is provided, add prefix queries
	if g.Prefixes != nil {
		for _, query := range g.Prefixes.GetQueries() {
			boolQuery.Must(query)
		}
	}
	// add time range query
	boolQuery.Must(time.GetQuery())
	// get date histogram agg
	timeAgg := time.GetAggregation()
	// if histogram param is provided, add histogram agg
	if g.Histogram != nil {
		timeAgg.SubAggregation(histogramAggName, g.Histogram.GetAggregation())
	}
	// build query
	query := client.
		Search(tileReq.Index).
		Size(0).
		Query(boolQuery).
		Aggregation(termsAggName, g.TopTerms.GetAggregation().
		SubAggregation(timeAggName, timeAgg))
	// send query through equalizer
	result, err := throttle.Send(query)
	if err != nil {
		return nil, err
	}
	// build map of topics and frequency arrays
	topTermFrequencies := make(map[string][]interface{})
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
		timeAgg, ok := bucket.Aggregations.DateHistogram(timeAggName)
		if !ok {
			return nil, fmt.Errorf("DateHistogram aggregation '%s' was not found in response for request %s", timeAggName, tileReq.String())
		}
		termCounts := make([]interface{}, len(timeAgg.Buckets))
		for i, bucket := range timeAgg.Buckets {
			if g.Histogram != nil {
				histogramAgg, ok := bucket.Aggregations.Histogram(histogramAggName)
				if !ok {
					return nil, fmt.Errorf("Histogram aggregation '%s' was not found in response for request %s",
						histogramAggName,
						tileReq.String())
				}
				termCounts[i] = g.Histogram.GetBucketMap(histogramAgg)
			} else {
				termCounts[i] = bucket.DocCount
			}
		}
		// add counts to frequencies map
		topTermFrequencies[term] = termCounts
	}
	// marshal results map
	return json.Marshal(topTermFrequencies)
}
