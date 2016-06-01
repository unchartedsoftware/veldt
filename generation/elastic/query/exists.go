package query

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/util/json"
)

// Exists represents an elasticsearch exists query.
type Exists struct {
	Field string
}

// NewExists instantiates and returns an exists query object.
func NewExists(params map[string]interface{}) (*Exists, error) {
	field, ok := json.GetString(params, "field")
	if !ok {
		return nil, fmt.Errorf("Exists `field` parameter missing from tiling param %v", params)
	}
	return &Exists{
		Field: field,
	}, nil
}

// GetHash returns a string hash of the query.
func (q *Exists) GetHash() string {
	return fmt.Sprintf("%s:%s", q.Field)
}

// GetQuery returns the elastic query object.
func (q *Exists) GetQuery() elastic.Query {
	return elastic.NewExistsQuery(q.Field)
}
