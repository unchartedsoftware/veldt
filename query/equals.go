package query

import (
	"fmt"

	"github.com/unchartedsoftware/prism"
	"github.com/unchartedsoftware/prism/util/json"
)

// Equals represents an equality query, checking if a field equals a provided
// value.
type Equals struct {
	Field string
	Value interface{}
}

// Parse parses the provided JSON object and populates the querys attributes.
func (q *Equals) Parse(params map[string]interface{}) error {
	field, ok := json.GetString(params, "field")
	if !ok {
		return nil, fmt.Errorf("`field` parameter missing from query params")
	}
	value, ok := json.Get(params, "value")
	if !ok {
		return nil, fmt.Errorf("`value` parameter missing from query params")
	}
	q.Field = field
	q.Value = value
	return nil
}

// Apply adds the query to the tiling job.
func (q *Equals) Apply(arg interface{}) error {
	return fmt.Errorf("Not implemented")
}
