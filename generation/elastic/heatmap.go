package elastic

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/util/json"
)

// HeatmapParams represents the parameters passed for a heatmap tiling tile.Request.
type HeatmapParams struct {
	X          string
	Y          string
	Extents    *binning.Bounds
	Resolution int64
}

func float64ToBytes(bytes []byte, float float64) {
	bits := math.Float64bits(float)
	binary.LittleEndian.PutUint64(bytes, bits)
}

func extractHeatmapParams(params map[string]interface{}) *HeatmapParams {
	return &HeatmapParams{
		X: json.GetStringDefault(params, "x", "pixel.x"),
		Y: json.GetStringDefault(params, "y", "pixel.y"),
		Extents: &binning.Bounds{
			TopLeft: &binning.Coord{
				X: json.GetNumberDefault(params, "minX", 0.0),
				Y: json.GetNumberDefault(params, "maxY", 0.0),
			},
			BottomRight: &binning.Coord{
				X: json.GetNumberDefault(params, "maxX", binning.MaxPixels),
				Y: json.GetNumberDefault(params, "minY", binning.MaxPixels),
			},
		},
		Resolution: int64(json.GetNumberDefault(params, "resolution", binning.MaxTileResolution)),
	}
}

// GetHeatmapHash returns a unique hash for a heatmap tile.
func GetHeatmapHash(tileReq *tile.Request) string {
	params := extractHeatmapParams(tileReq.Params)
	return fmt.Sprintf("%s:%s:%s:%d:%d:%d:%s:%s:%f:%f:%f:%f:%d",
		tileReq.Endpoint,
		tileReq.Index,
		tileReq.Type,
		tileReq.TileCoord.X,
		tileReq.TileCoord.Y,
		tileReq.TileCoord.Z,
		params.X,
		params.Y,
		params.Extents.TopLeft.X,
		params.Extents.TopLeft.Y,
		params.Extents.BottomRight.X,
		params.Extents.BottomRight.Y,
		params.Resolution,
	)
}

// GetHeatmapTile returns a marshalled tile containing a flat array of bins.
func GetHeatmapTile(tileReq *tile.Request) ([]byte, error) {
	params := extractHeatmapParams(tileReq.Params)
	bounds := binning.GetTileBounds(&tileReq.TileCoord, params.Extents)
	xBinSize := int64(bounds.BottomRight.X-bounds.TopLeft.X) / params.Resolution
	yBinSize := int64(bounds.BottomRight.Y-bounds.TopLeft.Y) / params.Resolution
	xMin := int64(bounds.TopLeft.X)
	xMax := int64(bounds.BottomRight.X - 1)
	yMin := int64(bounds.TopLeft.Y)
	yMax := int64(bounds.BottomRight.Y - 1)
	// get client
	client, err := getClient(tileReq.Endpoint)
	if err != nil {
		return nil, err
	}
	// query
	result, err := client.
		Search(tileReq.Index).
		Size(0).
		Query(elastic.NewBoolQuery().Must(
		elastic.NewRangeQuery(params.X).
			Gte(xMin).
			Lte(xMax),
		elastic.NewRangeQuery(params.Y).
			Gte(yMin).
			Lte(yMax))).
		Aggregation("x",
		elastic.NewHistogramAggregation().
			Field(params.X).
			Interval(xBinSize).
			MinDocCount(1).
			SubAggregation("y",
			elastic.NewHistogramAggregation().
				Field(params.Y).
				Interval(yBinSize).
				MinDocCount(1))).
		Do()
	if err != nil {
		return nil, err
	}
	// parse aggregations
	xAgg, ok := result.Aggregations.Histogram("x")
	if !ok {
		return nil, errors.New("Histogram aggregation 'x' was not found in response.")
	}
	// allocate count buffer
	bytes := make([]byte, params.Resolution*params.Resolution*8)
	// fill count buffer
	for _, xBucket := range xAgg.Buckets {
		x := xBucket.Key
		xBin := (x - xMin) / xBinSize
		yAgg, ok := xBucket.Histogram("y")
		if !ok {
			return nil, errors.New("Histogram aggregation 'y' was not found in response.")
		}
		for _, yBucket := range yAgg.Buckets {
			y := yBucket.Key
			yBin := (y - yMin) / yBinSize
			index := xBin + params.Resolution*yBin
			// encode count
			float64ToBytes(bytes[index*8:index*8+8], float64(yBucket.DocCount))
		}
	}
	return bytes, nil
}
