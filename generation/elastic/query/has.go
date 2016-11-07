package query

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/query"
)

// Has represents an elasticsearch terms query.
type Has struct {
	query.Has
}

// Apply adds the query to the tiling job.
func (q *Has) Apply(arg interface{}) error {
	query, ok := arg.(*elastic.BoolQuery)
	if !ok {
		return fmt.Errorf("`%v` is not of type *elastic.BoolQuery", arg)
	}
	query.Must(elastic.NewTermsQuery(q.Field, q.Values...))
	return nil
}
