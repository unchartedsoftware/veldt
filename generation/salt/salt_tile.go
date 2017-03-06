package salt

import (
	"fmt"
	"strings"
	"encoding/json"
	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
)

// Tile represents any tile served to Veldt by Salt
type Tile struct {
	rmqConfig *Configuration // The configuration defining how we connect to the RabbitMQ server
	tileConfiguration map[string]interface{} // The JSON description of the tile configuration
}

func joinErrors (input []error) error {
	errorStrings := make([]string, len(input))
	for i, e := range input {
		errorStrings[i] = e.Error()
	}
	return fmt.Errorf(strings.Join(errorStrings, "\n"))
}

// NewSaltTile returns a constructor for salt-based tiles of all sorts.  It also initializes the
// salt server with the datasets it expects to use.
func NewSaltTile (rmqConfig *Configuration, datasetConfigurations ...[]byte) veldt.TileCtor {
	return func() (veldt.Tile, error) {
		t := &Tile{}
		t.rmqConfig = rmqConfig

		// Send any dataset configurations to salt immediately
		// Need a connection for that
		connection, err := NewConnection(t.rmqConfig)
		if err != nil {
			return nil, err
		}

		errors := make([]error, 0)

		for _, datasetConfig := range datasetConfigurations {
			_, err = connection.Dataset(datasetConfig)
			if nil != err {
				errors = append(errors, fmt.Errorf("Error registering dataset: %v", err))
			}
		}

		if len(errors) > 0 {
			return t, joinErrors(errors)
		}

		return t, nil
	}
}

// Parse stores tile parameters so that they can be sent to Salt when the tile request is made
func (t *Tile) Parse (params map[string]interface{}) error {
	t.tileConfiguration = params
	return nil
}

// Create generates a tile from the provided URI, tile coordinate, and query parameters
func (t *Tile) Create (uri string, coord *binning.TileCoord, query veldt.Query) ([]byte, error) {
	connection, err := NewConnection(t.rmqConfig)
	if err != nil {
		return nil, err
	}

	saltQuery, ok := query.(*Query)
	if !ok {
		return nil, fmt.Errorf("query is not salt.Query")
	}

	fullConfiguration := make(map[string]interface{})
	fullConfiguration["tile"] = t.tileConfiguration
	fullConfiguration["query"] = saltQuery.GetQueryConfiguration()

	configBytes, err := json.Marshal(fullConfiguration)
	if nil != err {
		return nil, fmt.Errorf("Tile or query configuration could not be marshalled into JSON for transport to salt")
	}

	return connection.QueryTiles(configBytes)
}
