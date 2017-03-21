package salt

import (
	"fmt"
	"encoding/json"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/generation/batch"
	"github.com/unchartedsoftware/veldt/tile"
)

// MicroTile represents a Salt implementation of the micro tile
type MicroTile struct {
	TileData
	tile.Bivariate
	tile.Micro
	tile.TopHits
}

// NewMicroTile instantiates and returns a new tile struct.
func NewMicroTile (rmqConfig *Configuration, datasetConfigs ...[]byte) veldt.TileCtor {
	setupConnection(rmqConfig, datasetConfigs...)

	return func() (veldt.Tile, error) {
		saltInfof("New micro tile constructor request")
		t := &MicroTile{}
		t.rmqConfig = rmqConfig
		t.buildConfig = func () (map[string]interface{}, error) {
			return t.getTileConfiguration()
		}
		t.convert = func (input []byte) ([]byte, error) {
			return t.convertResults(input)
		}
		return t, nil
	}
}

// NewMicroTileFactory instantiates and returns a factory for creating batched micro tiles.
func NewMicroTileFactory (rmqConfig *Configuration, datasetConfigs ...[]byte) batch.TileFactoryCtor {
	setupConnection(rmqConfig, datasetConfigs...)

	return func() (batch.TileFactory, error) {
		saltInfof("New micro tile factory constructor request")
		tf := &MicroTile{}
		tf.rmqConfig = rmqConfig
		tf.buildConfig = func () (map[string]interface{}, error) {
			return tf.getTileConfiguration()
		}
		tf.convert = func (input []byte) ([]byte, error) {
			return input, nil
		}
		return tf, nil
	}
}

// Parse does the standard salt tile parsing of parameters - i.e., saving them for later
func (m *MicroTile) Parse (params map[string]interface{}) error {
	return m.TileData.Parse(params)
}

// parseMicroParameters actually parses the provided JSON object, and
// populates the tile attributes.
func (m *MicroTile) parseMicroParameters(params map[string]interface{}) error {
	if err := m.Bivariate.Parse(params); nil != err {
		return err
	}
	if err := m.TopHits.Parse(params); nil != err {
		return err
	}
	if err := m.Micro.Parse(params); nil != err {
		return err
	}
	// parse includes
	m.TopHits.IncludeFields = m.Micro.ParseIncludes(
		m.TopHits.IncludeFields,
		m.Bivariate.XField,
		m.Bivariate.YField)
	return nil
}

// GetTileConfiguration gets the configuration to send to Salt, so that it can
// construct the currently requested tile
func (m *MicroTile) getTileConfiguration () (map[string]interface{}, error) {
	err := m.parseMicroParameters(*m.parameters)
	if nil != err {
		return nil, err
	}

	result := make(map[string]interface{})

	setProperty("type", "micro", result)

	// Bivariate properties
	setProperty("xField", m.XField, result)
	setProperty("yField", m.YField, result)
	// Resolution is ignored for micro, as it is irrelevant
	// setProperty("resolution", m.Resolution, result)
	// Bounds are ignored - salt needs the dataset bounds, not the tile bounds
	// in visualization space
	// setProperty("bounds.left",   m.Left, result)
	// setProperty("bounds.right",  m.Right, result)
	// setProperty("bounds.top",    m.Top, result)
	// setProperty("bounds.bottom", m.Bottom, result)

	// TopHits properties
	setProperty("sortField", m.SortField, result)
	setProperty("sortOrder", m.SortOrder, result)
	setProperty("hitsCount", m.HitsCount, result)
	setProperty("includeFields", m.IncludeFields, result)
	
	return result, nil
}


func (m *MicroTile) convertResults (input []byte) ([]byte, error) {
	var rawHits []map[string]interface{}
	err := json.Unmarshal(input, &rawHits)
	if nil != err {
		return nil, err
	}

	numHits := len(rawHits)
	points := make([]float32, numHits * 2)
	hits := make([]map[string]interface{}, numHits)

	for i, hit := range rawHits {
		x, err := getFloat32Property("x", hit)
		if nil != err {
			return nil, err
		}
		y, err := getFloat32Property("y", hit)
		if nil != err {
			return nil, err
		}
		points[2*i+0] = x
		points[2*i+1] = y

		hitMapRaw, err := getProperty("values", hit)
		if nil != err {
			return nil, err
		}
		hitMap, ok := hitMapRaw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("values didn't form a map: %v", hit)
		}
		hits[i] = hitMap
	}
	
	return m.Micro.Encode(hits, points)
}

