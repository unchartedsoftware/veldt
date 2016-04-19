package query

import (
	"fmt"
	"sort"
	"strings"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/util/json"
)

// Prefix represents an elasticsearch prefix query.
type Prefix struct {
	Field    string
	Prefixes []string
}

// NewPrefix instantiates and returns a a prefix query object.
func NewPrefix(params map[string]interface{}) (*Prefix, error) {
	field, ok := json.GetString(params, "field")
	if !ok {
		return nil, fmt.Errorf("Prefix `field` parameter missing from tiling param %v", params)
	}
	prefixes, ok := json.GetStringArray(params, "prefixes")
	if !ok {
		return nil, fmt.Errorf("Prefix `prefixes` parameter missing from tiling param %v", params)
	}
	// sort prefixes
	sort.Strings(prefixes)
	return &Prefix{
		Field:    field,
		Prefixes: prefixes,
	}, nil
}

// GetHash returns a slice of
func (q *Prefix) GetHash() string {
	return fmt.Sprintf("%s:%s",
		q.Field,
		strings.Join(q.Prefixes, ":"))
}

// GetQuery returns a slice of prefix queries.
func (q *Prefix) GetQuery() elastic.Query {
	query := elastic.NewBoolQuery()
	for _, prefix := range q.Prefixes {
		query.Should(elastic.NewPrefixQuery(q.Field, prefix))
	}
	return query
}
