package citus

import (
	"fmt"

	"github.com/unchartedsoftware/prism"
	"github.com/unchartedsoftware/prism/binning"
)

type Count struct {
	Bivariate
	Tile
}

func NewCountTile(host, port string) prism.TileCtor {
	return func() (prism.Tile, error) {
		t := &Count{}
		t.Host = host
		t.Port = port
		return t, nil
	}
}

func (t *Count) Create(uri string, coord *binning.TileCoord, query prism.Query) ([]byte, error) {
	// Initialize the tile processing.
	client, citusQuery, err := t.InitliazeTile(uri, query)

	// add tiling query
	citusQuery = t.Bivariate.AddQuery(coord, citusQuery)

	citusQuery.Select("CAST(COUNT(*) AS FLOAT) AS value")
	// send query
	res, err := client.Query(citusQuery.GetQuery(false), citusQuery.QueryArgs...)
	if err != nil {
		return nil, err
	}

	value := float64(0.0)
	for res.Next() {
		err = res.Scan(&value)
		if err != nil {
			return nil, fmt.Errorf("Error parsing count: %v",
				err)
		}
	}

	return []byte(fmt.Sprintf(`{"count":%d}`, uint64(value))), nil
}
