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

// // NewRange instantiates and returns a range query object.
// func NewRange(params map[string]interface{}) (*Range, error) {
// 	field, ok := json.GetString(params, "field")
// 	if !ok {
// 		return nil, fmt.Errorf("Range `field` parameter missing from tiling param %v", params)
// 	}
// 	from, hasFrom := json.Get(params, "from")
// 	to, hasTo := json.Get(params, "to")
// 	if !hasFrom && !hasTo {
// 		return nil, fmt.Errorf("Range both `from` and `to` parameters missing from tiling param %v", params)
// 	}
// 	return &Range{
// 		Field: field,
// 		From:  from,
// 		To:    to,
// 	}, nil
// }

// Apply adds the query to the tiling job.
func (q *Range) Apply(query *elastic.Query) error {
	rang := elastic.NewRangeQuery(q.Field)
	if q.GTE != nil {
		rang.Gte(q.GTE)
	}
	if q.GT != nil {
		rang.Gt(q.GT)
	}
	if q.LTE != nil {
		rang.Lte(q.LTE)
	}
	if q.LT != nil {
		rang.Lte(q.LT)
	}
	query.Must(rang)
	return nil
}
