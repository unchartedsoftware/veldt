package elastic

import (
	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/query"
)

// Has represents an elasticsearch terms query.
type Has struct {
	query.Has
}

// NewHas instantiates and returns a new tile struct.
func NewHas() (veldt.Query, error) {
	return &Has{}, nil
}

// Get returns the appropriate elasticsearch query for the query.
func (q *Has) Get() (elastic.Query, error) {
	return elastic.NewTermsQuery(q.Field, q.Values...), nil
}
