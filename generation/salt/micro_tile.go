package salt

import (
	"encoding/json"
	"fmt"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
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
func NewMicroTile(rmqConfig *Configuration, datasetConfigs ...[]byte) veldt.TileCtor {
	setupConnection(rmqConfig, datasetConfigs...)

	return func() (veldt.Tile, error) {
		saltInfof("New micro tile constructor request")
		return newMicroTile(rmqConfig), nil
	}
}

// NewMicroTileFactory instantiates and returns a factory for creating batched micro tiles.
func NewMicroTileFactory(rmqConfig *Configuration, datasetConfigs ...[]byte) batch.TileFactoryCtor {
	setupConnection(rmqConfig, datasetConfigs...)

	return func() (batch.TileFactory, error) {
		saltInfof("New micro tile factory constructor request")
		return newMicroTile(rmqConfig), nil
	}
}

func newMicroTile(rmqConfig *Configuration) *MicroTile {
	mt := &MicroTile{}
	mt.tileType = "micro"
	mt.rmqConfig = rmqConfig
	mt.buildConfig = func() (map[string]interface{}, error) {
		return mt.getTileConfiguration()
	}
	mt.convert = func(coord *binning.TileCoord, input []byte) ([]byte, error) {
		return mt.convertTile(coord, input)
	}
	mt.buildDefault = func() ([]byte, error) {
		return mt.buildDefaultTile()
	}
	return mt
}

// Parse does the standard salt tile parsing of parameters - i.e., saving them for later
func (m *MicroTile) Parse(params map[string]interface{}) error {
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
func (m *MicroTile) getTileConfiguration() (map[string]interface{}, error) {
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

func (m *MicroTile) convertTile(coord *binning.TileCoord, input []byte) ([]byte, error) {
	var rawHits []map[string]interface{}
	err := json.Unmarshal(input, &rawHits)
	if nil != err {
		return nil, err
	}

	numHits := len(rawHits)
	points := make([]float32, numHits*2)
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
		points[2*i+1] = 256 - y

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

func (m *MicroTile) buildDefaultTile() ([]byte, error) {
	return m.Micro.Encode(make([]map[string]interface{}, 0), make([]float32, 0))
}
