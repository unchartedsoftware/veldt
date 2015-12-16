package elastic

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/generation/tile"
	jsonutil "github.com/unchartedsoftware/prism/util/json"
)

// TopicCountParams represents the parameters passed for a topic count tiling tile.Request.
type TopicCountParams struct {
	X       string
	Y       string
	Extents *binning.Bounds
	Text	string
	Topics  []string
}

func extractTopicCountParams(params map[string]interface{}) *TopicCountParams {
	topics := jsonutil.GetStringDefault(params, "topics", "")
	return &TopicCountParams{
		X: jsonutil.GetStringDefault(params, "x", "pixel.x"),
		Y: jsonutil.GetStringDefault(params, "y", "pixel.y"),
		Extents: &binning.Bounds{
			TopLeft: &binning.Coord{
				X: jsonutil.GetNumberDefault(params, "minX", 0.0),
				Y: jsonutil.GetNumberDefault(params, "maxY", 0.0),
			},
			BottomRight: &binning.Coord{
				X: jsonutil.GetNumberDefault(params, "maxX", binning.MaxPixels),
				Y: jsonutil.GetNumberDefault(params, "minY", binning.MaxPixels),
			},
		},
		Text: jsonutil.GetStringDefault(params, "text", "text"),
		Topics: strings.Split(topics, ","),
	}
}

// GetTopicCountHash returns a unique hash for a topic count tile.
func GetTopicCountHash(tileReq *tile.Request) string {
	params := extractTopicCountParams(tileReq.Params)
	return fmt.Sprintf("%s:%s:%s:%d:%d:%d:%s:%s:%f:%f:%f:%f:%s",
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
		strings.Join(params.Topics, ":"),
	)
}

// GetTopicCountTile returns a marshalled tile containing topics and counts.
func GetTopicCountTile(tileReq *tile.Request) ([]byte, error) {
	params := extractTopicCountParams(tileReq.Params)
	bounds := binning.GetTileBounds(&tileReq.TileCoord, params.Extents)
	xMin := int64(bounds.TopLeft.X)
	xMax := int64(bounds.BottomRight.X - 1)
	yMin := int64(bounds.TopLeft.Y)
	yMax := int64(bounds.BottomRight.Y - 1)
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
		elastic.NewRangeQuery(params.X).
			Gte(xMin).
			Lte(xMax),
		elastic.NewRangeQuery(params.Y).
			Gte(yMin).
			Lte(yMax)))
	// add all filter aggregations
	for _, topic := range params.Topics {
		query.Aggregation(topic,
			elastic.NewFilterAggregation().
				Filter(elastic.NewTermQuery(params.Text, topic)))
	}
	// send query
	result, err := query.Do()
	if err != nil {
		return nil, err
	}
	// build map of topics and counts
	topicCounts := make(map[string]int64)
	for _, topic := range params.Topics {
		filter, ok := result.Aggregations.Filter(topic)
		if ok && filter.DocCount > 0 {
			topicCounts[topic] = filter.DocCount
		}
	}
	// marshal results map
	return json.Marshal(topicCounts)
}
