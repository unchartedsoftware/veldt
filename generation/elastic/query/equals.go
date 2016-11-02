package query

import (
	"fmt"
	"sort"
	"strings"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/query"
	"github.com/unchartedsoftware/prism/util/json"
)

// Equals represents an elasticsearch term query.
type Equals struct {
	query.Equals
}

// // NewEquals instantiates and returns a terms query object.
// func NewEquals(params map[string]interface{}) (*Terms, error) {
// 	field, ok := json.GetString(params, "field")
// 	if !ok {
// 		return nil, fmt.Errorf("Terms `field` parameter missing from tiling param %v", params)
// 	}
// 	terms, ok := json.GetStringArray(params, "terms")
// 	if !ok {
// 		return nil, fmt.Errorf("Terms `terms` parameter missing from tiling param %v", params)
// 	}
// 	// sort terms
// 	sort.Strings(terms)
// 	return &Terms{
// 		Field: field,
// 		Terms: terms,
// 	}, nil
// }

// Apply adds the query to the tiling job.
func (q *Equals) Apply(query *elastic.Query) error {
	query.Must(elastic.NewTermQuery(q.Field, q.Value))
	return nil
}
