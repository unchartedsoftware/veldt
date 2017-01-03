package citus

import (
	"github.com/unchartedsoftware/prism"
	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/generation/common"
	jsonutil "github.com/unchartedsoftware/prism/util/json"
)

type MicroTile struct {
	Bivariate
	Tile
	TopHits
	LOD       int
	XIncluded bool
	YIncluded bool
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
	if !common.ExistsIn(xField, includes) {
		includes = append(includes, xField)
	} else {
		m.XIncluded = true
	}
	if !common.ExistsIn(yField, includes) {
		includes = append(includes, yField)
	} else {
		m.YIncluded = true
	}
	m.TopHits.IncludeFields = includes
	return nil
}

func (m *MicroTile) Create(uri string, coord *binning.TileCoord, query prism.Query) ([]byte, error) {
	// Initialize the tile processing.
	client, citusQuery, err := m.InitliazeTile(uri, query)

	// add tiling query
	citusQuery = m.Bivariate.AddQuery(coord, citusQuery)

	// get aggs
	citusQuery = m.TopHits.AddAggs(citusQuery)

	// send query
	res, err := client.Query(citusQuery.GetQuery(false), citusQuery.QueryArgs...)
	if err != nil {
		return nil, err
	}

	// get top hits
	hits, err := m.TopHits.GetTopHits(res)
	if err != nil {
		return nil, err
	}

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
		x := common.CastPixelResult(ix)
		y := common.CastPixelResult(iy)

		// convert to tile pixel coords
		tx := m.Bivariate.GetX(x)
		ty := m.Bivariate.GetY(y)
		// add to point array
		points[i*2] = common.ToFixed(float32(tx), 2)
		points[i*2+1] = common.ToFixed(float32(ty), 2)
		// remove fields if they weren't explicitly included
		if !m.XIncluded {
			delete(hit, m.Bivariate.XField)
		}
		if !m.YIncluded {
			delete(hit, m.Bivariate.YField)
		}
	}

	// check if there is any hit info to include at all
	if !m.XIncluded && !m.YIncluded && len(m.TopHits.IncludeFields) == 2 {
		// no point returning an array of empty hits
		hits = nil
	}

	return common.EncodeMicroTileResult(hits, points, m.LOD)
}
