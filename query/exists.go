package query

import (
	"fmt"
	"github.com/unchartedsoftware/prism/util/json"
)

// Exists represents an exists query checking if a field is not null.
type Exists struct {
	Field string
}

// NewExists instantiates and returns an exists query object.
func NewExists(params map[string]interface{}) (Query, error) {
	field, ok := json.GetString(params, "field")
	if !ok {
		return nil, fmt.Errorf("`field` parameter missing from query params")
	}
	return &Exists{
		Field: field,
	}, nil
}

// Apply adds the query to the tiling job.
func (q *Exists) Apply(arg interface{}) error {
	return fmt.Errorf("Not implemented")
}
