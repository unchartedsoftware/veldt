package elastic

import (
	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/query"
)

// MatchesString represents an elasticsearch query-string query.
type MatchesString struct {
	query.MatchesString
}

// NewMatchesString instantiates and returns a new struct.
func NewMatchesString() (veldt.Query, error) {
	return &MatchesString{}, nil
}

// Get returns the appropriate elasticsearch query for the query.
func (q *MatchesString) Get() (elastic.Query, error) {
	query := elastic.NewQueryStringQuery(q.Match)
	for _, f := range q.Fields {
		query = query.Field(f)
	}
	return query, nil
}
