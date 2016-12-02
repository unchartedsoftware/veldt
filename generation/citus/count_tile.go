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
	// get client
	client, err := NewClient(t.Host, t.Port)
	if err != nil {
		return nil, err
	}

	// create root query
	citusQuery, err := t.CreateQuery(query)
	if err != nil {
		return nil, err
	}
	citusQuery.From(uri)

	// add tiling query
	citusQuery = t.Bivariate.AddQuery(coord, citusQuery)

	citusQuery.Select("CAST(COUNT(*) AS FLOAT) AS value")
	// send query
	res, err := client.Query(citusQuery.GetQuery(false), citusQuery.QueryArgs...)
	if err != nil {
		return nil, err
	}

	var value float64
	err = res.Scan(&value)
	if err != nil {
		return nil, fmt.Errorf("Error parsing count: %v",
			err)
	}

	return []byte(fmt.Sprintf("{\"count\":%d}\n", value)), nil
}
