package citus

import (
	"math"

	"github.com/unchartedsoftware/prism"
	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/tile"
)

type MacroTile struct {
	Tile
	Bivariate
	tile.Macro
}

func NewMacroTile(host, port string) prism.TileCtor {
	return func() (prism.Tile, error) {
		m := &MacroTile{}
		m.Host = host
		m.Port = port
		return m, nil
	}
}

func (m *MacroTile) Parse(params map[string]interface{}) error {
	err := m.Bivariate.Parse(params)
	if err != nil {
		return err
	}
	return m.Macro.Parse(params)
}

func (m *MacroTile) Create(uri string, coord *binning.TileCoord, query prism.Query) ([]byte, error) {
	// Initialize the tile processing.
	client, citusQuery, err := m.InitliazeTile(uri, query)

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
