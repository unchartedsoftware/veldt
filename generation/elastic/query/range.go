package query

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/util/json"
)

// Range represents an elasticsearch range query.
type Range struct {
	Field string
	From  interface{}
	To    interface{}
}

// NewRange instantiates and returns a range query object.
func NewRange(params map[string]interface{}) (*Range, error) {
	field, ok := json.GetString(params, "field")
	if !ok {
		return nil, fmt.Errorf("Range `field` parameter missing from tiling param %v", params)
	}
	from, ok := json.Get(params, "from")
	if !ok {
		return nil, fmt.Errorf("Range `from` parameter missing from tiling param %v", params)
	}
	to, ok := json.Get(params, "to")
	if !ok {
		return nil, fmt.Errorf("Range `to` parameter missing from tiling param %v", params)
	}
	return &Range{
		Field: field,
		From:  from,
		To:    to,
	}, nil
}

// GetHash returns a string hash of the query.
func (q *Range) GetHash() string {
	return fmt.Sprintf("%s:%f:%f",
		q.Field,
		q.From,
		q.To)
}

// GetQuery returns the elastic query object.
func (q *Range) GetQuery() elastic.Query {
	return elastic.NewRangeQuery(q.Field).
		Gte(q.From).
		Lte(q.To)
}
