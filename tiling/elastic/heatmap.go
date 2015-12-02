package elastic

import (
	"encoding/binary"
	"errors"
	"math"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/binning"
)

// "aggregations": {
//     "x": {
//         "buckets": [
//             {
// 			       "key": 1261961216,
// 	        	   "y": {
// 		               "buckets": [
// 		            	   {
// 							   "key": 1615331328,"
// 			                   "doc_count": 10
// 			               },
// 						   ...
// 					   ]
// 				   }
// 			   },
// 			   ...
// 		   ]
// 	   }
// }

func float64ToBytes(bytes []byte, float float64) {
	bits := math.Float64bits(float)
	binary.LittleEndian.PutUint64(bytes, bits)
}

func getByteArray(data []float64) []byte {
	buf := make([]byte, len(data)*8)
	for i := 0; i < len(data); i++ {
		float64ToBytes(buf[i*8:i*8+8], data[i])
	}
	return buf[0:]
}

// GetHeatmapTile returns a marshalled tile containing a flat array of bins.
func GetHeatmapTile(endpoint string, index string, tile *binning.TileCoord) ([]byte, error) {
	pixelBounds := binning.GetTilePixelBounds(tile)
	tileResolution := int64(binning.MaxTileResolution)
	xBinSize := int64(pixelBounds.BottomRight.X-pixelBounds.TopLeft.X) / tileResolution
	yBinSize := int64(pixelBounds.BottomRight.Y-pixelBounds.TopLeft.Y) / tileResolution
	xMin := int64(pixelBounds.TopLeft.X)
	xMax := int64(pixelBounds.BottomRight.X - 1)
	yMin := int64(pixelBounds.TopLeft.Y)
	yMax := int64(pixelBounds.BottomRight.Y - 1)
	// get client
	client, err := getClient(endpoint)
	if err != nil {
		return nil, err
	}
	// query
	result, err := client.
		Search(index).
		Size(0).
		Query(elastic.NewBoolQuery().Must(
		elastic.NewRangeQuery("pixel.x").
			Gte(xMin).
			Lte(xMax),
		elastic.NewRangeQuery("pixel.y").
			Gte(yMin).
			Lte(yMax),
	)).
		Aggregation("x",
		elastic.NewHistogramAggregation().
			Field("pixel.x").
			Interval(xBinSize).
			MinDocCount(1).
			SubAggregation("y",
			elastic.NewHistogramAggregation().
				Field("pixel.y").
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
	bytes := make([]byte, tileResolution*tileResolution*8)
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
			index := xBin + tileResolution*yBin
			// encode count
			float64ToBytes(bytes[index*8:index*8+8], float64(yBucket.DocCount))
		}
	}
	return bytes, nil
}
