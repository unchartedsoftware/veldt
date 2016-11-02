package query

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/query"
	"github.com/unchartedsoftware/prism/util/json"
)

// Exists represents an elasticsearch exists query.
type Exists struct {
	query.Exists
}

// // NewExists instantiates and returns an exists query object.
// func NewExists(params map[string]interface{}) (query.Query, error) {
// 	field, ok := json.GetString(params, "field")
// 	if !ok {
// 		return nil, fmt.Errorf("Exists `field` parameter missing from tiling param %v", params)
// 	}
// 	return &Exists{
// 		Field: field,
// 	}, nil
// }

// Apply adds the query to the tiling job.
func (q *Exists) Apply(query *elastic.Query) error {
	query.Must(elastic.NewExistsQuery(q.Field))
	return nil
}
