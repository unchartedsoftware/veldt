package query

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/util/json"
)

// HasParent represents an elasticsearch has_parent query.
type HasParent struct {
	ParentType string
	Query      Bool
}

// NewHasParent instantiates and returns an has_parent query object.
func NewHasParent(params map[string]interface{}) (*HasParent, error) {
	parentType, ok := json.GetString(params, "parent_type")
	if !ok {
		return nil, fmt.Errorf("has_parent `parent_type` parameter missing from tiling param %v", params)
	}
	query, ok := json.GetChild(params, "query")
	if !ok {
		return nil, fmt.Errorf("has_parent `query` parameter missing from tiling param %v", params)
	}
	bool_, ok := json.GetChild(query, "bool")
	if !ok {
		return nil, fmt.Errorf("has_child `bool` parameter missing from tiling param %v", params)
	}
	boolQuery, err := NewBool(bool_)
	if err != nil {
		return nil, err
	}
	return &HasParent{
		ParentType: parentType,
		Query:      *boolQuery,
	}, nil
}

// GetHash returns a string hash of the query.
func (q *HasParent) GetHash() string {
	return fmt.Sprintf("%s::%s", q.ParentType, q.Query.GetHash())
}

// GetQuery returns the elastic query object.
func (q *HasParent) GetQuery() elastic.Query {
	return elastic.NewHasParentQuery(q.ParentType, q.Query.GetQuery())
}
