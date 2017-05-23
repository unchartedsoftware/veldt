package elastic

import (
	"fmt"
	"strings"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/veldt"
)

// Tile represents an elasticsearch tile type.
type Tile struct {
	Host string
	Port string
}

// CreateSearchService creates the elasticsearch search service from the provided uri. 
func (t *Tile) CreateSearchService(uri string) (*elastic.SearchService, error) { 
  // get client 
  client, err := NewClient(t.Host, t.Port) 
  if err != nil { 
    return nil, err 
  } 
  split := strings.Split(uri, "/") 
  if len(split) < 2 { 
    index := split[0] 
    return client.Search(). 
      Index(index). 
      Size(0), nil 
  } 
  index := split[0] 
  typ := split[1] 
  return client.Search(). 
    Index(index). 
    Type(typ). 
    Size(0), nil 
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
