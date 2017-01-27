package citus

import (
	"fmt"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/query"
)

// Equals represents an = query.
type Equals struct {
	query.Equals
}

// NewEquals instantiates and returns a new query struct.
func NewEquals() (veldt.Query, error) {
	return &Equals{}, nil
}

// Get adds the parameters to the query and returns the string representation.
func (q *Equals) Get(query *Query) (string, error) {
	valueParam := query.AddParameter(q.Value)
	return fmt.Sprintf("%s = %s", q.Field, valueParam), nil
}
