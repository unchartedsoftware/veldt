package citus

import (
    "fmt"

	"github.com/unchartedsoftware/prism/query"
)

// Has represents an cituss query on an array.
type Has struct {
	query.Has
}

// Get adds the parameters to the query and returns the string representation.
func (q *Has) Get(query *Query) (string, error) {
    // Check that the array contains the values.
    // Use the column @> ARRAY[value1, value2] notation.
    clause := ""

    //Generate the array values.
    for _, value := range q.Values {
        valueParam := query.AddParameter(value)
        clause = clause + fmt.Sprintf(", %s", valueParam)
    }

    //Remove the leading ", " from the array contents.
    clause = fmt.Sprintf("%s @> ARRAY[%s]", q.Field, clause[2:])
	return clause, nil
}
