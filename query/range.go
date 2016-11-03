package query

import (
	"fmt"
	"strings"
	"github.com/unchartedsoftware/prism/util/json"
)

// Range represents a range query, check that the values are within the defined
// range.
type Range struct {
	Field string
	GT interface{}
	GTE interface{}
	LT interface{}
	LTE interface{}
}

// NewRange instantiates and returns a range query object.
func NewRange(params map[string]interface{}) (Query, error) {
	field, ok := json.GetString(params, "field")
	if !ok {
		return nil, fmt.Errorf("Range `field` parameter missing from tiling param %v", params)
	}
    gte, gteOk := json.Get(params, "gte")
    gt, gtOk := json.Get(params, "gt")
    lte, lteOk := json.Get(params, "lte")
    lt, ltOk := json.Get(params, "lt")
    if !gteOk && !gtOk && !lteOk && !ltOk {
		return nil, fmt.Errorf("Range has no valid range parameters in tiling param %v", params)
	}
	return &Range{
		Field: field,
        GTE: gte,
        GT: gt,
        LTE: lte,
        LT: lt,
	}, nil
}

// Apply adds the query to the tiling job.
func (q *Range) Apply(arg interface{}) error {
	return fmt.Errorf("Range has not been implemented")
}

// GetHash returns a string hash of the query.
func (q *Range) GetHash() string {
	values := make([]string, 0)
	if q.GT != nil {
		values = append(values, fmt.Sprintf("%v", q.GT))
	}
	if q.GTE != nil {
		values = append(values, fmt.Sprintf("%v", q.GTE))
	}
	if q.LT != nil {
		values = append(values, fmt.Sprintf("%v", q.LT))
	}
	if q.LTE != nil {
		values = append(values, fmt.Sprintf("%v", q.LTE))
	}
	return fmt.Sprintf("%s:%s",
		q.Field,
		strings.Join(values, ":"))
}
