package elastic

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
	// get client
	client, err := NewClient(m.Host, m.Port)
	if err != nil {
		return nil, err
	}
	// create search service
	search := client.Search().
		Index(uri).
		Size(0)

	// create root query
	q, err := m.CreateQuery(query)
	if err != nil {
		return nil, err
	}
	// add tiling query
	q.Must(m.Bivariate.GetQuery(coord))
	// set the query
	search.Query(q)

	// get aggs
	aggs := m.Bivariate.GetAggs(coord)
	// set the aggregation
	search.Aggregation("x", aggs["x"])

	// send query
	res, err := search.Do()
	if err != nil {
		return nil, err
	}

	// get bins
	bins, err := m.Bivariate.GetBins(&res.Aggregations)
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
		if bin != nil {
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
