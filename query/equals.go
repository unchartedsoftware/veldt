package query

import (
	"fmt"

	"github.com/unchartedsoftware/veldt/util/json"
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
		return fmt.Errorf("`field` parameter missing from query")
	}
	value, ok := json.Get(params, "value")
	if !ok {
		return fmt.Errorf("`value` parameter missing from query")
	}
	q.Field = field
	q.Value = value
	return nil
}
