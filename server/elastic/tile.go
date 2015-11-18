package elastic

import (
	"errors"
	"fmt"
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

var targetWords = []string{
	"cool",
	"awesome",
	"amazing",
	"badass",
	"killer",
	"dank",
	"superfly",
	"smooth",
	"radical",
	"wicked",
	"neato",
	"nifty",
	"primo",
	"gnarly",
	"crazy",
	"insane",
	"sick",
	"mint",
	"nice",
	"nasty",
	"classic",
	"tight",
	"rancid",
}

func buildTermsFilter( terms []string ) string {
	var filters []string
	for _, term := range terms {
		filters = append( filters, `
			"` + term + `": {
				"filter": {
					"term": {
						"text": "` + term + `"
					}
				}
			}`)
	}
	return strings.Join( filters, "," )
}

func GetJSONTile( tile *binning.TileCoord ) ( []byte, error ) {
	pixelBounds := binning.GetTilePixelBounds( tile, maxLevelSupported, tileResolution )
	xMin := strconv.FormatUint( pixelBounds.TopLeft.X, 10 )
	xMax := strconv.FormatUint( pixelBounds.BottomRight.X - 1, 10 )
	yMin := strconv.FormatUint( pixelBounds.TopLeft.Y, 10 )
	yMax := strconv.FormatUint( pixelBounds.BottomRight.Y - 1, 10 )
	request := gorequest.New()
	termsFilters := buildTermsFilter( targetWords[0:] )
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
		"aggs": {` + termsFilters + `}
	}`
	searchSize := "size=0"
	_, body, errs := request.
		Post( esHost + "/" + esIndex + "/_search?" + searchSize ).
		Send( query ).
		End()
	if errs != nil {
		return nil, errors.New( "Unable to retrieve tile data" )
	}
	//
	payload := &TopicPayload{}
	err := json.Unmarshal( []byte(body), &payload )
	if err != nil {
		fmt.Println("err")
	    return nil, err
	}
	// marshal into map then unmarshal again into string
	topicCounts := make( map[string]uint64 )
	for topic, value := range payload.Aggs {
		if value.Count > 0 {
			topicCounts[topic] = value.Count
		}
	}
	result, err := json.Marshal( topicCounts )
	if err != nil {
		fmt.Println("err")
	    return nil, err
	}
	return result, nil
}

func GetHeatmapTile( tile *binning.TileCoord ) ( []byte, error ) {
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
	                "min_doc_count": 1
	            },
	            "aggs": {
	                "y": {
	                    "histogram": {
	                        "field": "locality.pixel.y",
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
	_, body, errs := request.
		Post( esHost + "/" + esIndex + "/_search?" + searchSize + "&" + filterPath ).
		Send( query ).
		End()
	if errs != nil {
		return nil, errors.New( "Unable to retrieve tile data" )
	}
	//
	payload := &HeatmapPayload{}
	err := json.Unmarshal( []byte(body), &payload )
	if err != nil {
		fmt.Println("err")
	    return nil, err
	}
	//
	counts := make([]float64, tileResolution * tileResolution )
	cols := payload.Aggs.X.Columns
	for i := 0; i < len( cols ); i++ {
		x := cols[i].PixelX
		xBin := ( x - pixelBounds.TopLeft.X ) / xBinSize
		rows := cols[i].Y.Rows;
		for j := 0; j < len( rows ); j++ {
			y := rows[j].PixelY
			yBin := ( y - pixelBounds.TopLeft.Y ) / yBinSize
			counts[ xBin + tileResolution * yBin ] = float64( rows[j].Count )
		}
	}
	bytes := getByteArray( counts )
	return bytes, nil
}
