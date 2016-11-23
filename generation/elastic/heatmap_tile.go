package elastic

import (
	"encoding/binary"
	//"math"

	"github.com/unchartedsoftware/prism"
	"github.com/unchartedsoftware/prism/binning"
)

type Heatmap struct {
	Bivariate
	Tile
}

func NewHeatmapTile(host, port string) prism.TileCtor {
	return func() (prism.Tile, error) {
		h := &Heatmap{}
		h.Host = host
		h.Port = port
		return h, nil
	}
}

func (h *Heatmap) Parse(params map[string]interface{}) error {
	return h.Bivariate.Parse(params)
}

func (h *Heatmap) Create(uri string, coord *binning.TileCoord, query prism.Query) ([]byte, error) {
	// get client
	client, err := NewClient(h.Host, h.Port)
	if err != nil {
		return nil, err
	}
	// create search service
	search := client.Search().
		Index(uri).
		Size(0)

	// create root query
	q, err := h.CreateQuery(query)
	if err != nil {
		return nil, err
	}
	// add tiling query
	q.Must(h.Bivariate.GetQuery(coord))
	// set the query
	search.Query(q)

	// get aggs
	aggs := h.Bivariate.GetAggs(coord)
	// set the aggregation
	search.Aggregation("x", aggs["x"])

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
	bits := make([]byte, len(bins)*4)
	for i, bin := range bins {
		if bin != nil {
			binary.LittleEndian.PutUint32(
				bits[i*4:i*4+4],
				uint32(bin.DocCount))
		}
	}
	return bits[0:], nil
}
