package citus

import (
	"math"

	"github.com/unchartedsoftware/prism"
	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/tile"
	"github.com/unchartedsoftware/prism/util/json"
)

type MacroTile struct {
	Bivariate
	Tile
	LOD int
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
	m.LOD = int(json.GetNumberDefault(params, 0, "lod"))
	return m.Bivariate.Parse(params)
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
	if m.LOD > 0 {
		return tile.EncodeLOD(points[0:numPoints*2], m.LOD), nil
	}
	return tile.Encode(points[0 : numPoints*2]), nil
}
