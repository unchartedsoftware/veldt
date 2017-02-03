package remote

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/parnurzeal/gorequest"
)

var (
	mutex    = sync.Mutex{} //May want to create a mutex / client.
	handlers = make(map[string]*ServiceClient)
)

// A service client exists to handle a batch of tile requests.
func getServiceClient(tile *TopicTile) *ServiceClient {
	requestId := tile.requestId

	// Not the most efficient use of locks. Will probably need to be optimized.
	mutex.Lock()
	defer mutex.Unlock()
	if client, ok := handlers[requestId]; ok {
		return client
	}

	client := NewServiceClient()
	client.requestId = requestId
	handlers[requestId] = client

	go client.HandleRequests()

	return client
}

// Handles a batch of topic tile requests.
type ServiceClient struct {
	requestId        string
	requests         []*TopicTile
	responseChannels map[string]chan interface{}
	requestCount     int
	processing       bool
}

// Create the new service client. Only initializes the response channels
// since we do not know how many topic tiles are expected yet.
func NewServiceClient() *ServiceClient {
	service := ServiceClient{}
	service.requests = make([]*TopicTile, 0)
	service.responseChannels = make(map[string]chan interface{})
	service.processing = false

	return &service
}

// Add a request to the batch.
func (c *ServiceClient) AddRequest(tile *TopicTile) (chan interface{}, error) {
	// Could get duplicate requests to return the same channel that is already setup.
	// TODO: Need to actually handle the duplcate request case because right now Only
	// a single message gets written to the channel.
	hash := c.getTileCoordinateHash(tile.x, tile.y, tile.z)
	mutex.Lock()
	defer mutex.Unlock()
	if channel, ok := c.responseChannels[hash]; ok {
		return channel, nil
	}

	if c.processing {
		return nil, fmt.Errorf("Cannot add tile to request. Service is already processing.")
	}

	// Create the response channel.
	channel := make(chan interface{})
	c.responseChannels[hash] = channel

	// Append the request.
	if c.requestCount == 0 {
		c.requestCount = tile.tileCount
	}
	c.requests = append(c.requests, tile)

	return channel, nil
}

// Main function to be called on initialization in a separate thread.
func (c *ServiceClient) HandleRequests() {
	// Read the requests data.
	requests := c.getClientRequestsData()
	mutex.Lock()
	c.processing = true
	mutex.Unlock()

	// Send the request to the server.
	response, err := c.sendRequest(requests)
	if err != nil {
		c.handleError(err)
		return
	}

	// Parse the response from the server.
	parsedResponse, err := c.parseResponse([]byte(response))
	if err != nil {
		c.handleError(err)
		return
	}

	// Write the topics back to the tile handlers.
	// May be a bit abusive of the mutex.
	mutex.Lock()
	for hash, topics := range parsedResponse {
		tileChannel := c.responseChannels[hash]
		tileChannel <- topics
		// Remove the handler.
		delete(c.responseChannels, hash)
	}
	delete(handlers, c.requestId)
	mutex.Unlock()
}

func (c *ServiceClient) sendRequest(requestData map[string]interface{}) (string, error) {
	// TODO: Have the URL configurable!
	request := gorequest.New()
	_, result, errs := request.Post("http://163.152.20.64:5000/GET_TOPICS/test").Send(requestData).End()

	// TODO: Check the return values (the response) to make sure it is 200.
	// For now return the first error.
	if errs != nil && len(errs) > 0 {
		return "", errs[0]
	}

	return result, nil
}

// Handle errors raised during the processing of the batch.
func (c *ServiceClient) handleError(err error) {
	// Send the error to all tiles.
	for _, tile := range c.responseChannels {
		tile <- err
	}
}

// Create the request to the remote service.
func (c *ServiceClient) getClientRequestsData() map[string]interface{} {
	// Wait until all the requests have come in!
	c.waitForRequests()
	initialRequest := c.requests[0]

	// All tiles have the same parameters except for tile coordinates.
	terms := make(map[string][]string)
	terms["include"] = initialRequest.inclusionTerms
	terms["exclude"] = initialRequest.exclusionTerms
	tileData := make(map[string]interface{})
	tileData["terms"] = terms

	// Add simple parameters.
	tileData["exclusiveness"] = initialRequest.exclusiveness
	tileData["clusterCount"] = initialRequest.clusterCount
	tileData["wordCount"] = initialRequest.wordCount

	// Get the tile coordinates.
	coordinates := make([]interface{}, c.requestCount)
	for i, t := range c.requests {
		coordinates[i] = c.getTileCoordinate(t)
	}
	tileData["tiles"] = coordinates

	return tileData
}

// Parse the response from the remote service.
func (c *ServiceClient) parseResponse(response []byte) (map[string]string, error) {
	// Unmarshall the json and go through the expected structure.
	var raw interface{}
	err := json.Unmarshal([]byte(response), &raw)
	if err != nil {
		return nil, err
	}

	tiles, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("Unable to parse response from Topic Modelling service.")
	}

	// Extract the information by tile.
	parsed := make(map[string]string)
	for _, t := range tiles {
		// Pull the coordinate information to build the tile hash.
		tileData, ok := t.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Unable to parse response from Topic Modelling service.")
		}

		tileCoord, ok := tileData["tile"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Unable to parse response from Topic Modelling service.")
		}

		// Create the json for the tile response.
		tileString, err := json.Marshal(tileData)
		if err != nil {
			return nil, err
		}

		x := uint32(tileCoord["x"].(float64))
		y := uint32(tileCoord["y"].(float64))
		z := uint32(tileCoord["level"].(float64))
		parsed[c.getTileCoordinateHash(x, y, z)] = string(tileString)
	}

	return parsed, nil
}

// Create the coordinate structure from the tile.
func (c *ServiceClient) getTileCoordinate(tile *TopicTile) interface{} {
	coordinates := make(map[string]uint32)

	coordinates["x"] = tile.x
	coordinates["y"] = tile.y
	coordinates["level"] = tile.z

	return coordinates
}

// Create the hash for the tile.
func (c *ServiceClient) getTileCoordinateHash(x, y, z uint32) string {
	return fmt.Sprintf("%v/%v/%v", z, x, y)
}

// Wait until we have all the expected requests.
func (c *ServiceClient) waitForRequests() {
	for c.requestCount != len(c.requests) {
		time.Sleep(500 * time.Millisecond)
	}
}
