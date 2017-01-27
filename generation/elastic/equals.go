package elastic

import (
	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism"
	"github.com/unchartedsoftware/prism/query"
)

// Equals represents an elasticsearch term query.
type Equals struct {
	query.Equals
}

// NewEquals instantiates and returns a new query struct.
func NewEquals() (prism.Query, error) {
	return &Equals{}, nil
}

// Get returns the appropriate elasticsearch query for the query.
func (q *Equals) Get() (elastic.Query, error) {
	return elastic.NewTermQuery(q.Field, q.Value), nil
}
