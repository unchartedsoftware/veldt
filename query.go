package veldt

// Query represents a base query interface.
type Query interface {
	Parse(map[string]interface{}) error
}

// QueryCtor represents a function that instantiates and returns a new query
// type.
type QueryCtor func() (Query, error)
