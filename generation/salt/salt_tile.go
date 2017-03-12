package salt

import (
	"fmt"
	"encoding/json"

	"github.com/liyinhgqw/typesafe-config/parse"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
	"github.com/unchartedsoftware/plog"
)

// Tile represents any tile served to Veldt by Salt
type Tile struct {
	rmqConfig *Configuration // The configuration defining how we connect to the RabbitMQ server
	tileType string          // The type of tile to produce, as interpretted by the salt server
	tileConfiguration map[string]interface{} // The JSON description of the tile configuration
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
func NewSaltTile (rmqConfig *Configuration, datasetConfigurations ...[]byte) veldt.TileCtor {
	// Send any dataset configurations to salt immediately
	// Need a connection for that
	connection, err := NewConnection(rmqConfig)
	if err != nil {
		log.Errorf("Error connecting to salt server to configure datasets: %v", err)
	} else {
		for _, datasetConfig := range datasetConfigurations {
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

	return func() (veldt.Tile, error) {
		log.Infof(preLog+"new tile constructor request")
		t := &Tile{}
		t.rmqConfig = rmqConfig
		return t, nil
	}
}

// Parse stores tile parameters so that they can be sent to Salt when the tile request is made
func (t *Tile) Parse (tileType string, params map[string]interface{}) error {
	t.tileType = tileType
	t.tileConfiguration = params
	return nil
}

// Create generates a tile from the provided URI, tile coordinate, and query parameters
func (t *Tile) Create (uri string, coord *binning.TileCoord, query veldt.Query) ([]byte, error) {
	connection, err := NewConnection(t.rmqConfig)
	if err != nil {
		return nil, err
	}

	var saltQueryConfig map[string]interface{}
	if nil != query {
		saltQuery, ok := query.(*Query)
		if !ok {
			return nil, fmt.Errorf("query is not salt.Query")
		}
		saltQueryConfig = saltQuery.GetQueryConfiguration()
	}

	coordMap := make(map[string]interface{})
	coordMap["x"] = coord.X
	coordMap["y"] = coord.Y
	coordMap["level"] = coord.Z
	tileSpec := make(map[string]interface{})
	tileSpec["type"] = t.tileType
	tileSpec["coordinates"] = coordMap

	fullConfiguration := make(map[string]interface{})
	fullConfiguration["tile-spec"] = tileSpec
	fullConfiguration["tile"] = t.tileConfiguration
	fullConfiguration["query"] = saltQueryConfig
	fullConfiguration["dataset"] = datasets[uri]

	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println("Sending tile request")
	fmt.Println("Full configuration:")
	fmt.Println(fullConfiguration)
	fmt.Println()
	fmt.Printf("Dataset is \"%v\" (%v)\n", uri, []byte(uri))
	fmt.Println("Dataset keys are:")
	for k := range datasets {
		fmt.Printf("\t\"%v\" = (%v)\n", k, []byte(k))
	}
	fmt.Println()
	fmt.Printf("uri's dataset: %v\n", datasets[uri])
	fmt.Println()
	fmt.Println()

	configBytes, err := json.Marshal(fullConfiguration)
	if nil != err {
		return nil, fmt.Errorf("Tile or query configuration could not be marshalled into JSON for transport to salt")
	}

	return connection.QueryTiles(configBytes)
}
