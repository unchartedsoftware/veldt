package elastic

import (
	"encoding/json"
	"fmt"
	"sort"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/elastic/agg"
	"github.com/unchartedsoftware/prism/generation/elastic/param"
	"github.com/unchartedsoftware/prism/generation/elastic/query"
	"github.com/unchartedsoftware/prism/generation/tile"
)

// TopTrailsTile represents a tiling generator that produces heatmaps.
type TopTrailsTile struct {
	TileGenerator
	Binning *param.Binning
	Query   *query.Bool
	Terms   *agg.Terms
}

func rankByCount(counts map[string]int64) pairs {
	pl := make(pairs, len(counts))
	i := 0
	for k, v := range counts {
		pl[i] = pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

type pair struct {
	Key   string
	Value int64
}

type pairs []pair

func (p pairs) Len() int           { return len(p) }
func (p pairs) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p pairs) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// NewTopTrailsTile instantiates and returns a pointer to a new generator.
func NewTopTrailsTile(host, port string) tile.GeneratorConstructor {
	return func(tileReq *tile.Request) (tile.Generator, error) {
		client, err := NewClient(host, port)
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

func (g *TopTrailsTile) getQuery() elastic.Query {
	return elastic.NewBoolQuery().
		Must(g.Binning.Tiling.GetXQuery()).
		Must(g.Binning.Tiling.GetYQuery()).
		Must(g.Query.GetQuery())
}

func (g *TopTrailsTile) getAgg() elastic.Aggregation {
	// create x aggregation
	xAgg := g.Binning.GetXAgg()
	// create y aggregation, add it as a sub-agg to xAgg
	yAgg := g.Binning.GetYAgg()
	xAgg.SubAggregation(yAggName, yAgg)
	// create top hits aggregation
	yAgg.SubAggregation(termsAggName, g.Terms.GetAgg())
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
	counts := make(map[string]int64)
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
			// extract terms
			terms, ok := yBucket.Terms(termsAggName)
			if !ok {
				return nil, fmt.Errorf("Terms aggregation '%s' was not found in response for request %s",
					termsAggName,
					g.req.String())
			}
			// get term buckets
			for _, bucket := range terms.Buckets {
				key, ok := bucket.Key.(string)
				if !ok {
					return nil, fmt.Errorf("Term bucket key was not of type `string` for request %s",
						g.req.String())
				}
				// add bin location under key
				if bins[key] == nil {
					bins[key] = make(map[int64]map[int64]bool)
				}
				if bins[key][xBin] == nil {
					bins[key][xBin] = make(map[int64]bool)
				}
				bins[key][xBin][yBin] = true
				counts[key] += bucket.DocCount
			}
		}
	}
	// rank the counts
	ranked := rankByCount(counts)
	length := g.Terms.Size
	if len(ranked) < g.Terms.Size {
		length = len(ranked)
	}
	// create map of bin positions for top N docs
	top := make(map[string][][]int64)
	for i := 0; i < length; i++ {
		key := ranked[i].Key
		top[key] = make([][]int64, 0)
		for x, xs := range bins[key] {
			for y := range xs {
				top[key] = append(top[key], []int64{x, y})
			}
		}
	}
	// marshal results map
	return json.Marshal(top)
}

// GetTile returns the marshalled tile data.
func (g *TopTrailsTile) GetTile() ([]byte, error) {
	// build query
	query := g.client.
		Search(g.req.Index).
		Size(0).
		Query(g.getQuery()).
		Aggregation(xAggName, g.getAgg())
	// send query through equalizer
	res, err := query.Do()
	if err != nil {
		return nil, err
	}
	// parse and return results
	return g.parseResult(res)
}
