package query

import (
	"fmt"
	"strings"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/util/json"
)

// Terms represents an elasticsearch terms query.
type Terms struct {
	Field string
	Terms []string
}

// NewTerms instantiates and returns a terms query object.
func NewTerms(params map[string]interface{}) (*Terms, error) {
	field, ok := json.GetString(params, "field")
	if !ok {
		return nil, fmt.Errorf("Terms `field` parameter missing from tiling param %v", params)
	}
	terms, ok := json.GetStringArray(params, "terms")
	if !ok {
		return nil, fmt.Errorf("Terms `terms` parameter missing from tiling param %v", params)
	}
	return &Terms{
		Field: field,
		Terms: terms,
	}, nil
}

// GetHash returns a string hash of the query.
func (q *Terms) GetHash() string {
	return fmt.Sprintf("%s:%s",
		q.Field,
		strings.Join(q.Terms, ":"))
}

// GetQuery returns the elastic query object.
func (q *Terms) GetQuery() elastic.Query {
	terms := make([]interface{}, len(q.Terms))
	for i, term := range q.Terms {
		terms[i] = term
	}
	return elastic.NewTermsQuery(q.Field, terms...)
}
