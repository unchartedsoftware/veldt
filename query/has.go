package query

import (
	"fmt"
	"strings"
	"github.com/unchartedsoftware/prism/util/json"
)

// Has quiery represents a query checking if the field has one or more of the
// provided values. represents an elasticsearch terms query.
type Has struct {
	Field string
	Values []interface{}
}

// NewHas instantiates and returns a has query object.
func NewHas(params map[string]interface{}) (Query, error) {
	field, ok := json.GetString(params, "field")
	if !ok {
		return nil, fmt.Errorf("Has `field` parameter missing from tiling param %v", params)
	}
	values, ok := json.GetArray(params, "values")
	if !ok {
		return nil, fmt.Errorf("Has `values` parameter missing from tiling param %v", params)
	}
	return &Has{
		Field: field,
		Values: values,
	}, nil
}

// Apply adds the query to the tiling job.
func (q *Has) Apply(arg interface{}) error {
	return fmt.Errorf("Has has not been implemented")
}

// GetHash returns a string hash of the query.
func (q *Has) GetHash() string {
	hashes := make([]string, len(q.Values))
	for i, val := range q.Values {
		hashes[i] = fmt.Sprintf("%v", val)
	}
	return fmt.Sprintf("%s:%s",
		q.Field,
		strings.Join(hashes, ":"))
}
