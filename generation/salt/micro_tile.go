package salt

import (
	"fmt"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
	"github.com/unchartedsoftware/veldt/generation/batch"
	"github.com/unchartedsoftware/veldt/tile"
	"github.com/unchartedsoftware/veldt/util/json"
)

// MicroTile represents a Salt implementation of the micro tile
type MicroTile struct {
	TileData
	tile.Bivariate
	tile.Micro
	tile.TopHits
}

// NewMicroTile instantiates and returns a new tile struct.
func NewMicroTile(rmqConfig *Config, datasetConfigs ...[]byte) veldt.TileCtor {
	setupConnection(rmqConfig, datasetConfigs...)

	return func() (veldt.Tile, error) {
		Infof("New micro tile constructor request")
		return newMicroTile(rmqConfig), nil
	}
}

// NewMicroTileFactory instantiates and returns a factory for creating batched micro tiles.
func NewMicroTileFactory(rmqConfig *Config, datasetConfigs ...[]byte) batch.TileFactoryCtor {
	setupConnection(rmqConfig, datasetConfigs...)

	return func() (batch.TileFactory, error) {
		Infof("New micro tile factory constructor request")
		return newMicroTile(rmqConfig), nil
	}
}

func newMicroTile(rmqConfig *Config) *MicroTile {
	mt := &MicroTile{}
	mt.tileType = "micro"
	mt.rmqConfig = rmqConfig
	mt.buildConfig = func() (map[string]interface{}, error) {
		return mt.getTileConfig()
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

// parseMicroParams actually parses the provided JSON object, and
// populates the tile attributes.
func (m *MicroTile) parseMicroParams(params map[string]interface{}) error {
	if err := m.Bivariate.Parse(params); err != nil {
		return err
	}
	if err := m.TopHits.Parse(params); err != nil {
		return err
	}
	if err := m.Micro.Parse(params); err != nil {
		return err
	}
	// parse includes
	m.TopHits.IncludeFields = m.Micro.ParseIncludes(
		m.TopHits.IncludeFields,
		m.Bivariate.XField,
		m.Bivariate.YField)
	return nil
}

// GetTileConfig gets the configuration to send to Salt, so that it can
// construct the currently requested tile
func (m *MicroTile) getTileConfig() (map[string]interface{}, error) {
	err := m.parseMicroParams(*m.parameters)
	if err != nil {
		return nil, err
	}
	// Bounds are ignored - salt needs the dataset bounds, not the tile bounds
	// in visualization space
	return map[string]interface{}{
		"type":          "micro",
		"xField":        m.XField,
		"yField":        m.YField,
		"sortField":     m.SortField,
		"sortOrder":     m.SortOrder,
		"hitsCount":     m.HitsCount,
		"includeFields": m.IncludeFields,
	}, nil
}

func (m *MicroTile) convertTile(coord *binning.TileCoord, input []byte) ([]byte, error) {
	// Make sure our parameters are in sync with the current tile
	err := m.parseMicroParams(*m.parameters)
	if nil != err {
		return nil, err
	}

	rawHits, err := json.UnmarshalArray(input)
	if nil != err {
		return nil, err
	}

	numHits := len(rawHits)
	points := make([]float32, numHits*2)
	hits := make([]map[string]interface{}, numHits)

	for i, hit := range rawHits {
		x, ok := json.GetFloat(hit, "x")
		if !ok {
			return nil, fmt.Errorf("could not parse `x` from hit: %v", hit)
		}
		y, ok := json.GetFloat(hit, "y")
		if !ok {
			return nil, fmt.Errorf("could not parse `y` from hit: %v", hit)
		}
		points[2*i+0] = float32(x)
		points[2*i+1] = 256 - float32(y)

		hitMap, ok := json.GetChild(hit, "values")
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
