package query

// Query represents a base query interface.
type Query interface {
	GetHash() string
	Apply(interface{}) error
}
