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

func NewEquals() (prism.Query, error) {
	return &Equals{}, nil
}

// Apply adds the query to the tiling job.
func (q *Equals) Get() (elastic.Query, error) {
	return elastic.NewTermQuery(q.Field, q.Value), nil
}
