package query

import (
	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/query"
)

// Has represents an elasticsearch terms query.
type Has struct {
	query.Has
}

// Apply adds the query to the tiling job.
func (q *Has) Get() (elastic.Query, error) {
	return elastic.NewTermsQuery(q.Field, q.Values...), nil
}
