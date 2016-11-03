package query

// Query represents a base query interface.
type Query interface {
	GetHash() string
	Apply(interface{}) error
}

// Constructor represents a query constructor.
type Constructor func(params map[string]interface{}) (Query, error)
