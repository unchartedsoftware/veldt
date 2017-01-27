package elastic

import (
	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/query"
)

// Range represents an elasticsearch range query.
type Range struct {
	query.Range
}

// NewRange instantiates and returns a new query struct.
func NewRange() (veldt.Query, error) {
	return &Range{}, nil
}

// Get returns the appropriate elasticsearch query for the query.
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
