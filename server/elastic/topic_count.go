package elastic

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"encoding/json"

	"github.com/parnurzeal/gorequest"

	"github.com/unchartedsoftware/prism/binning"
)

// "aggregations": {
//     "term0": {
//     		"doc_count": 234
// 	    },
// 	    "term1": {
//     		"doc_count": 234
// 	    },
// 	    "term2": {
//     		"doc_count": 234
// 	    }
//		...
// }

// Topic represents the count for a specific topic.
type Topic struct {
	Count uint64 `json:"doc_count"`
}

// TopicPayload represents the aggregations payload of the elasticsearch response.
type TopicPayload struct {
	Aggs map[string]Topic `json:"aggregations"`
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

// GetTopicCountTile returns a marshalled tile containing topics and their counts.
func GetTopicCountTile( tile *binning.TileCoord ) ( []byte, error ) {
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
