package elastic

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/elastic/param"
	"github.com/unchartedsoftware/prism/generation/elastic/query"
	"github.com/unchartedsoftware/prism/tile"
)

// CountTile represents a tiling generator that produces a tile.
type CountTile struct {
	TileGenerator
	Binning *param.Binning
	Query   *query.Bool
}

// NewCountTile instantiates and returns a pointer to a new generator.
func NewCountTile(host, port string) tile.GeneratorConstructor {
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
		t := &CountTile{}
		t.Elastic = elastic
		t.Binning = binning
		t.Query = query
		t.req = tileReq
		t.host = host
		t.port = port
		t.client = client
		return t, nil
	}
}

// GetParams returns a slice of tiling parameters.
func (g *CountTile) GetParams() []tile.Param {
	return []tile.Param{
		g.Binning,
		g.Query,
	}
}

func (g *CountTile) getQuery() elastic.Query {
	return elastic.NewBoolQuery().
		Must(g.Binning.Tiling.GetXQuery()).
		Must(g.Binning.Tiling.GetYQuery()).
		Must(g.Query.GetQuery())
}

// GetTile returns the marshalled tile data.
func (g *CountTile) GetTile() ([]byte, error) {
	res, err := g.Elastic.GetSearchService(g.client).
		Index(g.req.URI).
		Size(0).
		Query(g.getQuery()).
		Do()
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("{\"count\":%d}\n", res.Hits.TotalHits)), nil
}
