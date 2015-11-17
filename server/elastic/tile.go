package elastic

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"encoding/binary"
	"encoding/json"

	"github.com/parnurzeal/gorequest"

	"github.com/unchartedsoftware/prism/binning"
)

const esHost = "http://10.64.16.120:9200"
const esIndex = "nyc_twitter"
const tileResolution = 256
const maxLevelSupported = 24

func float64ToBytes( bytes []byte, float float64 ) {
    bits := math.Float64bits(float)
    binary.LittleEndian.PutUint64(bytes, bits)
}

func getByteArray( data []float64 ) []byte {
	buf := make([]byte, tileResolution * tileResolution * 8 )
	for i := 0; i < len( data ); i++ {
		float64ToBytes( buf[ i*8 : i*8+8 ], data[i] )
	}
	return buf
}

func parseLonLat( location *string ) ( *binning.LonLat, error ) {
	latLonStr := strings.Split( *location, "," )
	lat, elat := strconv.ParseFloat( latLonStr[0], 64 )
	lon, elon := strconv.ParseFloat( latLonStr[1], 64 )
	if elat == nil || elon == nil {
		return &binning.LonLat{
			Lon: lon,
			Lat: lat,
		}, nil
	}
	return nil, errors.New( "Unable to parse tile coordinate from URL" )
}

func GetTile( tile *binning.TileCoord ) ( []byte, error ) {
	pixelBounds := binning.GetTilePixelBounds( tile, maxLevelSupported, tileResolution )
	xBinSize := ( pixelBounds.BottomRight.X - pixelBounds.TopLeft.X ) / tileResolution
	yBinSize := ( pixelBounds.BottomRight.Y - pixelBounds.TopLeft.Y ) / tileResolution
	xMin := strconv.FormatUint( pixelBounds.TopLeft.X, 10 )
	xMax := strconv.FormatUint( pixelBounds.BottomRight.X - 1, 10 )
	yMin := strconv.FormatUint( pixelBounds.TopLeft.Y, 10 )
	yMax := strconv.FormatUint( pixelBounds.BottomRight.Y - 1, 10 )
	xInterval := strconv.FormatUint( xBinSize, 10 )
	yInterval := strconv.FormatUint( yBinSize, 10 )
	request := gorequest.New()
	query := `{
		"query": {
			"bool" : {
		        "must" : [
					{
			            "range": {
							"locality.pixel.x": {
								"gte":` + xMin + `,
								"lte":` + xMax + `
							}
						}
					},
					{
			            "range": {
							"locality.pixel.y": {
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
	                "field": "locality.pixel.x",
	                "interval":` + xInterval + `,
	                "min_doc_count": 0,
					"extended_bounds": {
						"min":` + xMin + `,
						"max":` + xMax + `
                    }
	            },
	            "aggs": {
	                "y": {
	                    "histogram": {
	                        "field": "locality.pixel.y",
	                        "interval":` + yInterval + `,
	                        "min_doc_count": 0,
							"extended_bounds": {
								"min":` + yMin + `,
								"max":` + yMax + `
                            }
	                    }
	                }
	            }
	        }
	    }
	}`
	//fmt.Println( query )
	searchSize := "size=0"
	filterPath := "filter_path=aggregations.x.buckets.y.buckets.doc_count"
	_, body, errs := request.
		Post( esHost + "/" + esIndex + "/_search?" + searchSize + "&" + filterPath ).
		Send( query ).
		End()
	if errs != nil {
		return nil, errors.New( "Unable to retrieve tile data" )
	}
	//
	payload := &Payload{}
	err := json.Unmarshal( []byte(body), &payload )
	if err != nil {
	    return nil, err
	}
	//
	counts := make([]float64, tileResolution * tileResolution )
	cols := payload.Aggs.X.Columns
	for i := 0; i < tileResolution; i++ {
		rows := cols[i].Y.Rows;
		for j := 0; j < tileResolution; j++ {
			counts[ i + tileResolution * j ] = float64( rows[j].Count )
		}
	}
	/*
		ISSUE: EXTENDED_BOUNDS DOESN'T WORK IF THE BOUNDS IS OUTSIDE INDEX MIN / MAX,
			CURRENTLY STATIC SIZING THE NUM BINS DOESN'T WORK CORRECTLY. FIX THIS.
	*/
	bytes := getByteArray( counts )
	return bytes, nil
}
