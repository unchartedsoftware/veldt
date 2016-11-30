package elastic

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
	tileSize := 256.0
	binSize := tileSize / float64(m.Resolution)
	halfSize := float64(binSize / 2)

	// convert to byte array
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
	if m.LOD > 0 {
		return tile.EncodeLOD(points[0:numPoints*2], m.LOD), nil
	}
	return tile.Encode(points[0 : numPoints*2]), nil

	/*
		// bin width
		tileSize := 256.0
		binSize := tileSize / float64(m.Resolution)
		halfSize := float64(binSize / 2)

		// convert to byte array
		bits := make([]byte, len(bins)*8)
		numPoints := 0
		for i, bin := range bins {
			if bin != nil {
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
	*/
}
