package salt

import (
	"fmt"
	"math"
	"encoding/json"
	"encoding/binary"

	"github.com/liyinhgqw/typesafe-config/parse"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
	"github.com/unchartedsoftware/veldt/generation/batch"
 	"github.com/unchartedsoftware/plog"
)

// Tile represents any tile served to Veldt by Salt
type Tile struct {
	 // The configuration defining how we connect to the RabbitMQ server
	rmqConfig *Configuration
	// The default configuration to merge into any parameters that are passed in
	defaultTileConfig map[string]interface{}
	// The The full configuration of the requested tile
	tileConfig map[string]interface{}
}

var datasets = make(map[string]string)

func getDatasetName (datasetConfigRaw []byte) (string, error) {
	datasetConfigMap, err := parse.Parse("dataset", string(datasetConfigRaw))
	if err != nil {
		return "", err
	}
	datasetConfig := datasetConfigMap.GetConfig()
	result, err := datasetConfig.GetString("name")
	if err != nil {
		return "", err
	}
	return stripTerminalQuotes(result), nil
}

// NewSaltTile returns a constructor for salt-based tiles of all sorts.  It also initializes the
// salt server with the datasets it expects to use.
//
// rmqConfig The configuration information needed to connect to RabbitMQ, through
//           which to connect to the salt server
// defaultTileConfig A default configuration that will get merged into any
//                   parameters passed into tile requests
// datasetConfigurations Any dataset configurations that will be needed for tiles using this constructor
func NewSaltTile (rmqConfig *Configuration,
	defaultTileConfig map[string]interface{},
	datasetConfigs ...[]byte) veldt.TileCtor {
		setupConnection(rmqConfig, datasetConfigs...)

		return func() (veldt.Tile, error) {
			log.Infof(preLog+"new tile constructor request")
			t := &Tile{}
			t.rmqConfig = rmqConfig
			t.defaultTileConfig = defaultTileConfig
			return t, nil
		}
	}

// NewSaltTileFactory returns a constructor for a tile factory for salt-based
// tiles of all sorts.  It also initializes the salt server with the datasets
// it expects to use.
//
// This is identical to NewSaltTile except in the type it assigns to its
// return value.
// TODO: Work this type into TileCTor perhaps?
//
// rmqConfig The configuration information needed to connect to RabbitMQ, through
//           which to connect to the salt server
// defaultTileConfig A default configuration that will get merged into any
//                   parameters passed into tile requests
// datasetConfigurations Any dataset configurations that will be needed for tiles using this constructor
func NewSaltTileFactory (rmqConfig *Configuration,
	defaultTileConfig map[string]interface{},
	datasetConfigs ...[]byte) batch.TileFactoryCtor {
		setupConnection(rmqConfig, datasetConfigs...)

		return func() (batch.TileFactory, error) {
			log.Infof(preLog+"new tile factory constructor request")
			tf := &Tile{}
			tf.rmqConfig = rmqConfig
			tf.defaultTileConfig = defaultTileConfig
			return tf, nil
		}
	}

func setupConnection (rmqConfig *Configuration, datasetConfigs ...[]byte) {
	// Send any dataset configurations to salt immediately
	// Need a connection for that
	connection, err := NewConnection(rmqConfig)
	if err != nil {
		log.Errorf("Error connecting to salt server to configure datasets: %v", err)
	} else {
		for _, datasetConfig := range datasetConfigs {
			name, err := getDatasetName(datasetConfig)
			if nil != err {
				log.Errorf("Error registering dataset: can't find name of dataset %v", string(datasetConfig))
			} else {
				_, err = connection.Dataset(datasetConfig)
				if nil != err {
					log.Errorf("Error registering dataset %v: %v", name, err)
				} else {
						datasets[name] = string(datasetConfig)
				}
			}
		}
	}
}

// Parse stores tile parameters so that they can be sent to Salt when the tile request is made
func (t *Tile) Parse (params map[string]interface{}) error {
	t.tileConfig = params
	return nil
}

// Create generates a single tile from the provided URI, tile coordinate, and query parameters
func (t *Tile) Create (uri string, coord *binning.TileCoord, query veldt.Query) ([]byte, error) {
	responseChan := make(chan batch.TileResponse, 1)
	request := &batch.TileRequest{t.tileConfig, uri, coord, query, responseChan}
	t.CreateTiles([]*batch.TileRequest{request})
	response := <-responseChan
	if nil != response.Tile {
		log.Debugf("Create: Got response tile of length %d", len(response.Tile))
	} else {
		log.Debugf("Create: Got nil response tile")
	}
	if nil != response.Err {
		log.Debugf("Create: Got non-nil error")
	} else {
		log.Debugf("Create: no error")
	}
	return response.Tile, response.Err
}



type separateTileRequest struct {
	coord *binning.TileCoord
	sendTo chan batch.TileResponse
}
type jointRequest struct {
	tileConfig map[string]interface{}
	query map[string]interface{}
	dataset string
	tiles []*separateTileRequest
}

func canMerge (a, b *jointRequest) bool {
	if !configurationsEqual(a.tileConfig, b.tileConfig) {
		return false
	}
	if !configurationsEqual(a.query, b.query) {
		return false
	}
	return a.dataset == b.dataset
}

func (j *jointRequest) merge (from *jointRequest) {
	for _, tile := range from.tiles {
		j.tiles = append(j.tiles, tile)
	}
}

func (t *Tile) extractJointRequest (request *batch.TileRequest) *jointRequest {
	tileConfig := mergeConfigurations(request.Parameters, t.defaultTileConfig)

	var queryConfig map[string]interface{}
	if nil != request.Query {
		saltQuery, ok := request.Query.(*Query)
		if !ok {
			log.Errorf(preLog+"Query for salt tile was not a salt query")
		} else {
			queryConfig = saltQuery.GetQueryConfiguration()
		}
	}

	separateRequest := separateTileRequest{request.Coordinates, request.ResultChannel}
	separateRequests := []*separateTileRequest{&separateRequest}
	
	return &jointRequest{tileConfig, queryConfig, request.URI, separateRequests}
}

// CreateTiles generates multiple tiles from the provided information
func (t *Tile) CreateTiles (requests []*batch.TileRequest) {
	log.Infof(preLog+"CreateTiles: Processing %d requests\n", len(requests))
	// Create our connection
	connection, err := NewConnection(t.rmqConfig)
	if err != nil {
		for _, request := range requests {
			request.ResultChannel <- batch.TileResponse{nil, err}
		}
		return
	}

	// For requests to be grouped, they have to share tile configuration, query
	// configuration, and dataset - i.e., uri - for the moment.
	//
	// Eventually, we should be able to eliminate the need to share tile
	// configuration, so as to get multiple layers in a single tiling, but
	// that's a secondary consideration
	// 
	// Ideally, we'd just do this with maps, but GO doesn't support complex map
	// keys, so we're stuck doing this the hard way
	consolidatedRequests := make([]*jointRequest, 0)
	for _, tileRequest := range requests {
		request := t.extractJointRequest(tileRequest)
		requestMerged := false
		for _, currentRequest := range consolidatedRequests {
			if !requestMerged && canMerge(request, currentRequest) {
				currentRequest.merge(request)
				requestMerged = true
			}
		}
		if !requestMerged {
			consolidatedRequests = append(consolidatedRequests, request)
		}
	}

	// Requests are all merged
	// Now actually make our requests of the server
	for _, request := range consolidatedRequests {
		log.Infof(preLog+"Request for %d tiles for dataset %s\n", len(request.tiles), request.dataset)
		// Create our consolidated configuration
		fullConfig := make(map[string]interface{})
		fullConfig["tile"] = request.tileConfig
		fullConfig["query"] = request.query
		fullConfig["dataset"] = datasets[request.dataset]
		// Put in all our tile requests, recording our response channel for each as we go
		responseChannels := make(map[string]chan batch.TileResponse)
		tileSpecs := make([]interface{}, 0)
		for _, tileReq := range request.tiles {
			c := tileReq.coord
			tileSpec := make(map[string]interface{})
			tileSpec["level"] = int(c.Z)
			tileSpec["x"] = int(c.X)
			tileSpec["y"] = int(c.Y)
			tileSpecs = append(tileSpecs, tileSpec)
			responseChannels[coordToString(int(c.Z), int(c.X), int(c.Y))] = tileReq.sendTo
		}
		fullConfig["tile-specs"] = tileSpecs

		// Marshal the consolidated request into a string
		requestBytes, err := json.Marshal(fullConfig)
		if nil != err {
			err := fmt.Errorf("Tile request(s) could not be marshalled into JSON for transport to salt")
			for _, tileReq := range request.tiles {
				tileReq.sendTo <- batch.TileResponse{nil, err}
			}
			return
		}

		// Send the marshalled request to Salt, and await a response
		result, err := connection.QueryTiles(requestBytes)
		// Unpack the results
		tiles := unpackTiles(result)
		for key, channel := range responseChannels {
			tile, ok := tiles[key]
			if ok {
				log.Infof("Found tile for key %s of length ", key, len(tile))
				channel <- batch.TileResponse{tile, nil}
			} else {
				// No tile, but no error either
				log.Infof("No tile found for key %s", key)
				channel <- batch.TileResponse{nil, nil}
			}
		}
	}
}

// Get a unique string ID for use in maps for a tile coordinate
func coordToString (level, x, y int) string {
    max := 1 << uint64(level)
	digits := int64(math.Floor(math.Log10(float64(max)))) + 1
	format := fmt.Sprintf("%%02d:%%0%dd:%%0%dd", digits, digits)
	return fmt.Sprintf(format, level, x, y)
}

// unpackTiles unpacks the message sent to us by salt into a series of tiles,
// keyed by the coordToString function above
func unpackTiles (saltMsg []byte) map[string][]byte {
	p := 0
	maxP := len(saltMsg)
	results := make(map[string][]byte)
	for p < maxP {
		level := binary.BigEndian.Uint64(saltMsg[p:p+8])
		p = p + 8
		x     := binary.BigEndian.Uint64(saltMsg[p:p+8])
		p = p + 8
		y     := binary.BigEndian.Uint64(saltMsg[p:p+8])
		p = p + 8
		size  := int(binary.BigEndian.Uint64(saltMsg[p:p+8]))
		p = p + 8
		key := coordToString(int(level), int(x), int(y))
		log.Infof("Unpacking tile [%d: %d, %d] = %s", level, x, y, key)
		results[key] = saltMsg[p:p+size]
		p = p + size
	}
	return results
}

// MergeConfigurations takes two configuration mappings (created from JSON,
// presumably) and merges them into one.
// 
// If both input maps have values for a given key, the first map wins.
func mergeConfigurations (a, b map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// Copy in values from A
	for k, v := range a {
		// generally, start by just taking the value from A
		result[k] = v

		// but if the value from A is a map, and there is a value from B, and
		// it's a map too, we need to merge them.
		subMapA, ok := v.(map[string]interface{})
		if ok {
			valB, ok := b[k]
			if ok {
				subMapB, ok := valB.(map[string]interface{})
				if ok {
					result[k] = mergeConfigurations(subMapA, subMapB)
				}
			}
		}
	}

	// Copy in values from B that won't overwrite existing values
	for k, v := range b {
		_, ok := result[k]
		if !ok {
			result[k] = v
		}
	}

	return result
}

func configurationsEqual (a, b map[string]interface{}) bool {
	// Check keys
	if len(a) != len(b) {
		return false
	}
	for k, valA := range a {
		valB, ok := b[k]
		if !ok {
			return false
		}
		subMapA, isSubMapA := valA.(map[string]interface{})
		subMapB, isSubMapB := valB.(map[string]interface{})
		if isSubMapA && isSubMapB {
			if !configurationsEqual(subMapA, subMapB) {
				return false
			}
		} else if isSubMapA || isSubMapB {
			return false
		} else if valA != valB {
			return false
		}
	}
	return true
}
