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

// MacroMicroTile represents a tiling generator that produces a tile.
type MacroMicroTile struct {
	TileGenerator
	Binning    *param.Binning
	Query      *query.Bool
	MacroMicro *param.MacroMicro
	TopHits    *agg.TopHits
}

// NewMacroMicroTile instantiates and returns a pointer to a new generator.
func NewMacroMicroTile(host, port string) tile.GeneratorConstructor {
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
		topHits, err := agg.NewTopHits(tileReq.Params)
		if err != nil {
			return nil, err
		}
		t := &MacroMicroTile{}
		t.Binning = binning
		t.MacroMicro = macromicro
		t.Query = query
		t.TopHits = topHits
		t.req = tileReq
		t.host = host
		t.port = port
		t.client = client
		return t, nil
	}
}

// GetParams returns a slice of tiling parameters.
func (g *MacroMicroTile) GetParams() []tile.Param {
	return []tile.Param{
		g.Binning,
		g.MacroMicro,
		g.Query,
		g.TopHits,
	}
}

func (g *MacroMicroTile) getQuery() elastic.Query {
	return elastic.NewBoolQuery().
		Must(g.Binning.Tiling.GetXQuery()).
		Must(g.Binning.Tiling.GetYQuery()).
		Must(g.Query.GetQuery())
}

func (g *MacroMicroTile) getMacroAgg() elastic.Aggregation {
	// create x aggregation
	return g.Binning.GetXAgg().
		SubAggregation(yAggName, g.Binning.GetYAgg())
}

func (g *MacroMicroTile) parseMacroResult(res *elastic.SearchResult) ([]byte, error) {
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

func (g *MacroMicroTile) getMicroAgg() elastic.Aggregation {
	// create x aggregation
	xAgg := g.Binning.GetXAgg()
	// create y aggregation, add it as a sub-agg to xAgg
	yAgg := g.Binning.GetYAgg()
	xAgg.SubAggregation(yAggName, yAgg)
	// if there is a z field to sum, add sum agg to yAgg
	yAgg.SubAggregation(topHitsAggName, g.TopHits.GetAgg())
	return xAgg
}

func (g *MacroMicroTile) parseMicroResult(res *elastic.SearchResult) ([]byte, error) {
	binning := g.Binning
	// parse aggregations
	xAggRes, ok := res.Aggregations.Histogram(xAggName)
	if !ok {
		return nil, fmt.Errorf("Histogram aggregation '%s' was not found in response for request %s",
			xAggName,
			g.req.String())
	}
	// allocate bins buffer
	bins := make([][]map[string]interface{}, binning.Resolution*binning.Resolution)
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
			// extract results from each bucket
			topHitsResult, ok := yBucket.TopHits(topHitsAggName)
			if !ok {
				return nil, fmt.Errorf("Top hits were not found in response for request %s", g.req.String())
			}
			// loop over raw hit results for the bin and unmarshall them into a list
			topHitsList := make([]map[string]interface{}, len(topHitsResult.Hits.Hits))
			for idx, hit := range topHitsResult.Hits.Hits {
				if err := json.Unmarshal(*hit.Source, &(topHitsList[idx])); err != nil {
					return nil, fmt.Errorf("Top hits could not be unmarshalled from response for request %s",
						g.req.String())
				}
			}
			bins[index] = topHitsList
		}
	}
	return json.Marshal(bins)
}

// GetTile returns the marshalled tile data.
func (g *MacroMicroTile) GetTile() ([]byte, error) {
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
			Aggregation(xAggName, g.getMacroAgg()).
			Do()
		if err != nil {
			return nil, err
		}
		// parse and return results
		return g.parseMacroResult(res)
	}
	// generate micro tile
	query := g.client.
		Search(g.req.Index).
		Size(0).
		Query(g.getQuery()).
		Aggregation(xAggName, g.getMicroAgg())
	// send query through equalizer
	res, err = query.Do()
	if err != nil {
		return nil, err
	}
	// parse and return results
	return g.parseMicroResult(res)
}
