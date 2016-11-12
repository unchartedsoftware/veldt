package query

import (
	"fmt"
	"strings"

	"github.com/unchartedsoftware/prism"
	"github.com/unchartedsoftware/prism/util/json"
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
		return nil, fmt.Errorf("`field` parameter missing from query params")
	}
	gte, gteOk := json.Get(params, "gte")
	gt, gtOk := json.Get(params, "gt")
	lte, lteOk := json.Get(params, "lte")
	lt, ltOk := json.Get(params, "lt")
	if !gteOk && !gtOk && !lteOk && !ltOk {
		return nil, fmt.Errorf("Range has no valid range parameters")
	}
	q.Field = field
	q.GTE = gte
	q.GT = gt
	q.LTE = lte
	q.LT = lt
	return nil
}

// Apply adds the query to the tiling job.
func (q *Range) Apply(arg interface{}) error {
	return fmt.Errorf("Not implemented")
}
