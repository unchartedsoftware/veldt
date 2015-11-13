package elastic

import (
	"errors"
	//"fmt"
	"math"
	"strconv"
	"strings"
	"encoding/binary"
	"encoding/json"

	"github.com/parnurzeal/gorequest"

	"github.com/unchartedsoftware/prism/server/binning"
)

const esHost = "http://memex3:9200"
const esIndex = "nyc_twitter_july" //"isil_twitter_dec2may"

const tileResolution = 256

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


// GetTile takes a tile coord and returns data
func GetGeoTile( tile *binning.TileCoord ) ( []byte, error ) {
	bounds := binning.GetTileGeoBounds( tile )
	request := gorequest.New()
	query := `{
	    "_source": "locality.location",
	    "query":{
	        "filtered": {
	            "filter": {
	                "bool": {
	                    "must": [
	                        {
	                            "exists" : {
	                                "field" : "locality.location"
	                            }
	                        },
	                        {
	                            "geo_bounding_box": {
	                                "locality.location": {
	                                    "top_left" : {
	                                        "lon" : ` + strconv.FormatFloat( bounds.BottomLeft.Lon, 'f', 6, 64 ) + `,
	                                        "lat" : ` + strconv.FormatFloat( bounds.TopRight.Lat, 'f', 6, 64 ) + `
	                                    },
	                                    "bottom_right" : {
	                                        "lon" : ` + strconv.FormatFloat( bounds.TopRight.Lon, 'f', 6, 64 ) + `,
	                                        "lat" : ` + strconv.FormatFloat( bounds.BottomLeft.Lat, 'f', 6, 64 ) + `
	                                    }
	                                }
	                            }
	                        }
	                    ]
	                }
	            }
	        }
	    }
	}`
	_, body, errs := request.
		Post( esHost + "/" + esIndex + "/_search?size=10000" ).
		Send( query ).
		End()
	if errs != nil {
		return nil, errors.New( "Unable to retrieve tile data" )
	}
	//
	payload := &JsonPayload{}
	err := json.Unmarshal( []byte(body), &payload )
	if err != nil {
	    return nil, err
	}

	counts := make([]float64, tileResolution * tileResolution )
	max := 0.0
	for i := 0; i < len( payload.Hits.Bins ); i++ {
		lonLat, err := parseLonLat( &payload.Hits.Bins[i].Source.Locality.Location )
		if err == nil {
			binIndex := binning.LonLatToFlatBin( lonLat, tile.Z, tileResolution )
			counts[ binIndex ]++
			if counts[ binIndex ] > max {
				max = counts[ binIndex ]
			}
		}
	}
	// for i := 0; i < len( counts ); i++ {
	// 	if counts[i] > 0 {
	// 		fmt.Printf("%f\n", counts[i])
	// 	}
	// }
	bytes := getByteArray( counts )
	return bytes, nil
}
