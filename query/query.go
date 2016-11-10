package query

// Query represents a base query interface.
type Query interface {
	Parse(map[string]interface{})
	Apply(interface{}) error
}
