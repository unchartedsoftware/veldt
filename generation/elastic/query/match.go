package query

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/util/json"
)

// Match represents an elasticsearch match query.
type Match struct {
	Field  string
	String string
	Type   string
}

// NewMatch instantiates and returns a match query object.
func NewMatch(params map[string]interface{}) (*Match, error) {
	field, ok := json.GetString(params, "field")
	if !ok {
		return nil, fmt.Errorf("Match `field` parameter missing from tiling param %v", params)
	}
	str, ok := json.GetString(params, "string")
	if !ok {
		return nil, fmt.Errorf("Match `string` parameter missing from tiling param %v", params)
	}
	matchType := json.GetStringDefault(params, "boolean", "type")
	return &Match{
		Field:  field,
		String: str,
		Type: matchType,
	}, nil
}

// GetHash returns a string hash of the query.
func (q *Match) GetHash() string {
	return fmt.Sprintf("%s:%s",
		q.Field,
		q.String)
}

// GetQuery returns the elastic query object.
func (q *Match) GetQuery() elastic.Query {
	matchQuery := elastic.NewMatchQuery(q.Field, q.String)
	matchQuery.Type(q.Type)
	return matchQuery
}
