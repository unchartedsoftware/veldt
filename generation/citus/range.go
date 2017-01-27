package citus

import (
	"fmt"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/query"
)

// Range represents a citus range query.
type Range struct {
	query.Range
}

// NewRange instantiates and returns a new query struct.
func NewRange() (veldt.Query, error) {
	return &Range{}, nil
}

// Get adds the parameters to the query and returns the string representation.
func (q *Range) Get(query *Query) (string, error) {
	clause := ""
	valueParam := ""
	if q.GTE != nil {
		valueParam = query.AddParameter(q.GTE)
		clause = clause + fmt.Sprintf(" AND %s >= %v", q.Field, valueParam)
	}
	if q.GT != nil {
		valueParam = query.AddParameter(q.GT)
		clause = clause + fmt.Sprintf(" AND %s > %v", q.Field, valueParam)
	}
	if q.LTE != nil {
		valueParam = query.AddParameter(q.LTE)
		clause = clause + fmt.Sprintf(" AND %s <= %v", q.Field, valueParam)
	}
	if q.LT != nil {
		valueParam = query.AddParameter(q.LT)
		clause = clause + fmt.Sprintf(" AND %s < %v", q.Field, valueParam)
	}
	//Remove leading " AND "
	return clause[5:], nil
}
