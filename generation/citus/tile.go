package citus

import (
	"fmt"

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

		root.AddWhereClause(q)
	}

	return root, nil
}
