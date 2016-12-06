package elastic

import (
	"encoding/json"

	"github.com/unchartedsoftware/prism"
	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/tile"
	jsonutil "github.com/unchartedsoftware/prism/util/json"
)

type MicroTile struct {
	Bivariate
	Tile
	TopHits
	LOD int
}

func NewMicroTile(host, port string) prism.TileCtor {
	return func() (prism.Tile, error) {
		m := &MicroTile{}
		m.Host = host
		m.Port = port
		return m, nil
	}
}

func (m *MicroTile) Parse(params map[string]interface{}) error {
	m.LOD = int(jsonutil.GetNumberDefault(params, 0, "lod"))
	err := m.Bivariate.Parse(params)
	if err != nil {
		return err
	}
	err = m.TopHits.Parse(params)
	if err != nil {
		return err
	}
	// ensure that the x / y field are included
	xField := m.Bivariate.XField
	yField := m.Bivariate.YField
	includes := m.TopHits.IncludeFields
	if !existsIn(xField, includes) {
		includes = append(includes, xField)
	}
	if !existsIn(yField, includes) {
		includes = append(includes, yField)
	}
	m.TopHits.IncludeFields = includes
	return nil
}

func (m *MicroTile) Create(uri string, coord *binning.TileCoord, query prism.Query) ([]byte, error) {
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
	aggs := m.TopHits.GetAggs()
	// set the aggregation
	search.Aggregation("top-hits", aggs["top-hits"])

	// send query
	res, err := search.Do()
	if err != nil {
		return nil, err
	}

	// get top hits
	hits, err := m.TopHits.GetTopHits(&res.Aggregations)
	if err != nil {
		return nil, err
	}

	// bin width
	binSize := binning.MaxTileResolution / float64(m.Resolution)
	halfSize := float32(binSize / 2)

	// convert to point array
	points := make([]float32, len(hits)*2)
	for i, hit := range hits {
		ix, ok := hit[m.Bivariate.XField]
		if !ok {
			continue
		}
		iy, ok := hit[m.Bivariate.YField]
		if !ok {
			continue
		}
		x, ok := ix.(float64)
		if !ok {
			continue
		}
		y, ok := iy.(float64)
		if !ok {
			continue
		}
		points[i*2] = float32(x) + halfSize
		points[i*2+1] = float32(y) + halfSize
	}

	var buffer []byte
	if m.LOD > 0 {
		buffer = tile.EncodeLOD(points, m.LOD)
	} else {
		buffer = tile.Encode(points)
	}

	return json.Marshal(map[string]interface{}{
		"buffer": buffer,
		"hits":   hits,
	})
}

func existsIn(val string, arr []string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}
