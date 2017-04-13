package salt

import (
	"fmt"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
	"github.com/unchartedsoftware/veldt/generation/batch"
	"github.com/unchartedsoftware/veldt/tile"
	"github.com/unchartedsoftware/veldt/util/json"
)

// HeatmapTile represents a Salt implementation of the heatmap tile
type HeatmapTile struct {
	tile.Bivariate
	TileData
	valueField string
}

// NewHeatmapTile instantiates and returns a new tile struct.
func NewHeatmapTile(rmqConfig *Config, datasetConfigs ...[]byte) veldt.TileCtor {
	setupConnection(rmqConfig, datasetConfigs...)

	return func() (veldt.Tile, error) {
		Infof("New heatmap tile constructor request")
		return newHeatmapTile(rmqConfig), nil
	}
}

// NewHeatmapTileFactory instantiates and returns a factory for creating batched heatmap tiles.
func NewHeatmapTileFactory(rmqConfig *Config, datasetConfigs ...[]byte) batch.TileFactoryCtor {
	setupConnection(rmqConfig, datasetConfigs...)

	return func() (batch.TileFactory, error) {
		Infof("New heatmap tile factory constructor request")
		return newHeatmapTile(rmqConfig), nil
	}
}

func newHeatmapTile(rmqConfig *Config) *HeatmapTile {
	ht := &HeatmapTile{}
	ht.tileType = "heatmap"
	ht.rmqConfig = rmqConfig
	ht.buildConfig = func() (map[string]interface{}, error) {
		return ht.getTileConfig()
	}
	ht.convert = func(coord *binning.TileCoord, input []byte) ([]byte, error) {
		return ht.convertTile(coord, input)
	}
	ht.buildDefault = func() ([]byte, error) {
		return ht.buildDefaultTile()
	}
	return ht
}

// Parse does the standard salt tile parsing of parameters - i.e., saving them for later
func (h *HeatmapTile) Parse(params map[string]interface{}) error {
	return h.TileData.Parse(params)
}

// parseHeatmapParams actually parses the provided JSON object, and
// populates the tile attributes.
func (h *HeatmapTile) parseHeatmapParams(params map[string]interface{}) error {
	valueField, ok := json.GetString(params, "valueField")
	if ok {
		h.valueField = valueField
	} else {
		h.valueField = ""
	}
	return h.Bivariate.Parse(params)
}

// GetTileConfig gets the configuration to send to Salt, so that it can
// construct the currently requested tile
func (h *HeatmapTile) getTileConfig() (map[string]interface{}, error) {
	err := h.parseHeatmapParams(*h.parameters)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})

	setProperty("type", "heatmap", result)
	setProperty("xField", h.XField, result)
	setProperty("yField", h.YField, result)
	if 0 < len(h.valueField) {
		setProperty("valueField", h.valueField, result)
	}
	setProperty("resolution", h.Resolution, result)
	// Bounds are ignored - salt needs the dataset bounds, not the tile bounds
	// in visualization space
	// setProperty("bounds.left",   h.Left, result)
	// setProperty("bounds.right",  h.Right, result)
	// setProperty("bounds.top",    h.Top, result)
	// setProperty("bounds.bottom", h.Bottom, result)

	return result, nil
}

func (h *HeatmapTile) convertTile(coord *binning.TileCoord, input []byte) ([]byte, error) {
	err := h.parseHeatmapParams(*h.parameters)
	if err != nil {
		return nil, err
	}

	res := h.Resolution
	numPoints := len(input) / 4
	if res*res != numPoints {
		return nil, fmt.Errorf("wrong number of points returned.  Expected %d, got %d", res*res, numPoints)
	}

	// Copy to output buffer, flipping the y
	output := make([]byte, len(input))
	stride := res * 4
	for y := 0; y < res; y++ {
		// Copy lines of floats in bulk, wholesale, rather than parsing and rewriting each number
		copy(output[(y+0)*stride:(y+1)*stride], input[(res-y-1)*stride:(res-y)*stride])
	}

	return output, nil
}

func (h *HeatmapTile) buildDefaultTile() ([]byte, error) {
	err := h.parseHeatmapParams(*h.parameters)
	if err != nil {
		return nil, err
	}

	bins := h.Resolution * h.Resolution
	bits := make([]byte, bins*4)
	return bits, nil
}
