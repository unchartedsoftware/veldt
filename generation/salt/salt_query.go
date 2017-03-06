package salt

import (
	"github.com/unchartedsoftware/veldt"
)

// Query encapsulates a query to the salt server
type Query struct {
	queryConfiguration map[string]interface{}
}

// NewSaltQuery instantiates and returns a new query structure for use with Salt
func NewSaltQuery() (veldt.Query, error) {
	return &Query{}, nil
}

// GetQueryConfiguration provides the full configuration with which the query was configured
func (q *Query) GetQueryConfiguration () map[string]interface{} {
	return q.queryConfiguration
}

// Parse stores the provided JSON object for shipping to Salt
func (q *Query) Parse (params map[string]interface{}) error {
	q.queryConfiguration = params
	return nil
}
