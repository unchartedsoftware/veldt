package remote

import (
	"encoding/json"
	"fmt"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
	jsonUtil "github.com/unchartedsoftware/veldt/util/json"
)

type TopicTile struct {
	inclusionTerms []string
	exclusionTerms []string
	requestId      string
	tileCount      int
	exclusiveness  int
	clusterCount   int
	wordCount      int
	x              uint32
	y              uint32
	z              uint32
}

func NewTopicTile() veldt.TileCtor {
	return func() (veldt.Tile, error) {
		return &TopicTile{}, nil
	}
}

func (t *TopicTile) Parse(params map[string]interface{}) error {
	// get inclusion terms
	include, ok := jsonUtil.GetStringArray(params, "include")
	if !ok {
		return fmt.Errorf("`include` parameter missing from topic tile")
	}
	// get exclusion terms
	exclude, ok := jsonUtil.GetStringArray(params, "exclude")
	if !ok {
		return fmt.Errorf("`exclude` parameter missing from topic tile")
	}
	// get request id
	requestId, ok := jsonUtil.GetString(params, "requestId")
	if !ok {
		return fmt.Errorf("`requestId` parameter missing from topic tile")
	}
	// get tile count
	tileCount, ok := jsonUtil.GetNumber(params, "tileCount")
	if !ok {
		return fmt.Errorf("`tileCount` parameter missing from topic tile")
	}
	// get exclusiveness
	exclusiveness, ok := jsonUtil.GetNumber(params, "exclusiveness")
	if !ok {
		return fmt.Errorf("`tileCount` parameter missing from topic tile")
	}
	// get topic word count
	wordCount, ok := jsonUtil.GetNumber(params, "topicWordCount")
	if !ok {
		return fmt.Errorf("`tileCount` parameter missing from topic tile")
	}
	// get topic cluster count
	clusterCount, ok := jsonUtil.GetNumber(params, "topicClusterCount")
	if !ok {
		return fmt.Errorf("`tileCount` parameter missing from topic tile")
	}

	t.inclusionTerms = include
	t.exclusionTerms = exclude
	t.exclusiveness = int(exclusiveness)
	t.wordCount = int(wordCount)
	t.clusterCount = int(clusterCount)
	t.requestId = requestId
	t.tileCount = int(tileCount)

	return nil
}

func (t *TopicTile) Create(uri string, coord *binning.TileCoord, query veldt.Query) ([]byte, error) {
	// Setup the tile coordinate information.
	t.x = coord.X
	t.y = coord.Y
	t.z = coord.Z

	// Send the request to the batching client and wait for the response.
	client := getServiceClient(t)
	resChan, err := client.AddRequest(t)
	if err != nil {
		return nil, err
	}

	res := <-resChan

	// Either an error, or the response from the remote service.
	err, ok := res.(error)
	if ok {
		return nil, err
	}

	// Encode the results. Extract all the topics and use the score for weighing.
	// Result is a string containing the JSON. Need to get to the topics.
	var tmpTopicParsed interface{}
	err = json.Unmarshal([]byte(res.(string)), &tmpTopicParsed)
	if err != nil {
		return nil, err
	}

	topicsParsed, ok := tmpTopicParsed.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Unexpected response format from topic modelling service: incorrect structure in %v", res)
	}

	counts := make(map[string]uint32)
	topics, ok := jsonUtil.GetArray(topicsParsed, "topic")
	for _, topic := range topics {
		topicMap, ok := topic.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Unexpected response format from topic modelling service: incorrect topic structure in %v", res)
		}
		words, ok := jsonUtil.GetArray(topicMap, "words")
		if !ok {
			return nil, fmt.Errorf("Unexpected response format from topic modelling service: cannot find 'words' in %v", res)
		}

		for _, wordEntry := range words {
			wordParsed, ok := wordEntry.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("Unexpected response format from topic modelling service: incorrect word structure in %v", res)
			}
			word, ok := jsonUtil.GetString(wordParsed, "word")
			if !ok {
				return nil, fmt.Errorf("Unexpected response format from topic modelling service: cannot find 'word' in %v", res)
			}
			weight, ok := jsonUtil.GetNumber(wordParsed, "score")
			if !ok {
				return nil, fmt.Errorf("Unexpected response format from topic modelling service: cannot find 'score' in %v", res)
			}

			counts[word] = counts[word] + uint32(weight)
		}
	}

	// marshal results
	return json.Marshal(counts)
}

func (t *TopicTile) GetTileId() string {
	return fmt.Sprintf("%v/%v/%v", t.z, t.x, t.y)
}
