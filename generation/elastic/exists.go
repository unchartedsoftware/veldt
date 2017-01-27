package elastic

import (
	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism"
	"github.com/unchartedsoftware/prism/query"
)

// Exists represents an elasticsearch exists query.
type Exists struct {
	query.Exists
}

// NewExists instantiates and returns a new query struct.
func NewExists() (prism.Query, error) {
	return &Exists{}, nil
}

// Get returns the appropriate elasticsearch query for the query.
func (q *Exists) Get() (elastic.Query, error) {
	return elastic.NewExistsQuery(q.Field), nil
}
