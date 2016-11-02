package query

import (
	"fmt"
	"sort"
	"strings"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/util/json"
)

// Has represents an elasticsearch terms query.
type Has struct {
	query.Has
}

// // NewHas instantiates and returns a terms query object.
// func NewHas(params map[string]interface{}) (*Has, error) {
// 	field, ok := json.GetString(params, "field")
// 	if !ok {
// 		return nil, fmt.Errorf("Has `field` parameter missing from tiling param %v", params)
// 	}
// 	terms, ok := json.GetStringArray(params, "terms")
// 	if !ok {
// 		return nil, fmt.Errorf("Has `terms` parameter missing from tiling param %v", params)
// 	}
// 	// sort terms
// 	sort.Strings(terms)
// 	return &Has{
// 		Field: field,
// 		Has: terms,
// 	}, nil
// }

// Apply adds the query to the tiling job.
func (q *Equals) Apply(query *elastic.Query) error {
	query.Must(elastic.NewTermsQuery(q.Field, q.Values...))
	return nil
}
