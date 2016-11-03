package query

import (
	"fmt"
	"github.com/unchartedsoftware/prism/util/json"
)

// Equals represents an equality query, checking if a field equals a provided
// value.
type Equals struct {
	Field string
	Value interface{}
}

// NewEquals instantiates and returns an equals query object.
func NewEquals(params map[string]interface{}) (Query, error) {
	field, ok := json.GetString(params, "field")
	if !ok {
		return nil, fmt.Errorf("Equals `field` parameter missing from tiling param %v", params)
	}
	value, ok := json.Get(params, "value")
	if !ok {
		return nil, fmt.Errorf("Equals `value` parameter missing from tiling param %v", params)
	}
	return &Equals{
		Field: field,
		Value: value,
	}, nil
}

// Apply adds the query to the tiling job.
func (q *Equals) Apply(arg interface{}) error {
	return fmt.Errorf("Equals has not been implemented")
}

// GetHash returns a string hash of the query.
func (q *Equals) GetHash() string {
	return fmt.Sprintf("%s:%v",
		q.Field,
		q.Value)
}
