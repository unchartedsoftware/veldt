package elastic

import (
	"encoding/json"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/tiling"
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

var terms = []string{
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

// GetTopicCountTile returns a marshalled tile containing topics and their counts.
func GetTopicCountTile(tileReq *tiling.TileRequest) ([]byte, error) {
	pixelBounds := binning.GetTilePixelBounds(&tileReq.TileCoord)
	xMin := int64(pixelBounds.TopLeft.X)
	xMax := int64(pixelBounds.BottomRight.X - 1)
	yMin := int64(pixelBounds.TopLeft.Y)
	yMax := int64(pixelBounds.BottomRight.Y - 1)
	// get client
	client, err := getClient(tileReq.Endpoint)
	if err != nil {
		return nil, err
	}
	// build query
	query := client.
		Search(tileReq.Index).
		Size(0).
		Query(elastic.NewBoolQuery().Must(
		elastic.NewRangeQuery("pixel.x").
			Gte(xMin).
			Lte(xMax),
		elastic.NewRangeQuery("pixel.y").
			Gte(yMin).
			Lte(yMax)))
	// add all filter aggregations
	for _, term := range terms {
		query.Aggregation(term,
			elastic.NewFilterAggregation().
				Filter(elastic.NewTermQuery("text", term)))
	}
	// send query
	result, err := query.Do()
	if err != nil {
		return nil, err
	}
	// build map of topics and counts
	topicCounts := make(map[string]int64)
	for _, term := range terms {
		filter, ok := result.Aggregations.Filter(term)
		if ok && filter.DocCount > 0 {
			topicCounts[term] = filter.DocCount
		}
	}
	// marshal results map
	return json.Marshal(topicCounts)
}
