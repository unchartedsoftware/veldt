package citus

import (
	"fmt"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/query"
)

// Exists checks for the existence of the field (not null).
type Exists struct {
	query.Exists
}

// NewExists instantiates and returns a new query struct.
func NewExists() (veldt.Query, error) {
	return &Exists{}, nil
}

// Get adds the parameters to the query and returns the string representation.
func (q *Exists) Get(query *Query) (string, error) {
	return fmt.Sprintf("%s IS NOT NULL", q.Field), nil
}
