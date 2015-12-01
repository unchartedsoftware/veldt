package elastic

import (
	"encoding/json"
	"strconv"

	"gopkg.in/olivere/elastic.v3"

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
func GetTopicCountTile(endpoint string, index string, tile *binning.TileCoord) ([]byte, error) {
	pixelBounds := binning.GetTilePixelBounds(tile)
	xMin := strconv.FormatUint(pixelBounds.TopLeft.X, 10)
	xMax := strconv.FormatUint(pixelBounds.BottomRight.X-1, 10)
	yMin := strconv.FormatUint(pixelBounds.TopLeft.Y, 10)
	yMax := strconv.FormatUint(pixelBounds.BottomRight.Y-1, 10)
	// query
	query := client.
		Search(index).
		Size(0).
		Query(elastic.NewBoolQuery().Must(
		elastic.NewRangeQuery("pixel.x").
			Gte(xMin).
			Lte(xMax),
		elastic.NewRangeQuery("pixel.y").
			Gte(yMin).
			Lte(yMax),
	))
	// add all filter aggregations
	for _, term := range terms {
		query.Aggregation(term,
			elastic.NewFilterAggregation().
				Filter(elastic.NewTermQuery("text", term)))
	}

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
