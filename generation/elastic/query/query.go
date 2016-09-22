package query

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/util/json"
)

// Query represents a base query Query interface.
type Query interface {
	GetQuery() elastic.Query
	GetHash() string
}

func getQueryByType(query map[string]interface{}) (Query, error) {
	// bool
	params, ok := json.GetChild(query, "bool")
	if ok {
		return NewBool(params)
	}
	// exists
	params, ok = json.GetChild(query, "exists")
	if ok {
		return NewExists(params)
	}
	// terms
	params, ok = json.GetChild(query, "terms")
	if ok {
		return NewTerms(params)
	}
	// range
	params, ok = json.GetChild(query, "range")
	if ok {
		return NewRange(params)
	}
	// prefix
	params, ok = json.GetChild(query, "prefix")
	if ok {
		return NewPrefix(params)
	}
	// query_string
	params, ok = json.GetChild(query, "query_string")
	if ok {
		return NewString(params)
	}
	// match
	params, ok = json.GetChild(query, "match")
	if ok {
		return NewMatch(params)
	}
	// has_parent
	params, ok = json.GetChild(query, "has_parent")
	if ok {
		return NewHasParent(params)
	}
	// has_child
	params, ok = json.GetChild(query, "has_child")
	if ok {
		return NewHasChild(params)
	}
	return nil, fmt.Errorf("No testrecognized query type found in %v", query)
}
