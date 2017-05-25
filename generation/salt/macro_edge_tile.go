package salt

import (
	"encoding/binary"
	"math"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
	"github.com/unchartedsoftware/veldt/generation/batch"
	"github.com/unchartedsoftware/veldt/tile"
)

// MacroEdgeTile represents a salt implementation of the Edge tile
type MacroEdgeTile struct {
	TileData
	tile.TopHits
	tile.Edge
	tile.MacroEdge
}

// NewMacroEdgeTile instantiates and returns a new tile struct.
func NewMacroEdgeTile(rmqConfig *Config, datasetConfigs ...[]byte) veldt.TileCtor {
	setupConnection(rmqConfig, datasetConfigs...)

	return func() (veldt.Tile, error) {
		Infof("New macro edge tile constructor request")
		return newEdgeTile(rmqConfig), nil
	}
}

// NewMacroEdgeTileFactory instantiates and returns a factory for creating batched tiles.
func NewMacroEdgeTileFactory(rmqConfig *Config, datasetConfigs ...[]byte) batch.TileFactoryCtor {
	setupConnection(rmqConfig, datasetConfigs...)

	return func() (batch.TileFactory, error) {
		Infof("New macro edge tile factory constructor request")
		return newEdgeTile(rmqConfig), nil
	}
}

func newEdgeTile(rmqConfig *Config) *MacroEdgeTile {
	met := &MacroEdgeTile{}
	met.tileType = "macro-edge"
	met.rmqConfig = rmqConfig
	met.buildConfig = func() (map[string]interface{}, error) {
		return met.getTileConfig()
	}
	met.convert = func(coord *binning.TileCoord, input []byte) ([]byte, error) {
		return met.convertTile(coord, input)
	}
	met.buildDefault = func() ([]byte, error) {
		return met.buildDefaultTile()
	}
	return met
}

// Parse does the standard salt tile parsing of parameters - i.e., saving them for later
func (m *MacroEdgeTile) Parse(params map[string]interface{}) error {
	return m.TileData.Parse(params)
}

// parseEdgeParams actually parses the provided JSON object, and
// populates the tile attributes.
func (m *MacroEdgeTile) parseEdgeParams(params map[string]interface{}) error {
	err := m.Edge.Parse(params)
	if err != nil {
		return err
	}
	err = m.TopHits.Parse(params)
	if err != nil {
		return err
	}
	// parse includes
	m.TopHits.IncludeFields = m.MacroEdge.ParseIncludes(
		m.TopHits.IncludeFields,
		m.Edge.SrcXField,
		m.Edge.SrcYField,
		m.Edge.DstXField,
		m.Edge.DstYField,
		m.Edge.WeightField)
	return m.MacroEdge.Parse(params)
}

// GetTileConfig gets the configuration to send to Salt, so that it can
// construct the currently requested tile
func (m *MacroEdgeTile) getTileConfig() (map[string]interface{}, error) {
	err := m.parseEdgeParams(*m.parameters)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	setProperty("type", "macro-edge", result)
	setProperty("srcXField", m.SrcXField, result)
	setProperty("srcYField", m.SrcYField, result)
	setProperty("dstXField", m.DstXField, result)
	setProperty("dstYField", m.DstYField, result)
	setProperty("weightField", m.WeightField, result)
	setProperty("hitsCount", m.HitsCount, result)
	setProperty("lengthSorted", true, result)

	return result, nil
}

func (m *MacroEdgeTile) convertTile(coord *binning.TileCoord, input []byte) ([]byte, error) {
	err := m.parseEdgeParams(*m.parameters)
	if err != nil {
		return nil, err
	}

	// Salt returns us absolute bin coordinates; we need to convert them to
	// values relative to the current tile
	tileSize := uint32(binning.MaxTileResolution)
	offsetX := float32(coord.X * tileSize)
	offsetY := float32(coord.Y * tileSize)

	// Macro tiles are returned to us as a series of floating-point numbers
	bLen := len(input)
	fLen := bLen / 4 // number of floating-pont numbers
	edges := fLen / 5 // number of edges
	points := make([]float32, edges * 6)
	for i := 0; i < edges; i++ {
		srcXBits :=   binary.LittleEndian.Uint32(input[i*20 +  0 : i*20 +  4])
		srcYBits :=   binary.littleEndian.Uint32(input[i*20 +  4 : i*20 +  8])
		dstXBits :=   binary.LittleEndian.Uint32(input[i*20 +  8 : i*20 + 12])
		dstYBits :=   binary.littleEndian.Uint32(input[i*20 + 12 : i*20 + 16])
		weightBits := binary.littleEndian.Uint32(input[i*20 + 16 : i*20 + 20])

		srcX := math.Float32frombits(srcXBits) - offsetX
		srcY := math.Float32frombits(srcYBits) - offsetY
		dstX := math.Float32frombits(dstXBits) - offsetX
		dstY := math.Float32frombits(dstYBits) - offsetY
		weight := math.Float32frombits(weightBits)

		points[i*6 + 0] = srcX
		points[i*6 + 1] = srcY
		points[i*6 + 2] = weight
		points[i*6 + 3] = dstX
		points[i*6 + 4] = dstY
		points[i*6 + 5] = weight
	}

	return m.MacroEdge.Encode(points)
}

func (m *MacroEdgeTile) buildDefaultTile() ([]byte, error) {
	return m.MacroEdge.Encode(make([]float32, 0))
}
