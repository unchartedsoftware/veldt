package query

import (
	"fmt"

	"github.com/unchartedsoftware/veldt/util/json"
)

// Has quiery represents a query checking if the field has one or more of the
// provided values. represents an elasticsearch terms query.
type Has struct {
	Field  string
	Values []interface{}
}

// Parse parses the provided JSON object and populates the querys attributes.
func (q *Has) Parse(params map[string]interface{}) error {
	field, ok := json.GetString(params, "field")
	if !ok {
		return fmt.Errorf("`field` parameter missing from query")
	}
	values, ok := json.GetArray(params, "values")
	if !ok {
		return fmt.Errorf("`values` parameter missing from query")
	}
	q.Field = field
	q.Values = values
	return nil
}
