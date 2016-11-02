package citus

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/jackc/pgx"

	"github.com/unchartedsoftware/prism/generation/citus/param"
	"github.com/unchartedsoftware/prism/generation/citus/query"

	"github.com/unchartedsoftware/prism/tile"
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
	Query   *query.Query
	//Metric  *agg.Metric
}

// NewHeatmapTile instantiates and returns a pointer to a new generator.
func NewHeatmapTile(host, port string) tile.GeneratorConstructor {
	return func(tileReq *tile.Request) (tile.Generator, error) {
		client, err := NewClient(host, port)
		if err != nil {
			return nil, err
		}
		citus, err := param.NewCitus(tileReq)
		if err != nil {
			return nil, err
		}
		binning, err := param.NewBinning(tileReq)
		if err != nil {
			return nil, err
		}
		query, err := query.NewQuery()
		if err != nil {
			return nil, err
		}
		// optional
		//metric, err := agg.NewMetric(tileReq.Params)
		//if param.IsOptionalErr(err) {
		//	return nil, err
		//}
		t := &HeatmapTile{}
		t.Citus = citus
		t.Binning = binning
		t.Query = query
		//t.Metric = metric
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
		//g.Metric,
	}
}

func (g *HeatmapTile) getQuery() *query.Query {
	g.Binning.Tiling.AddXQuery(g.Query)
	g.Binning.Tiling.AddYQuery(g.Query)
	return g.Query
}

func (g *HeatmapTile) getAgg() *query.Query {
	// create x aggregation
	g.Binning.AddXAgg(g.Query)
	// create y aggregation
	g.Binning.AddYAgg(g.Query)
	// if there is a z field to sum, add sum agg to yAgg
	g.Query.AddField(fmt.Sprintf("CAST(COUNT(*) AS FLOAT) as %s", metricAggName))
	//if g.Metric != nil {
	//	yAgg.SubAggregation(metricAggName, g.Metric.GetAgg())
	//}
	return g.Query
}

func (g *HeatmapTile) parseResult(rows *pgx.Rows) ([]byte, error) {
	binning := g.Binning

	// allocate bins buffer
	bins := make([]float64, binning.Resolution*binning.Resolution)
	// fill bins buffer
	for rows.Next() {
		var x, y int64
		var value float64

		err := rows.Scan(&x, &y, &value)
		if err != nil {
			return nil, fmt.Errorf("Error parsing histogram aggregation: %s %v",
				g.req.String(), err)
		}

		xBin := binning.GetXBin(x)
		yBin := binning.GetYBin(y)

		index := xBin + binning.Resolution*yBin
		bins[index] += value
	}

	return float64ToByteSlice(bins), nil
}

// GetTile returns the marshalled tile data.
func (g *HeatmapTile) GetTile() ([]byte, error) {
	heatReq := g.req
	// send query
	query := g.getQuery()
	query = g.getAgg()
	query.AddTable(heatReq.URI)
	rows, err := g.client.Query(query.GetQuery(false), query.QueryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// parse and return results
	return g.parseResult(rows)
}
