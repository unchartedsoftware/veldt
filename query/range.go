package query

import (
	"fmt"

	"github.com/unchartedsoftware/veldt/util/json"
)

// Range represents a range query, check that the values are within the defined
// range.
type Range struct {
	Field string
	GT    interface{}
	GTE   interface{}
	LT    interface{}
	LTE   interface{}
}

// Parse parses the provided JSON object and populates the querys attributes.
func (q *Range) Parse(params map[string]interface{}) error {
	field, ok := json.GetString(params, "field")
	if !ok {
		return fmt.Errorf("`field` parameter missing from query")
	}
	gte, gteOk := json.Get(params, "gte")
	gt, gtOk := json.Get(params, "gt")
	lte, lteOk := json.Get(params, "lte")
	lt, ltOk := json.Get(params, "lt")
	if !gteOk && !gtOk && !lteOk && !ltOk {
		return fmt.Errorf("range has no valid range parameters")
	}
	if gteOk && gtOk {
		return fmt.Errorf("both `gte` and `gt` have been provided, only one upper bound may be provided")
	}
	if lteOk && ltOk {
		return fmt.Errorf("both `lte` and `lt` have been provided, only one lower bound may be provided")
	}
	q.Field = field
	q.GTE = gte
	q.GT = gt
	q.LTE = lte
	q.LT = lt
	return nil
}
