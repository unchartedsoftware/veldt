package elastic

import (
	"encoding/binary"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
)

// HeatmapTile represents an elasticsearch implementation of the heatmap tile.
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
	bins, err := h.Bivariate.GetBins(&res.Aggregations)
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
	return bits, nil
}
