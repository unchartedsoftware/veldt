package query

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/query"
)

// Range represents an elasticsearch range query.
type Range struct {
	query.Range
}

// Apply adds the query to the tiling job.
func (q *Range) Apply(arg interface{}) error {
	query, ok := arg.(*elastic.BoolQuery)
	if !ok {
		return fmt.Errorf("`%v` is not of type *elastic.BoolQuery", arg)
	}
	rang := elastic.NewRangeQuery(q.Field)
	if q.GTE != nil {
		rang.Gte(q.GTE)
	}
	if q.GT != nil {
		rang.Gt(q.GT)
	}
	if q.LTE != nil {
		rang.Lte(q.LTE)
	}
	if q.LT != nil {
		rang.Lte(q.LT)
	}
	query.Must(rang)
	return nil
}
