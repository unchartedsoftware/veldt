package citus

import (
	"encoding/binary"
	"math"

	"github.com/unchartedsoftware/prism"
	"github.com/unchartedsoftware/prism/binning"
)

type Heatmap struct {
	Bivariate
	Tile
}

func NewHeatmap(host, port string) prism.TileCtor {
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

	// create root query
	citusQuery, err := h.CreateQuery(query)
	if err != nil {
		return nil, err
	}

	// add tiling query
	citusQuery = h.Bivariate.GetQuery(coord, citusQuery)

	// add aggs
	citusQuery = h.Bivariate.GetAgg(coord, citusQuery)

	// set the aggregation
	//search.Aggregation("x", aggs["x"].SubAggregation("y", aggs["y"]))

	citusQuery.AddTable(uri)
	// send query
	res, err := client.Query(citusQuery.GetQuery(false), citusQuery.QueryArgs...)
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
