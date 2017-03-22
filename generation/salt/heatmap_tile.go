package salt

import (
	"github.com/unchartedsoftware/veldt"
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
func NewHeatmapTile (rmqConfig *Configuration, datasetConfigs ...[]byte) veldt.TileCtor {
	setupConnection(rmqConfig, datasetConfigs...)

	return func() (veldt.Tile, error) {
		saltInfof("New heatmap tile constructor request")
		return newHeatmapTile(rmqConfig), nil
	}
}

// NewHeatmapTileFactory instantiates and returns a factory for creating batched heatmap tiles.
func NewHeatmapTileFactory (rmqConfig *Configuration, datasetConfigs ...[]byte) batch.TileFactoryCtor {
	setupConnection(rmqConfig, datasetConfigs...)

	return func() (batch.TileFactory, error) {
		saltInfof("New heatmap tile factory constructor request")
		return newHeatmapTile(rmqConfig), nil
	}
}

func newHeatmapTile (rmqConfig *Configuration) *HeatmapTile {
	ht := &HeatmapTile{}
	ht.tileType = "heatmap"
	ht.rmqConfig = rmqConfig
	ht.buildConfig = func () (map[string]interface{}, error) {
		return ht.getTileConfiguration()
	}
	ht.convert = func (input []byte) ([]byte, error) {
		return ht.convertTile(input)
	}
	ht.buildDefault = func () ([]byte, error) {
		return ht.buildDefaultTile()
	}
	return ht
}

// Parse does the standard salt tile parsing of parameters - i.e., saving them for later
func (h *HeatmapTile) Parse (params map[string]interface{}) error {
	return h.TileData.Parse(params)
}

// parseHeatmapParameters actually parses the provided JSON object, and
// populates the tile attributes.
func (h *HeatmapTile) parseHeatmapParameters(params map[string]interface{}) error {
	valueField, ok := json.GetString(params, "valueField")
	if ok {
		h.valueField = valueField
	} else {
		h.valueField = ""
	}
	return h.Bivariate.Parse(params)
}

// GetTileConfiguration gets the configuration to send to Salt, so that it can
// construct the currently requested tile
func (h *HeatmapTile) getTileConfiguration () (map[string]interface{}, error) {
	err := h.parseHeatmapParameters(*h.parameters)
	if nil != err {
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

func (h *HeatmapTile) convertTile (input []byte) ([]byte, error) {
	return input, nil
}

func (h *HeatmapTile) buildDefaultTile () ([]byte, error) {
	err := h.parseHeatmapParameters(*h.parameters)
	if nil != err {
		return nil, err
	}

	bins := h.Resolution * h.Resolution
	bits := make([]byte, bins*4)
	return bits, nil
}
