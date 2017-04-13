package salt

import (
	"encoding/binary"
	"fmt"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
	"github.com/unchartedsoftware/veldt/generation/batch"
	"github.com/unchartedsoftware/veldt/tile"
	"github.com/unchartedsoftware/veldt/util/json"
)

// CountTile represents a Salt implementation of the count tile
type CountTile struct {
	tile.Bivariate
	TileData
	valueField string
}

// NewCountTile instantiates and returns a new tile struct.
func NewCountTile(rmqConfig *Config, datasetConfigs ...[]byte) veldt.TileCtor {
	setupConnection(rmqConfig, datasetConfigs...)

	return func() (veldt.Tile, error) {
		Infof("New count tile constructor request")
		return newCountTile(rmqConfig), nil
	}
}

// NewCountTileFactory instantiates and returns a factory for creating batched count tiles.
func NewCountTileFactory(rmqConfig *Config, datasetConfigs ...[]byte) batch.TileFactoryCtor {
	setupConnection(rmqConfig, datasetConfigs...)

	return func() (batch.TileFactory, error) {
		Infof("New count tile factory constructor request")
		return newCountTile(rmqConfig), nil
	}
}

func newCountTile(rmqConfig *Config) *CountTile {
	ct := &CountTile{}
	ct.tileType = "count"
	ct.rmqConfig = rmqConfig
	ct.buildConfig = func() (map[string]interface{}, error) {
		return ct.getTileConfig()
	}
	ct.convert = func(coord *binning.TileCoord, input []byte) ([]byte, error) {
		return ct.convertTile(coord, input)
	}
	ct.buildDefault = func() ([]byte, error) {
		return ct.buildDefaultTile()
	}
	return ct
}

// Parse does the standard salt tile parsing of parameters - i.e., saving them for later
func (c *CountTile) Parse(params map[string]interface{}) error {
	return c.TileData.Parse(params)
}

// parseCountParams actually parses the provided JSON object, and
// populates the tile attributes.
func (c *CountTile) parseCountParams(params map[string]interface{}) error {
	valueField, ok := json.GetString(params, "valueField")
	if ok {
		c.valueField = valueField
	} else {
		c.valueField = ""
	}
	return c.Bivariate.Parse(params)
}

// GetTileConfig gets the configuration to send to Salt, so that it can
// construct the currently requested tile
func (c *CountTile) getTileConfig() (map[string]interface{}, error) {
	err := c.parseCountParams(*c.parameters)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})

	setProperty("type", "heatmap", result)
	setProperty("xField", c.XField, result)
	setProperty("yField", c.YField, result)
	if 0 < len(c.valueField) {
		setProperty("valueField", c.valueField, result)
	}
	setProperty("resolution", 1, result)
	// Bounds are ignored - salt needs the dataset bounds, not the tile bounds
	// in visualization space
	// setProperty("bounds.left",   c.Left, result)
	// setProperty("bounds.right",  c.Right, result)
	// setProperty("bounds.top",    c.Top, result)
	// setProperty("bounds.bottom", c.Bottom, result)

	return result, nil
}

func (c *CountTile) convertTile(coord *binning.TileCoord, input []byte) ([]byte, error) {
	count := binary.LittleEndian.Uint32(input)
	return []byte(fmt.Sprintf(`{"count":%d}`, count)), nil
}

func (c *CountTile) buildDefaultTile() ([]byte, error) {
	return []byte(`{"count":0}`), nil
}
