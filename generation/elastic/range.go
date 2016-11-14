package elastic

import (
	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/query"
)

// Range represents an elasticsearch range query.
type Range struct {
	query.Range
}

// Apply adds the query to the tiling job.
func (q *Range) Get() (elastic.Query, error) {
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
	return rang, nil
}
