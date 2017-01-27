package citus

import (
	"math"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
	"github.com/unchartedsoftware/veldt/tile"
)

// MacroTile represents a citus implementation of the macro tile.
type MacroTile struct {
	Tile
	Bivariate
	tile.Macro
}

// NewMacroTile instantiates and returns a new tile struct.
func NewMacroTile(host, port string) veldt.TileCtor {
	return func() (veldt.Tile, error) {
		m := &MacroTile{}
		m.Host = host
		m.Port = port
		return m, nil
	}
}

// Parse parses the provided JSON object and populates the tiles attributes.
func (m *MacroTile) Parse(params map[string]interface{}) error {
	err := m.Bivariate.Parse(params)
	if err != nil {
		return err
	}
	return m.Macro.Parse(params)
}

// Create generates a tile from the provided URI, tile coordinate and query
// parameters.
func (m *MacroTile) Create(uri string, coord *binning.TileCoord, query veldt.Query) ([]byte, error) {
	// Initialize the tile processing.
	client, citusQuery, err := m.InitializeTile(uri, query)

	// add tiling query
	citusQuery = m.Bivariate.AddQuery(coord, citusQuery)

	// add aggs
	citusQuery = m.Bivariate.AddAggs(coord, citusQuery)

	citusQuery.Select("CAST(COUNT(*) AS FLOAT) AS value")

	// send query
	res, err := client.Query(citusQuery.GetQuery(false), citusQuery.QueryArgs...)
	if err != nil {
		return nil, err
	}

	// get bins
	bins, err := m.Bivariate.GetBins(res)
	if err != nil {
		return nil, err
	}

	// bin width
	binSize := binning.MaxTileResolution / float64(m.Resolution)
	halfSize := float64(binSize / 2)

	// convert to point array
	points := make([]float32, len(bins)*2)
	numPoints := 0
	for i, bin := range bins {
		if bin > 0 {
			x := float32(float64(i%m.Resolution)*binSize + halfSize)
			y := float32(math.Floor(float64(i/m.Resolution))*binSize + halfSize)
			points[numPoints*2] = x
			points[numPoints*2+1] = y
			numPoints++
		}
	}

	// encode the result
	return m.Macro.Encode(points[0 : numPoints*2])
}
