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

// CustomTile represents a tiling generator that produces top term counts.
type CustomTile struct {
	TileGenerator
	Tiling     *param.Tiling
	Query      *query.Bool
	CustomAggs  *agg.CustomAggs
}

// NewCustomTile instantiates and returns a pointer to a new generator.
func NewCustomTile(host, port string) tile.GeneratorConstructor {
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
		query, err := query.NewBool(tileReq.Params)
		if err != nil {
			return nil, err
		}
		customAggs, err := agg.NewCustomAggs(tileReq.Params)
		if err != nil {
			return nil, err
		}
		t := &CustomTile{}
		t.Elastic = elastic
		t.Tiling = tiling
		t.CustomAggs = customAggs
		t.Query = query
		t.req = tileReq
		t.host = host
		t.port = port
		t.client = client
		return t, nil
	}
}

// GetParams returns a slice of tiling parameters.
func (g *CustomTile) GetParams() []tile.Param {
	return []tile.Param{
		g.Tiling,
		g.Query,
		g.CustomAggs,
	}
}

func (g *CustomTile) getQuery() elastic.Query {
	return elastic.NewBoolQuery().
		Must(g.Tiling.GetXQuery()).
		Must(g.Tiling.GetYQuery()).
		Must(g.Query.GetQuery())
}

func (g *CustomTile) GetQuerySource() (interface{}, error) {
	querySource, err := g.getQuery().Source()
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"query": querySource,
		"aggs": g.CustomAggs.GetAgg(),
	}, nil
}

func (g *CustomTile) parseResult(res *elastic.SearchResult) ([]byte, error) {
	// Return the raw results
	return json.Marshal(res)
}

// GetTile returns the marshalled tile data.
func (g *CustomTile) GetTile() ([]byte, error) {
	source, err := g.GetQuerySource()
	if err != nil {
		return nil, err
	}
	msg, err := json.Marshal(source)
	fmt.Println(string(msg[:]))
	// send query
	res, err := g.Elastic.GetSearchService(g.client).
		Index(g.req.URI).
		Source(source).
		Size(0).
		Do()
	if err != nil {
		return nil, err
	}
	// parse and return results
	return g.parseResult(res)
}
