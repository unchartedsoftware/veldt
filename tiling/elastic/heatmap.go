package elastic

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"math"
	"strconv"

	"github.com/parnurzeal/gorequest"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/util/log"
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

// Row represents a row value in the y axis aggregations.
type Row struct {
	Count  uint64 `json:"doc_count"`
	PixelY uint64 `json:"key"`
}

// YAgg represents the y value aggregations inside a column.
type YAgg struct {
	Rows []Row `json:"buckets"`
}

// Column represents a column of values in an x aggregation.
type Column struct {
	Y      YAgg   `json:"y"`
	PixelX uint64 `json:"key"`
}

// XAgg represents the set of bucketed columns.
type XAgg struct {
	Columns []Column `json:"buckets"`
}

// Aggregation represents the histogram aggregations.
type Aggregation struct {
	X XAgg `json:"x"`
}

// HeatmapPayload represents the aggregations payload of the elasticsearch response.
type HeatmapPayload struct {
	Aggs Aggregation `json:"aggregations"`
}

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
func GetHeatmapTile(tile *binning.TileCoord) ([]byte, error) {
	tileResolution := binning.MaxTileResolution // TEMP
	pixelBounds := binning.GetTilePixelBounds(tile)
	xBinSize := (pixelBounds.BottomRight.X - pixelBounds.TopLeft.X) / tileResolution
	yBinSize := (pixelBounds.BottomRight.Y - pixelBounds.TopLeft.Y) / tileResolution
	xMin := strconv.FormatUint(pixelBounds.TopLeft.X, 10)
	xMax := strconv.FormatUint(pixelBounds.BottomRight.X-1, 10)
	yMin := strconv.FormatUint(pixelBounds.TopLeft.Y, 10)
	yMax := strconv.FormatUint(pixelBounds.BottomRight.Y-1, 10)
	xInterval := strconv.FormatUint(xBinSize, 10)
	yInterval := strconv.FormatUint(yBinSize, 10)
	query := `{
		"query": {
			"bool" : {
		        "must" : [
					{
			            "range": {
							"pixel.x": {
								"gte":` + xMin + `,
								"lte":` + xMax + `
							}
						}
					},
					{
			            "range": {
							"pixel.y": {
								"gte":` + yMin + `,
								"lte":` + yMax + `
							}
						}
					}
				]
			}
		},
		"aggs": {
	        "x": {
	            "histogram": {
	                "field": "pixel.x",
	                "interval":` + xInterval + `,
	                "min_doc_count": 1
	            },
	            "aggs": {
	                "y": {
	                    "histogram": {
	                        "field": "pixel.y",
	                        "interval":` + yInterval + `,
	                        "min_doc_count": 1
	                    }
	                }
	            }
	        }
	    }
	}`
	searchSize := "size=0"
	filterPath := "filter_path=aggregations.x.buckets.key,aggregations.x.buckets.y.buckets.key,aggregations.x.buckets.y.buckets.doc_count"
	_, body, errs := gorequest.
		New().
		Post(esHost + "/" + esIndex + "/_search?" + searchSize + "&" + filterPath).
		Send(query).
		End()
	if errs != nil {
		return nil, errors.New("Unable to retrieve tile data")
	}
	// unmarshal payload
	payload := &HeatmapPayload{}
	err := json.Unmarshal([]byte(body), &payload)
	if err != nil {
		log.Warn(err)
		return nil, err
	}
	// build array of counts
	counts := make([]float64, tileResolution*tileResolution)
	cols := payload.Aggs.X.Columns
	for i := 0; i < len(cols); i++ {
		x := cols[i].PixelX
		xBin := (x - pixelBounds.TopLeft.X) / xBinSize
		rows := cols[i].Y.Rows
		for j := 0; j < len(rows); j++ {
			y := rows[j].PixelY
			yBin := (y - pixelBounds.TopLeft.Y) / yBinSize
			counts[xBin+tileResolution*yBin] = float64(rows[j].Count)
		}
	}
	bytes := getByteArray(counts)
	return bytes, nil
}
