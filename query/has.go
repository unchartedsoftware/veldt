package query

import (
	"fmt"
	"github.com/unchartedsoftware/prism/util/json"
	"strings"
)

// Has quiery represents a query checking if the field has one or more of the
// provided values. represents an elasticsearch terms query.
type Has struct {
	Field  string
	Values []interface{}
}

// NewHas instantiates and returns a has query object.
func NewHas(params map[string]interface{}) (Query, error) {
	field, ok := json.GetString(params, "field")
	if !ok {
		return nil, fmt.Errorf("`field` parameter missing from query params")
	}
	values, ok := json.GetArray(params, "values")
	if !ok {
		return nil, fmt.Errorf("`values` parameter missing from query params")
	}
	return &Has{
		Field:  field,
		Values: values,
	}, nil
}

// Apply adds the query to the tiling job.
func (q *Has) Apply(arg interface{}) error {
	return fmt.Errorf("Not implemented")
}
