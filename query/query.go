package query

// Query represents a base query interface.
type Query interface {
	Apply(interface{}) error
}

// Constructor represents a query constructor.
type Constructor func(params map[string]interface{}) (Query, error)
