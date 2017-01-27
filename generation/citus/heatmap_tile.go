package citus

import (
	"encoding/binary"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
)

// HeatmapTile represents a citus implementation of the heatmap tile.
type HeatmapTile struct {
	Bivariate
	Tile
}

// NewHeatmapTile instantiates and returns a new tile struct.
func NewHeatmapTile(host, port string) veldt.TileCtor {
	return func() (veldt.Tile, error) {
		h := &HeatmapTile{}
		h.Host = host
		h.Port = port
		return h, nil
	}
}

// Parse parses the provided JSON object and populates the tiles attributes.
func (h *HeatmapTile) Parse(params map[string]interface{}) error {
	return h.Bivariate.Parse(params)
}

// Create generates a tile from the provided URI, tile coordinate and query
// parameters.
func (h *HeatmapTile) Create(uri string, coord *binning.TileCoord, query veldt.Query) ([]byte, error) {
	// Initialize the tile processing.
	client, citusQuery, err := h.InitializeTile(uri, query)

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
