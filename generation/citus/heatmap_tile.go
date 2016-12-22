package citus

import (
	"encoding/binary"

	"github.com/unchartedsoftware/prism"
	"github.com/unchartedsoftware/prism/binning"
)

type HeatmapTile struct {
	Bivariate
	Tile
}

func NewHeatmapTile(host, port string) prism.TileCtor {
	return func() (prism.Tile, error) {
		h := &HeatmapTile{}
		h.Host = host
		h.Port = port
		return h, nil
	}
}

func (h *HeatmapTile) Parse(params map[string]interface{}) error {
	return h.Bivariate.Parse(params)
}

func (h *HeatmapTile) Create(uri string, coord *binning.TileCoord, query prism.Query) ([]byte, error) {
	// Initialize the tile processing.
	client, citusQuery, err := h.InitliazeTile(uri, query)

	// add tiling query
	citusQuery = h.Bivariate.AddQuery(coord, citusQuery)

	// add aggs
	citusQuery = h.Bivariate.AddAggs(coord, citusQuery)

	//May support AVG (& others) in the future. May as well make it a float for now.
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
