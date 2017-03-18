package salt

import (
	"github.com/unchartedsoftware/veldt"
)

// Query represents a salt implementation of the veldt.Query interface
type Query interface {
	veldt.Query
	Get () (map[string]interface{}, error)
}

// GenericQuery represents any specific non-boolean query for use with salt
type GenericQuery struct {
	queryType string
	parameters map[string]interface{}
}


// NewGenericQuery instantiates and returns a new generic query structure for
// use with salt
func NewGenericQuery (queryType string) veldt.QueryCtor {
	return func () (veldt.Query, error) {
		return &GenericQuery{queryType, nil}, nil
	}
}

// Parse stores a query's configuration for later use by the salt server
func (q *GenericQuery) Parse (parameters map[string]interface{}) error {
	q.parameters = parameters
	return nil
}

// Get retrieves the configuration from a query for use by the salt server
func (q *GenericQuery) Get () (map[string]interface{}, error) {
	result := make(map[string]interface{})
	result["operation"] = q.queryType
	result["parameters"] = q.parameters
	return result, nil
}

