package elastic

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/veldt"
)

// Tile represents an elasticsearch tile type.
type Tile struct {
	Host string
	Port string
}

// CreateQuery creates the elasticsearch query from the query struct.
func (t *Tile) CreateQuery(query veldt.Query) (*elastic.BoolQuery, error) {
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
