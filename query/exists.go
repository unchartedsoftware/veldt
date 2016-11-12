package query

import (
	"fmt"

	"github.com/unchartedsoftware/prism"
	"github.com/unchartedsoftware/prism/util/json"
)

// Exists represents an exists query checking if a field is not null.
type Exists struct {
	Field string
}

// Parse parses the provided JSON object and populates the querys attributes.
func (q *Exists) Parse(params map[string]interface{}) error {
	field, ok := json.GetString(params, "field")
	if !ok {
		return nil, fmt.Errorf("`field` parameter missing from query params")
	}
	q.Field = field
	return nil
}

// Apply adds the query to the tiling job.
func (q *Exists) Apply(arg interface{}) error {
	return fmt.Errorf("Not implemented")
}
