package salt

import (
	"encoding/binary"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
	"github.com/unchartedsoftware/veldt/generation/batch"
	"github.com/unchartedsoftware/veldt/tile"
)

// MacroTile represents a Salt implementation of the macro tile
type MacroTile struct {
	TileData
	tile.Bivariate
	tile.Macro
}

// NewMacroTile instantiates and returns a new tile struct.
func NewMacroTile(rmqConfig *Config, datasetConfigs ...[]byte) veldt.TileCtor {
	setupConnection(rmqConfig, datasetConfigs...)

	return func() (veldt.Tile, error) {
		Infof("New macro tile constructor request")
		return newMacroTile(rmqConfig), nil
	}
}

// NewMacroTileFactory instantiates and returns a factory for creating batched macro tiles.
func NewMacroTileFactory(rmqConfig *Config, datasetConfigs ...[]byte) batch.TileFactoryCtor {
	setupConnection(rmqConfig, datasetConfigs...)

	return func() (batch.TileFactory, error) {
		Infof("New macro tile factory constructor request")
		return newMacroTile(rmqConfig), nil
	}
}

func newMacroTile(rmqConfig *Config) *MacroTile {
	mt := &MacroTile{}
	mt.tileType = "macro"
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
func (m *MacroTile) Parse(params map[string]interface{}) error {
	return m.TileData.Parse(params)
}

// parseMacroParams actually parses the provided JSON object, and
// populates the tile attributes.
func (m *MacroTile) parseMacroParams(params map[string]interface{}) error {
	if err := m.Bivariate.Parse(params); err != nil {
		return err
	}
	return m.Macro.Parse(params)
}

// GetTileConfig gets the configuration to send to Salt, so that it can
// construct the currently requested tile
func (m *MacroTile) getTileConfig() (map[string]interface{}, error) {
	err := m.parseMacroParams(*m.parameters)
	if err != nil {
		return nil, err
	}
	// Bounds are ignored - salt needs the dataset bounds, not the tile bounds
	// in visualization space
	return map[string]interface{}{
		"type":       "macro",
		"xField":     m.XField,
		"yField":     m.YField,
		"resolution": m.Resolution,
	}, nil
}

func (m *MacroTile) convertTile(coord *binning.TileCoord, input []byte) ([]byte, error) {
	err := m.parseMacroParams(*m.parameters)
	if err != nil {
		return nil, err
	}

	// Bin characteristics
	binSize := binning.MaxTileResolution / float64(m.Resolution)
	halfSize := float64(binSize / 2)

	// Macro tiles are returned to us as a series of integers which indicate
	// the x and y coordinates of the populated bins
	// Current position in our data
	p := 0
	// First we have one integer representing the number of points
	numPoints := int(binary.LittleEndian.Uint32(input))
	p = p + 4

	points := make([]float32, numPoints*2)
	for i := 0; i < numPoints; i++ {
		x := binary.LittleEndian.Uint32(input[p : p+4])
		p = p + 4
		y := uint32(m.Resolution) - binary.LittleEndian.Uint32(input[p:p+4])
		p = p + 4

		// Convert from bin number to location
		// X
		points[i*2+0] = float32(float64(x)*binSize + halfSize)
		// Y
		points[i*2+1] = float32(float64(y)*binSize + halfSize)
	}

	return m.Macro.Encode(points)
}

func (m *MacroTile) buildDefaultTile() ([]byte, error) {
	return m.Macro.Encode(make([]float32, 0))
}
