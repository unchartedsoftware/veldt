package citus

import (
	"encoding/binary"

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

	// create root query
	citusQuery, err := h.CreateQuery(query)
	if err != nil {
		return nil, err
	}
	citusQuery.From(uri)

	// add tiling query
	citusQuery = h.Bivariate.AddQuery(coord, citusQuery)

	// add aggs
	citusQuery = h.Bivariate.AddAgg(coord, citusQuery)

	// set the aggregation
	//search.Aggregation("x", aggs["x"].SubAggregation("y", aggs["y"]))

	citusQuery.Select("CAST(COUNT(*) AS FLOAT) AS value")
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
	bits := make([]byte, len(bins)*4)
	for i, bin := range bins {
		binary.LittleEndian.PutUint32(
			bits[i*4:i*4+4],
			uint32(bin))
	}
	return bits, nil
}
