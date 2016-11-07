package query

import (
	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/query"
)

// Equals represents an elasticsearch term query.
type Equals struct {
	query.Equals
}

// Apply adds the query to the tiling job.
func (q *Equals) Apply(arg interface{}) error {
	query, ok := arg.(*elastic.BoolQuery)
	if !ok {
		return fmt.Errorf("`%v` is not of type *elastic.BoolQuery", arg)
	}
	query.Must(elastic.NewTermQuery(q.Field, q.Value))
	return nil
}
