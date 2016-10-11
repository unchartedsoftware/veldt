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
	from, hasFrom := json.Get(params, "from")
	to, hasTo := json.Get(params, "to")
	if !hasFrom && !hasTo {
		return nil, fmt.Errorf("Range both `from` and `to` parameters missing from tiling param %v", params)
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
	query := elastic.NewRangeQuery(q.Field)
	if q.From != nil {
		query.Gte(q.From)
	}
	if q.To != nil {
		query.Lte(q.To)
	}
	return query
}
