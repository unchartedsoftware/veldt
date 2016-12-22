package citus

import (
	"fmt"

	"github.com/jackc/pgx"
	"github.com/unchartedsoftware/prism"
)

type Tile struct {
	Host string
	Port string
}

func (t *Tile) CreateQuery(query prism.Query) (*Query, error) {
	// create root query
	root, err := NewQuery()
	if err != nil {
		return nil, err
	}

	// add filter query
	if query != nil {
		// type assert
		citusQuery, ok := query.(QueryString)
		if !ok {
			return nil, fmt.Errorf("query is not citus.Query")
		}
		// get underlying query
		q, err := citusQuery.Get(root)
		if err != nil {
			return nil, err
		}

		root.Where(q)
	}

	return root, nil
}

func (t *Tile) InitliazeTile(uri string, query prism.Query) (*pgx.ConnPool, *Query, error) {
	// get client
	client, err := NewClient(t.Host, t.Port)
	if err != nil {
		return nil, nil, err
	}

	// create root query
	citusQuery, err := t.CreateQuery(query)
	if err != nil {
		return nil, nil, err
	}
	citusQuery.From(uri)

	return client, citusQuery, nil
}
