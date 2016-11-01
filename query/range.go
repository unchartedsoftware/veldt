package query

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

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
func NewRange(filters map[string]interface{}) (*Range, error) {
}

// GetHash returns a string hash of the query.
func (q *Range) GetHash() string {
	values := make([]interface{})
	if q.GT != nil {
		values := append(values, q.GT)
	}
	if q.GTE != nil {
		values := append(values, q.GTE)
	}
	if q.LT != nil {
		values := append(values, q.LT)
	}
	if q.LTE != nil {
		values := append(values, q.LTE)
	}
	return fmt.Sprintf("%s:%s",
		q.Field,
		strings.Join(hashes, ":"))
}
