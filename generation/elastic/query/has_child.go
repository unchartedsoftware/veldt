package query

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/util/json"
)

// HasChild represents an elasticsearch has_child query.
type HasChild struct {
	ChildType string
	Query     Bool
}

// NewHasChild instantiates and returns an has_child query object.
func NewHasChild(params map[string]interface{}) (*HasChild, error) {
	childType, ok := json.GetString(params, "type")
	if !ok {
		return nil, fmt.Errorf("has_child `type` parameter missing from tiling param %v", params)
	}
	query, ok := json.GetChild(params, "query", "bool")
	if !ok {
		return nil, fmt.Errorf("has_child `query` or `bool` parameter(s) missing from tiling param %v", params)
	}
	boolQuery, err := NewBool(query)
	if err != nil {
		return nil, err
	}
	return &HasChild{
		ChildType: childType,
		Query:     *boolQuery,
	}, nil
}

// GetHash returns a string hash of the query.
func (q *HasChild) GetHash() string {
	return fmt.Sprintf("%s::%s", q.ChildType, q.Query.GetHash())
}

// GetQuery returns the elastic query object.
func (q *HasChild) GetQuery() elastic.Query {
	return elastic.NewHasChildQuery(q.ChildType, q.Query.GetQuery())
}
