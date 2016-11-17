package citus

import (
	"fmt"

	"github.com/unchartedsoftware/prism"
	"github.com/unchartedsoftware/prism/query"
)

// Equals represents an = query.
type Equals struct {
	query.Equals
}

func NewEquals() (prism.Query, error) {
	return &Equals{}, nil
}

// Get adds the parameters to the query and returns the string representation.
func (q *Equals) Get(query *Query) (string, error) {
	valueParam := query.AddParameter(q.Value)
	//query.AddWhereClause(fmt.Sprintf("%s = %s", q.Field, valueParam))
	return fmt.Sprintf("%s = %s", q.Field, valueParam), nil
}
