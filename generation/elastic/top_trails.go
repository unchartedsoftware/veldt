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

// TopTrailsTile represents a tiling generator that produces a top trails tile.
type TopTrailsTile struct {
	TileGenerator
	Binning *param.Binning
	Query   *query.Bool
	Terms   *agg.Terms
	// NOTE: this param is generated internally, rather than inside the request
	Filter *agg.TermsFilter
}

// NewTopTrailsTile instantiates and returns a pointer to a new generator.
func NewTopTrailsTile(host, port string) tile.GeneratorConstructor {
	return func(tileReq *tile.Request) (tile.Generator, error) {
		client, err := NewClient(host, port)
		if err != nil {
			return nil, err
		}
		elastic, err := param.NewElastic(tileReq)
		if err != nil {
			return nil, err
		}
		binning, err := param.NewBinning(tileReq)
		if err != nil {
			return nil, err
		}
		query, err := query.NewBool(tileReq.Params)
		if err != nil {
			return nil, err
		}
		terms, err := agg.NewTerms(tileReq.Params)
		if err != nil {
			return nil, err
		}
		t := &TopTrailsTile{}
		t.Elastic = elastic
		t.Binning = binning
		t.Query = query
		t.Terms = terms
		t.req = tileReq
		t.host = host
		t.port = port
		t.client = client
		return t, nil
	}
}

// GetParams returns a slice of tiling parameters.
func (g *TopTrailsTile) GetParams() []tile.Param {
	return []tile.Param{
		g.Binning,
		g.Query,
		g.Terms,
	}
}

func (g *TopTrailsTile) getFirstQuery() elastic.Query {
	return elastic.NewBoolQuery().
		Must(g.Binning.Tiling.GetXQuery()).
		Must(g.Binning.Tiling.GetYQuery()).
		Must(g.Query.GetQuery())
}

func (g *TopTrailsTile) getSecondQuery() elastic.Query {
	return elastic.NewBoolQuery().
		Must(g.Binning.Tiling.GetXQuery()).
		Must(g.Binning.Tiling.GetYQuery()).
		Must(g.Filter.GetQuery()).
		Must(g.Query.GetQuery())
}

func (g *TopTrailsTile) getAgg() elastic.Aggregation {
	// create x aggregation
	xAgg := g.Binning.GetXAgg()
	// create y aggregation, add it as a sub-agg to xAgg
	yAgg := g.Binning.GetYAgg()
	xAgg.SubAggregation(yAggName, yAgg)
	// add all filter aggregations
	for id, agg := range g.Filter.GetAggs() {
		yAgg.SubAggregation(id, agg)
	}
	return xAgg
}

func (g *TopTrailsTile) parseResult(res *elastic.SearchResult) ([]byte, error) {
	binning := g.Binning
	// parse aggregations
	xAggRes, ok := res.Aggregations.Histogram(xAggName)
	if !ok {
		return nil, fmt.Errorf("Histogram aggregation '%s' was not found in response for request %s",
			xAggName,
			g.req.String())
	}
	// the bins coords per document key
	bins := make(map[string]map[int64]map[int64]bool)
	// fill bins buffer
	for _, xBucket := range xAggRes.Buckets {
		x := xBucket.Key
		xBin := binning.GetXBin(x)
		yAggRes, ok := xBucket.Histogram(yAggName)
		if !ok {
			return nil, fmt.Errorf("Histogram aggregation '%s' was not found in response for request %s",
				yAggName,
				g.req.String())
		}
		for _, yBucket := range yAggRes.Buckets {
			y := yBucket.Key
			yBin := binning.GetYBin(y)
			// extract ids
			for _, id := range g.Filter.Terms {
				filter, ok := yBucket.Aggregations.Filter(id)
				if !ok {
					return nil, fmt.Errorf("Filter aggregation '%s' was not found in response for request %s",
						id,
						g.req.String())
				}
				if filter.DocCount > 0 {
					// add bin location under key
					if bins[id] == nil {
						bins[id] = make(map[int64]map[int64]bool)
					}
					if bins[id][xBin] == nil {
						bins[id][xBin] = make(map[int64]bool)
					}
					bins[id][xBin][yBin] = true
				}
			}
		}
	}
	// create map of bin positions for top N docs
	top := make(map[string][][]int64)
	for _, id := range g.Filter.Terms {
		bin := bins[id]
		top[id] = make([][]int64, 0)
		for x, xs := range bin {
			for y := range xs {
				top[id] = append(top[id], []int64{x, y})
			}
		}
	}
	// marshal results map
	return json.Marshal(top)
}

// GetTile returns the marshalled tile data.
func (g *TopTrailsTile) GetTile() ([]byte, error) {
	// first pass to get the top N ids
	res, err := g.Elastic.GetSearchService(g.client).
		Index(g.req.URI).
		Size(0).
		Query(g.getFirstQuery()).
		Aggregation(termsAggName, g.Terms.GetAgg()).
		Do()
	if err != nil {
		return nil, err
	}
	terms, ok := res.Aggregations.Terms(termsAggName)
	if !ok {
		return nil, fmt.Errorf("Terms aggregation '%s' was not found in response for request %s",
			termsAggName,
			g.req.String())
	}
	// if no ids exit early
	if len(terms.Buckets) == 0 {
		return json.Marshal(make(map[string]interface{}))
	}
	// other wise, let's get the bin locations
	top := make([]string, len(terms.Buckets))
	// get term buckets
	for i, bucket := range terms.Buckets {
		top[i] = bucket.Key.(string)
	}
	// make id filter agg
	g.Filter = &agg.TermsFilter{
		Field: g.Terms.Field,
		Terms: top,
	}
	// second pass to pull by bin
	res, err = g.Elastic.GetSearchService(g.client).
		Index(g.req.URI).
		Size(0).
		Query(g.getSecondQuery()).
		Aggregation(xAggName, g.getAgg()).
		Do()
	if err != nil {
		return nil, err
	}
	// parse and return results
	return g.parseResult(res)
}
