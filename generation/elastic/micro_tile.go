package elastic

import (
	"encoding/json"
	"math"
	"sort"

	"github.com/unchartedsoftware/prism"
	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/tile"
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
	if !existsIn(xField, includes) {
		includes = append(includes, xField)
	} else {
		m.XIncluded = true
	}
	if !existsIn(yField, includes) {
		includes = append(includes, yField)
	} else {
		m.YIncluded = true
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

	// convert to point array and hits array
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
		// convert to tile pixel coords
		tx := m.Bivariate.GetX(x)
		ty := m.Bivariate.GetY(y)
		// add to point array
		points[i*2] = toFixed(float32(tx), 2)
		points[i*2+1] = toFixed(float32(ty), 2)
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

	if m.LOD > 0 {
		// sort hits by morton code so they align
		sortHitsArray(hits, points)
		// sort points and get offsets
		sortedPoints, offsets := tile.LOD(points, m.LOD)
		return json.Marshal(map[string]interface{}{
			"points":  sortedPoints,
			"offsets": offsets,
			"hits":    hits,
		})
	}

	return json.Marshal(map[string]interface{}{
		"points": points,
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

func toFixed(num float32, precision int) float32 {
	output := math.Pow(10, float64(precision))
	return float32(math.Floor(float64(num)*output+0.5)) / float32(output)
}

func sortHitsArray(hits []map[string]interface{}, points []float32) {
	if hits == nil {
		return
	}
	// sort hits by morton code so they align
	hitsArr := make(hitsArray, len(hits))
	for i, hit := range hits {
		// add to hits array
		hitsArr[i] = &hitWrapper{
			x:    points[i*2],
			y:    points[i*2],
			data: hit,
		}
	}
	sort.Sort(hitsArr)
	// copy back into same arr
	for i, hit := range hitsArr {
		hits[i] = hit.data
	}
}

type hitWrapper struct {
	x    float32
	y    float32
	data map[string]interface{}
}

type hitsArray []*hitWrapper

func (h hitsArray) Len() int {
	return len(h)
}
func (h hitsArray) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
func (h hitsArray) Less(i, j int) bool {
	return tile.Morton(h[i].x, h[i].y) < tile.Morton(h[j].x, h[j].y)
}
