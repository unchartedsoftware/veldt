package elastic

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/elastic/param"
	"github.com/unchartedsoftware/prism/generation/elastic/query"
	"github.com/unchartedsoftware/prism/generation/tile"
)

// MacroTile represents a tiling generator that produces a tile.
type MacroTile struct {
	TileGenerator
	Binning    *param.Binning
	Query      *query.Bool
	MacroMicro *param.MacroMicro
}

// NewMacroTile instantiates and returns a pointer to a new generator.
func NewMacroTile(host, port string) tile.GeneratorConstructor {
	return func(tileReq *tile.Request) (tile.Generator, error) {
		client, err := NewClient(host, port)
		if err != nil {
			return nil, err
		}
		binning, err := param.NewBinning(tileReq)
		if err != nil {
			return nil, err
		}
		macromicro, err := param.NewMacroMicro(tileReq)
		if err != nil {
			return nil, err
		}
		query, err := query.NewBool(tileReq.Params)
		if err != nil {
			return nil, err
		}
		t := &MacroTile{}
		t.Binning = binning
		t.MacroMicro = macromicro
		t.Query = query
		t.req = tileReq
		t.host = host
		t.port = port
		t.client = client
		return t, nil
	}
}

// GetParams returns a slice of tiling parameters.
func (g *MacroTile) GetParams() []tile.Param {
	return []tile.Param{
		g.Binning,
		g.MacroMicro,
		g.Query,
	}
}

func (g *MacroTile) getQuery() elastic.Query {
	return elastic.NewBoolQuery().
		Must(g.Binning.Tiling.GetXQuery()).
		Must(g.Binning.Tiling.GetYQuery()).
		Must(g.Query.GetQuery())
}

func (g *MacroTile) getAgg() elastic.Aggregation {
	// create x aggregation
	return g.Binning.GetXAgg().
		SubAggregation(yAggName, g.Binning.GetYAgg())
}

func (g *MacroTile) parseResult(res *elastic.SearchResult) ([]byte, error) {
	binning := g.Binning
	// parse aggregations
	xAggRes, ok := res.Aggregations.Histogram(xAggName)
	if !ok {
		return nil, fmt.Errorf("Histogram aggregation '%s' was not found in response for request %s",
			xAggName,
			g.req.String())
	}
	// allocate bins buffer
	bins := make([]float64, binning.Resolution*binning.Resolution)
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
			index := xBin + binning.Resolution*yBin
			// encode count
			bins[index] += float64(yBucket.DocCount)
		}
	}
	return float64ToByteSlice(bins), nil
}

// GetTile returns the marshalled tile data.
func (g *MacroTile) GetTile() ([]byte, error) {
	// first pass to get the count for the tile
	res, err := g.client.
		Search(g.req.Index).
		Size(0).
		Query(g.getQuery()).
		Do()
	if err != nil {
		return nil, err
	}
	if res.Hits.TotalHits > g.MacroMicro.Threshold {
		// generate macro tile
		res, err := g.client.
			Search(g.req.Index).
			Size(0).
			Query(g.getQuery()).
			Aggregation(xAggName, g.getAgg()).
			Do()
		if err != nil {
			return nil, err
		}
		// parse and return results
		return g.parseResult(res)
	}
	// below threshold, return nil
	return nil, nil
}
