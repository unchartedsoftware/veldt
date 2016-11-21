package citus

import (
	"encoding/binary"
	"math"

	"github.com/unchartedsoftware/prism"
	"github.com/unchartedsoftware/prism/binning"
)

type Macro struct {
	Bivariate
	Tile
}

func NewMacroTile(host, port string) prism.TileCtor {
	return func() (prism.Tile, error) {
		m := &Macro{}
		m.Host = host
		m.Port = port
		return m, nil
	}
}

func (m *Macro) Parse(params map[string]interface{}) error {
	return m.Bivariate.Parse(params)
}

func (m *Macro) Create(uri string, coord *binning.TileCoord, query prism.Query) ([]byte, error) {
	// get client
	client, err := NewClient(m.Host, m.Port)
	if err != nil {
		return nil, err
	}

	// create root query
    citusQuery, err := m.CreateQuery(query)
	if err != nil {
		return nil, err
	}

	// add tiling query
	citusQuery = m.Bivariate.GetQuery(coord, citusQuery)

	// add aggs
	citusQuery = m.Bivariate.GetAgg(coord, citusQuery)

	// set the aggregation
	//search.Aggregation("x", aggs["x"])

	citusQuery.AddTable(uri)
	citusQuery.AddField("CAST(COUNT(*) AS FLOAT) AS value")

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
	binSize := float64(256 / m.Resolution)
	halfSize := float64(binSize / 2)

	// convert to byte array
	bits := make([]byte, len(bins)*8)
	numPoints := 0
	for i, bin := range bins {
		if bin > 0 {
			x := float32(float64(i%m.Resolution)*binSize + halfSize)
			y := float32(math.Floor(float64(i/m.Resolution))*binSize + halfSize)
			binary.LittleEndian.PutUint32(
				bits[numPoints*8:numPoints*8+4],
				math.Float32bits(x))
			binary.LittleEndian.PutUint32(
				bits[numPoints*8+4:numPoints*8+8],
				math.Float32bits(y))
			numPoints++
		}
	}
	return bits[0 : numPoints*8], nil
}
