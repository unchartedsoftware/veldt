package elastic

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism"
)

type Tile struct {
	Host string
	Port string
}

func (t *Tile) CreateQuery(query prism.Query) (*elastic.BoolQuery, error) {
	// create root query
	root := elastic.NewBoolQuery()
	// add filter query
	if query != nil {
		// type assert
		esquery, ok := query.(Query)
		if !ok {
			return nil, fmt.Errorf("query is not elastic.Query")
		}
		// get underlying query
		q, err := esquery.Get()
		if err != nil {
			return nil, err
		}
		root.Must(q)
	}
	return root, nil
}
