package query

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/util/json"
)

// String represents an elasticsearch query_string query.
type String struct {
	Field  string
	String string
}

// NewString instantiates and returns a query string query object.
func NewString(params map[string]interface{}) (*String, error) {
	field, ok := json.GetString(params, "field")
	if !ok {
		return nil, fmt.Errorf("String `field` parameter missing from tiling param %v", params)
	}
	str, ok := json.GetString(params, "string")
	if !ok {
		return nil, fmt.Errorf("String `string` parameter missing from tiling param %v", params)
	}
	return &String{
		Field:  field,
		String: str,
	}, nil
}

// GetHash returns the hash for the query object.
func (q *String) GetHash() string {
	return fmt.Sprintf("%s:%s",
		q.Field,
		q.String)
}

// GetQuery returns the elastic query object.
func (q *String) GetQuery() elastic.Query {
	return elastic.NewQueryStringQuery(q.String).
		Field(q.Field)
}
