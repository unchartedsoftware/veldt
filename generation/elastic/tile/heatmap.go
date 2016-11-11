package tile

import (
	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/tile"
)

type Heatmap struct {
	tile.Bivariate
	Host string
	Port string
}

func NewHeatmap(host, port string) tile.Ctor {
	return func() {
		return &Heatmap{
			Host: host,
			Port: port,
		}
	}
}

func (h *Heatmap) Parse(params map[string]interface{}) error {
	return h.Bivariate.Parse(params)
}

func (h *Heatmap) CreateTile(uri string, coord tile.TileCoord, query tile.Query) ([]byte, error) {
	// get client
	client, err := NewClient(p.Host, p.Port)
	if err != nil {
		return nil, err
	}

	// create search service
	search := p.client.Search().
		Index(uri).
		Size(0)

	// create root query
	query := elastic.NewBoolQuery()

	// add tiling query
	query.Must(h.Bivariate.GetQuery(coord))

	// add filter query
	if query != nil {
		query.Must(query.GetQuery())
	}

	// add aggs
	aggs := h.Bivariate.GetAggs(coord)
	// set the aggregation
	search.Aggregation("x", aggs["x"].SubAggregation("y", aggs["y"]))

	// send query
	res, err := search.Do()
	if err != nil {
		return nil, err
	}

	// get bins
	bins, err := h.Bivariate.GetBins(res)
	if err != nil {
		return nil, err
	}

	// convert to byte array
	bits := make([]byte, len(bins)*8)
	for i, val := range bins {
		binary.LittleEndian.PutUint64(
			bits[i*8:i*8+8],
			math.Float64bits(val))
	}
	return bits[0:], nil
}
