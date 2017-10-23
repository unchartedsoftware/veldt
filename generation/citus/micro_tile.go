package citus

import (
	"fmt"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
	"github.com/unchartedsoftware/veldt/tile"
)

// MicroTile represents a citus implementation of the micro tile.
type MicroTile struct {
	Bivariate
	Tile
	TopHits
	tile.Micro
}

// NewMicroTile instantiates and returns a new tile struct.
func NewMicroTile(cfg *Config) veldt.TileCtor {
	return func() (veldt.Tile, error) {
		m := &MicroTile{}
		m.Config = cfg
		return m, nil
	}
}

// Parse parses the provided JSON object and populates the tiles attributes.
func (m *MicroTile) Parse(params map[string]interface{}) error {
	err := m.Bivariate.Parse(params)
	if err != nil {
		return err
	}
	err = m.TopHits.Parse(params)
	if err != nil {
		return err
	}
	err = m.Micro.Parse(params)
	if err != nil {
		return err
	}
	// parse includes
	m.TopHits.IncludeFields = m.Micro.ParseIncludes(
		m.TopHits.IncludeFields,
		m.Bivariate.XField,
		m.Bivariate.YField)
	return nil
}

// Create generates a tile from the provided URI, tile coordinate and query
// parameters.
func (m *MicroTile) Create(uri string, coord *binning.TileCoord, query veldt.Query) ([]byte, error) {
	// Initialize the tile processing.
	client, citusQuery, err := m.InitializeTile(uri, query)
	if err != nil {
		return nil, err
	}

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
		// get hit x/y in tile coords
		x, y, ok := m.Bivariate.GetXY(coord, hit)
		if !ok {
			return nil, fmt.Errorf("could not parse position from hit: %v", hit)
		}
		// add to point array
		points[i*2] = float32(x)
		points[i*2+1] = float32(y)
	}

	// encode and return results
	return m.Micro.Encode(hits, points)
}
