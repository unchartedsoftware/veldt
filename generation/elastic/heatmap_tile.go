package elastic

import (
	"encoding/binary"
	"fmt"
	"math"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/elastic/agg"
	"github.com/unchartedsoftware/prism/generation/elastic/param"
	"github.com/unchartedsoftware/prism/generation/elastic/query"
	"github.com/unchartedsoftware/prism/generation/tile"
)

const (
	xAggName      = "x"
	yAggName      = "y"
	metricAggName = "metric"
)

func float64ToBytes(bytes []byte, float float64) {
	bits := math.Float64bits(float)
	binary.LittleEndian.PutUint64(bytes, bits)
}

func float64ToByteSlice(arr []float64) []byte {
	bits := make([]byte, len(arr)*8)
	for i, a := range arr {
		float64ToBytes(bits[i*8:i*8+8], a)
	}
	return bits[0:]
}

// HeatmapTile represents a tiling generator that produces heatmaps.
type HeatmapTile struct {
	TileGenerator
	Binning *param.Binning
	Query   *query.Bool
	Metric  *agg.Metric
}

// NewHeatmapTile instantiates and returns a pointer to a new generator.
func NewHeatmapTile(host, port string) tile.GeneratorConstructor {
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
		// optional
		metric, err := agg.NewMetric(tileReq.Params)
		if param.IsOptionalErr(err) {
			return nil, err
		}
		t := &HeatmapTile{}
		t.Binning = binning
		t.Query = query
		t.Metric = metric
		t.req = tileReq
		t.host = host
		t.port = port
		t.client = client
		return t, nil
	}
}

// GetParams returns a slice of tiling parameters.
func (g *HeatmapTile) GetParams() []tile.Param {
	return []tile.Param{
		g.Binning,
		g.Query,
		g.Metric,
	}
}

func (g *HeatmapTile) getQuery() elastic.Query {
	return elastic.NewBoolQuery().
		Must(g.Binning.Tiling.GetXQuery()).
		Must(g.Binning.Tiling.GetYQuery()).
		Must(g.Query.GetQuery())
}

func (g *HeatmapTile) getAgg() elastic.Aggregation {
	// create x aggregation
	xAgg := g.Binning.GetXAgg()
	// create y aggregation, add it as a sub-agg to xAgg
	yAgg := g.Binning.GetYAgg()
	xAgg.SubAggregation(yAggName, yAgg)
	// if there is a z field to sum, add sum agg to yAgg
	if g.Metric != nil {
		yAgg.SubAggregation(metricAggName, g.Metric.GetAgg())
	}
	return xAgg
}

func (g *HeatmapTile) parseResult(res *elastic.SearchResult) ([]byte, error) {
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
			if g.Metric != nil {
				// extract metric
				value, ok := g.Metric.GetAggValue(metricAggName, yBucket)
				if !ok {
					return nil, fmt.Errorf("'%s' aggregation '%s' was not found in response for request %s",
						g.Metric.Type,
						metricAggName,
						g.req.String())
				}
				// encode metric
				bins[index] += value
			} else {
				// encode count
				bins[index] += float64(yBucket.DocCount)
			}
		}
	}
	return float64ToByteSlice(bins), nil
}

// GetTile returns the marshalled tile data.
func (g *HeatmapTile) GetTile() ([]byte, error) {
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
