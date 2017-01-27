package query

import (
	"fmt"

	"github.com/unchartedsoftware/veldt/util/json"
)

// Exists represents an exists query checking if a field is not null.
type Exists struct {
	Field string
}

// Parse parses the provided JSON object and populates the querys attributes.
func (q *Exists) Parse(params map[string]interface{}) error {
	field, ok := json.GetString(params, "field")
	if !ok {
		return fmt.Errorf("`field` parameter missing from query")
	}
	q.Field = field
	return nil
}
