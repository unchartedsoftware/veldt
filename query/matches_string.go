package query

import (
	"fmt"

	"github.com/unchartedsoftware/veldt/util/json"
)

// MatchesString query represents a raw string query. The string could be
// A regular expression, or Lucene query, etc. (depends on the implementation
// support chosen.)
type MatchesString struct {
	Match  string
	Fields []string
}

// Parse parses the provided JSON object and populates the querys attributes.
func (q *MatchesString) Parse(params map[string]interface{}) error {
	match, ok := json.GetString(params, "match")
	if !ok {
		return fmt.Errorf("`match` parameter missing from query")
	}
	fields, ok := json.GetStringArray(params, "fields")
	if !ok {
		return fmt.Errorf("`fields` parameter missing from query")
	}
	q.Fields = fields
	q.Match = match
	return nil
}
