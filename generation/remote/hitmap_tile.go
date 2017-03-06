package remote

import (
	"encoding/json"
	"fmt"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
	jsonUtil "github.com/unchartedsoftware/veldt/util/json"
)

type HitmapTile struct {
    APITile
	timeFrom       int64
	timeTo         int64
}

func NewHitmapTile() veldt.TileCtor {
	return func() (veldt.Tile, error) {
		return &HitmapTile{}, nil
	}
}

func (t *HitmapTile) Parse(params map[string]interface{}) error {
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

	t.timeFrom = int64(timeFrom)
	t.timeTo = int64(timeTo)

	// Use hash of all parameters to identify request.
	requestId := fmt.Sprintf("%v::%v::%v::%v::%v::%v::%v",
		timeFrom, timeTo)

	t.requestId = requestId
    t.tileType = "hitmap"

	return nil
}

func (t *HitmapTile) GetApiUrl() string {
	// TODO: Have the URL configurable!
	return "http://163.152.20.64:5000/GET_HITMAP/test"
}

func (t *HitmapTile) Create(uri string, coord *binning.TileCoord, query veldt.Query) ([]byte, error) {
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
	var tmpHitmapParsed interface{}
	err = json.Unmarshal([]byte(res.(string)), &tmpHitmapParsed)
	if err != nil {
		return nil, err
	}

	hitmapParsed, ok := tmpHitmapParsed.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Unexpected response format from topic modelling service for hitmap tile: incorrect structure in %v", res)
	}

	exclusiveness := make(map[string]float64)
	exclusivenessValues, ok := jsonUtil.GetArray(hitmapParsed, "exclusiveness")
	for _, ex := range exclusivenessValues {
		exMap, ok := ex.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Unexpected response format from topic modelling service: incorrect exclusiveness structure in %v", res)
		}

        date, ok := jsonUtil.GetString(exMap, "date")
        if !ok {
            return nil, fmt.Errorf("Unexpected response format from topic modelling service: incorrect date structure in %v", res)
        }

        value, ok := jsonUtil.GetNumber(exMap, "value")
        if !ok {
            return nil, fmt.Errorf("Unexpected response format from topic modelling service: incorrect value structure in %v", res)
        }

        exclusiveness[date] = value
	}

	// marshal results
	return json.Marshal(exclusiveness)
}

// Create the request to the remote service.
func (t *HitmapTile) GetRequestParameters() map[string]interface{} {
	parameters := make(map[string]interface{})

	// Add time range parameters.
	time := make(map[string]int64)
	time["from"] = t.timeFrom
	time["to"] = t.timeTo
	parameters["date"] = time

	return parameters
}

func (t *HitmapTile) GetTileId() string {
	return fmt.Sprintf("%v/%v/%v", t.z, t.x, t.y)
}
