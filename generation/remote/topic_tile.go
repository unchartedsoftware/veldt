package remote

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
	jsonUtil "github.com/unchartedsoftware/veldt/util/json"
)

type TopicTile struct {
	APITile
	inclusionTerms []string
	exclusionTerms []string
	exclusiveness  int
	clusterCount   int
	wordCount      int
	timeFrom       int64
	timeTo         int64
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
	// get time from
	timeFrom, ok := jsonUtil.GetNumber(params, "timeFrom")
	if !ok {
		return fmt.Errorf("`timeFrom` parameter missing from topic tile")
	}
	// get time time
	timeTo, ok := jsonUtil.GetNumber(params, "timeTo")
	if !ok {
		return fmt.Errorf("`timeTo` parameter missing from topic tile")
	}

	t.inclusionTerms = include
	t.exclusionTerms = exclude
	t.exclusiveness = int(exclusiveness)
	t.wordCount = int(wordCount)
	t.clusterCount = int(clusterCount)
	t.timeFrom = int64(timeFrom)
	t.timeTo = int64(timeTo)

	// Use hash of all parameters to identify request.
	requestId := fmt.Sprintf("%v::%v::%v::%v::%v::%v::%v",
		wordCount, clusterCount, timeFrom, timeTo, exclusiveness,
		strings.Join(include, ","), strings.Join(exclude, ","))

	t.requestId = requestId
	t.tileType = "topic"

	return nil
}

func (t *TopicTile) GetApiUrl() string {
	// TODO: Have the URL configurable!
	return "http://163.152.20.64:5000/GET_TOPICS/test"
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

	counts := make(map[uint32]map[string]interface{})
	topics, ok := jsonUtil.GetArray(topicsParsed, "topic")
	topicGroup := uint32(0)
	for _, topic := range topics {
		counts[topicGroup] = make(map[string]interface{})
		topicWords := make(map[string]uint32)
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

			// Set the weight & topic group of the word.
			topicWords[word] = uint32(weight)
		}
		exclusiveness, ok := jsonUtil.GetNumber(topicMap, "exclusiveness")
		if !ok {
			return nil, fmt.Errorf("Unexpected response format from topic modelling service: cannot find 'exclusiveness' in %v", res)
		}
		counts[topicGroup]["exclusiveness"] = exclusiveness
		counts[topicGroup]["words"] = topicWords
		topicGroup = topicGroup + 1
	}

	// marshal results
	return json.Marshal(counts)
}

// Create the request to the remote service.
func (t *TopicTile) GetRequestParameters() map[string]interface{} {
	parameters := make(map[string]interface{})
	parameters["include_words"] = t.inclusionTerms
	parameters["exclude_words"] = t.exclusionTerms

	// Add simple parameters.
	parameters["exclusiveness"] = t.exclusiveness
	parameters["topic_count"] = t.clusterCount
	parameters["word_count"] = t.wordCount

	// Add time range parameters.
	time := make(map[string]int64)
	time["from"] = t.timeFrom
	time["to"] = t.timeTo
	parameters["time"] = time

	return parameters
}

func (t *TopicTile) GetTileId() string {
	return fmt.Sprintf("%v/%v/%v", t.z, t.x, t.y)
}
