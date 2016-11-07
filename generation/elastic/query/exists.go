package query

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/query"
)

// Exists represents an elasticsearch exists query.
type Exists struct {
	query.Exists
}

// Apply adds the query to the tiling job.
func (q *Exists) Apply(arg interface{}) error {
	query, ok := arg.(*elastic.BoolQuery)
	if !ok {
		return fmt.Errorf("`%v` is not of type *elastic.BoolQuery", arg)
	}
	query.Must(elastic.NewExistsQuery(q.Field))
	return nil
}
