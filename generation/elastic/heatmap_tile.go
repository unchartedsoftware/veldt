package elastic

import (
	"encoding/binary"
	"fmt"
	"math"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/elastic/param"
	"github.com/unchartedsoftware/prism/generation/elastic/throttle"
	"github.com/unchartedsoftware/prism/generation/tile"
)

const (
	xAggName = "x"
	yAggName = "y"
	zAggName = "z"
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
	Binning   *param.Binning
	Topic     *param.Topic
	TimeRange *param.TimeRange
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
		topic, _ := param.NewTopic(tileReq)
		time, _ := param.NewTimeRange(tileReq)
		t := &HeatmapTile{}
		t.Binning = binning
		t.Topic = topic
		t.TimeRange = time
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
		g.Topic,
		g.TimeRange,
	}
}

// GetTile returns the marshalled tile data.
func (g *HeatmapTile) GetTile() ([]byte, error) {
	binning := g.Binning
	tiling := binning.Tiling
	timeRange := g.TimeRange
	topic := g.Topic
	tileReq := g.req
	client := g.client
	// create x and y range queries
	boolQuery := elastic.NewBoolQuery().Must(
		tiling.GetXQuery(),
		tiling.GetYQuery())
	// if time params are provided, add time range query
	if timeRange != nil {
		boolQuery.Must(timeRange.GetTimeQuery())
	}
	// if topic params are provided, add terms query
	if topic != nil {
		boolQuery.Must(topic.GetTopicQuery())
	}
	// create x aggregation
	xAgg := binning.GetXAgg()
	// create y aggregation, add it as a sub-agg to xAgg
	yAgg := binning.GetYAgg()
	xAgg.SubAggregation(yAggName, yAgg)
	// if there is a z field to sum, add sum agg to yAgg
	if binning.Z != "" {
		zAgg := binning.GetZAgg()
		yAgg.SubAggregation(zAggName, zAgg)
	}
	// build query
	query := client.
		Search(tileReq.Index).
		Size(0).
		Query(boolQuery).
		Aggregation(xAggName, xAgg)
	// send query through equalizer
	result, err := throttle.Send(query)
	if err != nil {
		return nil, err
	}
	// parse aggregations
	xAggRes, ok := result.Aggregations.Histogram(xAggName)
	if !ok {
		return nil, fmt.Errorf("Histogram aggregation '%s' was not found in response for request %s",
			xAggName,
			tileReq.String())
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
				tileReq.String())
		}
		for _, yBucket := range yAggRes.Buckets {
			y := yBucket.Key
			yBin := binning.GetYBin(y)
			index := xBin + binning.Resolution*yBin
			if binning.Z != "" {
				// extract metric
				zAggRes, ok := binning.GetZAggValue(zAggName, yBucket)
				if !ok {
					return nil, fmt.Errorf("'%s' aggregation '%s' was not found in response for request %s",
						binning.Metric,
						zAggName,
						tileReq.String())
				}
				// encode metric
				if (zAggRes.Value != nil) {
					bins[index] += *zAggRes.Value
				}
			} else {
				// encode count
				bins[index] += float64(yBucket.DocCount)
			}
		}
	}
	return float64ToByteSlice(bins), nil
}
