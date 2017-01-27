package citus

import (
	"fmt"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/query"
)

// Has represents an citus query on an array.
type Has struct {
	query.Has
}

// NewHas instantiates and returns a new query struct.
func NewHas() (veldt.Query, error) {
	return &Has{}, nil
}

// Get adds the parameters to the query and returns the string representation.
func (q *Has) Get(query *Query) (string, error) {
	// Check that the array contains the values.
	// Use the column && ARRAY[value1, value2] notation.
	clause := ""

	//Generate the array values.
	for _, value := range q.Values {
		valueParam := query.AddParameter(value)
		clause = clause + fmt.Sprintf(", %s", valueParam)
	}

	//Remove the leading ", " from the array contents.
	clause = fmt.Sprintf("%s && ARRAY[%s]", q.Field, clause[2:])
	return clause, nil
}
