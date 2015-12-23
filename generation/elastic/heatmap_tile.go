package elastic

import (
	"encoding/binary"
	"fmt"
	"math"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/elastic/param"
	"github.com/unchartedsoftware/prism/generation/tile"
)

const (
	xAggName = "x"
	yAggName = "y"
)

func float64ToBytes(bytes []byte, float float64) {
	bits := math.Float64bits(float)
	binary.LittleEndian.PutUint64(bytes, bits)
}

// HeatmapTile represents a tiling generator that produces heatmaps.
type HeatmapTile struct {
	Binning   *param.Binning
	Topic     *param.Topic
	TimeRange *param.TimeRange
}

// NewHeatmapTile instantiates and returns a pointer to a new generator.
func NewHeatmapTile(tileReq *tile.Request) (tile.Generator, error) {
	binning, err := param.NewBinning(tileReq)
	if err != nil {
		return nil, err
	}
	topic, _ := param.NewTopic(tileReq)
	time, _ := param.NewTimeRange(tileReq)
	return &HeatmapTile{
		Binning:   binning,
		Topic:     topic,
		TimeRange: time,
	}, nil
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
func (g *HeatmapTile) GetTile(tileReq *tile.Request) ([]byte, error) {
	binning := g.Binning
	tiling := binning.Tiling
	timeRange := g.TimeRange
	topic := g.Topic
	// get client
	client, err := getClient(tileReq.Endpoint)
	if err != nil {
		return nil, err
	}
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
	// query
	result, err := client.
		Search(tileReq.Index).
		Size(0).
		Query(boolQuery).
		Aggregation(xAggName,
		binning.GetXAgg().
			SubAggregation(yAggName, binning.GetYAgg())).
		Do()
	if err != nil {
		return nil, err
	}
	// parse aggregations
	xAgg, ok := result.Aggregations.Histogram(xAggName)
	if !ok {
		return nil, fmt.Errorf("Histogram aggregation '%s' was not found in response for request %s", xAggName, tileReq.String())
	}
	// allocate count buffer
	bytes := make([]byte, binning.Resolution*binning.Resolution*8)
	// fill count buffer
	for _, xBucket := range xAgg.Buckets {
		x := xBucket.Key
		xBin := (x - tiling.Left) / binning.BinSizeX
		yAgg, ok := xBucket.Histogram(yAggName)
		if !ok {
			return nil, fmt.Errorf("Histogram aggregation '%s' was not found in response for request %s", yAggName, tileReq.String())
		}
		for _, yBucket := range yAgg.Buckets {
			y := yBucket.Key
			yBin := (y - tiling.Top) / binning.BinSizeY
			index := xBin + binning.Resolution*yBin
			// encode count
			float64ToBytes(bytes[index*8:index*8+8], float64(yBucket.DocCount))
		}
	}
	return bytes, nil
}