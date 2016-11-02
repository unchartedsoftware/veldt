package elastic

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"encoding/json"
	"github.com/unchartedsoftware/prism/generation/elastic/agg"
	"github.com/unchartedsoftware/prism/generation/elastic/param"
	"github.com/unchartedsoftware/prism/generation/elastic/query"
	"github.com/unchartedsoftware/prism/tile"
)

const (
	topHitsAggName = "tophits"
)

// PreviewTile represents a tiling generator that produces an binned tile where
// each bin contains the top documents sorted by some field.
type PreviewTile struct {
	TileGenerator
	Binning *param.Binning
	Query   *query.Bool
	TopHits *agg.TopHits
}

// NewPreviewTile instantiates and returns a pointer to a new generator.
func NewPreviewTile(host, port string) tile.GeneratorConstructor {
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
		// optional
		topHits, err := agg.NewTopHits(tileReq.Params)
		if param.IsOptionalErr(err) {
			return nil, err
		}
		t := &PreviewTile{}
		t.Elastic = elastic
		t.Binning = binning
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
func (g *PreviewTile) GetParams() []tile.Param {
	return []tile.Param{
		g.Binning,
		g.Query,
		g.TopHits,
	}
}

func (g *PreviewTile) getQuery() elastic.Query {
	return elastic.NewBoolQuery().
		Must(g.Binning.Tiling.GetXQuery()).
		Must(g.Binning.Tiling.GetYQuery()).
		Must(g.Query.GetQuery())
}

func (g *PreviewTile) getAgg() elastic.Aggregation {
	return g.Binning.GetXAgg().
		SubAggregation(yAggName, g.Binning.GetYAgg().
			SubAggregation(topHitsAggName, g.TopHits.GetAgg()))
}

func (g *PreviewTile) parseResult(res *elastic.SearchResult) ([]byte, error) {
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
				return nil, fmt.Errorf("Top hits were not found in response for request %s",
					g.req.String())
			}
			// loop over raw hit results for the bin and unmarshall them into a list
			topHits := make([]map[string]interface{}, len(topHitsResult.Hits.Hits))
			for index, hit := range topHitsResult.Hits.Hits {
				if err := json.Unmarshal(*hit.Source, &(topHits[index])); err != nil {
					return nil, fmt.Errorf("Top hits could not be unmarshalled from response for request %s",
						g.req.String())
				}
			}
			bins[index] = topHits
		}
	}
	return json.Marshal(bins)
}

// GetTile returns the marshalled tile data.
func (g *PreviewTile) GetTile() ([]byte, error) {
	// build query
	query := g.Elastic.GetSearchService(g.client).
		Index(g.req.URI).
		Size(0).
		Query(g.getQuery()).
		Aggregation(xAggName, g.getAgg())
	// send query
	res, err := query.Do()
	if err != nil {
		return nil, err
	}
	// parse and return results
	return g.parseResult(res)
}
