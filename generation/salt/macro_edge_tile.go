package salt

import (
	"encoding/binary"
	"math"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
	"github.com/unchartedsoftware/veldt/tile"
	"github.com/unchartedsoftware/veldt/generation/batch"
)

// MacroEdgeTile represents a salt implementation of the Edge tile
type MacroEdgeTile struct {
	TileData
	tile.TopHits
	tile.Edge
	tile.MacroEdge
}

// NewMacroEdgeTile instantiates and returns a new tile struct.
func NewMacroEdgeTile (rmqConfig *Configuration, datasetConfigs ...[]byte) veldt.TileCtor {
	setupConnection(rmqConfig, datasetConfigs...)

	return func() (veldt.Tile, error) {
		saltInfof("New macro edge tile constructor request")
		return newEdgeTile(rmqConfig), nil
	}
}

// NewMacroEdgeTileFactory instantiates and returns a factory for creating batched tiles.
func NewMacroEdgeTileFactory (rmqConfig *Configuration, datasetConfigs ...[]byte) batch.TileFactoryCtor {
	setupConnection(rmqConfig, datasetConfigs...)

	return func() (batch.TileFactory, error) {
		saltInfof("New macro edge tile factory constructor request")
		return newEdgeTile(rmqConfig), nil
	}
}

func newEdgeTile (rmqConfig *Configuration) *MacroEdgeTile {
	met := &MacroEdgeTile{}
	met.tileType = "macro-edge"
	met.rmqConfig = rmqConfig
	met.buildConfig = func () (map[string]interface{}, error) {
		return met.getTileConfiguration()
	}
	met.convert = func (coord *binning.TileCoord, input []byte) ([]byte, error) {
		return met.convertTile(coord, input)
	}
	met.buildDefault = func () ([]byte, error) {
		return met.buildDefaultTile()
	}
	return met
}

// Parse does the standard salt tile parsing of parameters - i.e., saving them for later
func (m *MacroEdgeTile) Parse (params map[string]interface{}) error {
	return m.TileData.Parse(params)
}

// parseEdgeParameters actually parses the provided JSON object, and
// populates the tile attributes.
func (m *MacroEdgeTile) parseEdgeParameters(params map[string]interface{}) error {
	err := m.Edge.Parse(params)
	if nil != err {
		return err
	}
	err = m.TopHits.Parse(params)
	if nil != err {
		return err
	}
	// parse includes
	m.TopHits.IncludeFields = m.MacroEdge.ParseIncludes(
		m.TopHits.IncludeFields,
		m.Edge.SrcXField,
		m.Edge.SrcYField,
		m.Edge.DstXField,
		m.Edge.DstYField)
	return m.MacroEdge.Parse(params)
}

// GetTileConfiguration gets the configuration to send to Salt, so that it can
// construct the currently requested tile
func (m *MacroEdgeTile) getTileConfiguration () (map[string]interface{}, error) {
	err := m.parseEdgeParameters(*m.parameters)
	if nil != err {
		return nil, err
	}

	result := make(map[string]interface{})
	setProperty("type", "macro-edge", result)
	setProperty("srcXField", m.SrcXField, result)
	setProperty("srcYField", m.SrcYField, result)
	setProperty("dstXField", m.DstXField, result)
	setProperty("dstYField", m.DstYField, result)
	setProperty("hitsCount", m.HitsCount, result)
	setProperty("lengthSorted", true, result)

	return result, nil
}

func (m *MacroEdgeTile) convertTile (coord *binning.TileCoord, input []byte) ([]byte, error) {
	err := m.parseEdgeParameters(*m.parameters)
	if nil != err {
		return nil, err
	}

	// Salt returns us absolute bin coordinates; we need to convert them to
	// values relative to the current tile
	tileSize := uint32(256)
	offsetX := float32(coord.X * tileSize)
	offsetY := float32(coord.Y * tileSize)
	
	// Macro tiles are returned to us as a series of floating-point numbers
	bLen := len(input)
	fLen := bLen / 4
	points := make([]float32, fLen)
	for i := 0; i < fLen; i++ {
		bits := binary.LittleEndian.Uint32(input[i*4:i*4+4])
		base := math.Float32frombits(bits)
		if 0 == (i % 2) {
			// X coordinate
			points[i] = base - offsetX
		} else {
			// Y coordinate
			points[i] = base - offsetY
		}
	}	

	return m.MacroEdge.Encode(points)
}

func (m *MacroEdgeTile) buildDefaultTile () ([]byte, error) {
	return m.MacroEdge.Encode(make([]float32, 0))
}
