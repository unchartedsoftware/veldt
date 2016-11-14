package query

import (
	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/query"
)

// Exists represents an elasticsearch exists query.
type Exists struct {
	query.Exists
}

// Apply adds the query to the tiling job.
func (q *Exists) Get() (elastic.Query, error) {
	return elastic.NewExistsQuery(q.Field), nil
}
