package elastic

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"

	"gopkg.in/olivere/elastic.v3"

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

// GetHeatmapParams returns a map of tiling parameters.
func GetHeatmapParams(tileReq *tile.Request) map[string]tile.Param {
	return map[string]tile.Param{
		"binning": NewBinningParams(tileReq),
		"topic": NewTopicParams(tileReq),
		"time": NewTimeParams(tileReq),
	}
}

// GetHeatmapTile returns a marshalled tile containing a flat array of bins.
func GetHeatmapTile(tileReq *tile.Request, params map[string]tile.Param) ([]byte, error) {
	binning, _ := params["binning"].(*BinningParams)
	time, _ := params["time"].(*TimeParams)
	topic, _ := params["topic"].(*TopicParams)
	if binning == nil {
		return nil, errors.New("No binning information has been provided")
	}
	// get client
	client, err := getClient(tileReq.Endpoint)
	if err != nil {
		return nil, err
	}
	// create x and y range queries
	boolQuery := elastic.NewBoolQuery().Must(
		binning.GetXQuery(),
		binning.GetYQuery())
	// if time params are provided, add time range query
	if time != nil {
		boolQuery.Must(time.GetTimeQuery())
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
		Aggregation(xAggName, binning.GetXAgg().
			SubAggregation(yAggName, binning.GetYAgg())).
		Do()
	if err != nil {
		return nil, err
	}
	// parse aggregations
	xAgg, ok := result.Aggregations.Histogram(xAggName)
	if !ok {
		return nil, fmt.Errorf("Histogram aggregation '%s' was not found in response", xAggName)
	}
	// allocate count buffer
	bytes := make([]byte, binning.Resolution*binning.Resolution*8)
	// fill count buffer
	for _, xBucket := range xAgg.Buckets {
		x := xBucket.Key
		xBin := (x - binning.MinX) / binning.BinSizeX
		yAgg, ok := xBucket.Histogram(yAggName)
		if !ok {
			return nil, fmt.Errorf("Histogram aggregation '%s' was not found in response", yAggName)
		}
		for _, yBucket := range yAgg.Buckets {
			y := yBucket.Key
			yBin := (y - binning.MinY) / binning.BinSizeY
			index := xBin + binning.Resolution*yBin
			// encode count
			float64ToBytes(bytes[index*8:index*8+8], float64(yBucket.DocCount))
		}
	}
	return bytes, nil
}
