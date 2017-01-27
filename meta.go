package veldt

// Meta represents an interface for generating meta data.
type Meta interface {
	Create(string) ([]byte, error)
	Parse(map[string]interface{}) error
}

// MetaCtor represents a function that instantiates and returns a new meta
// data type.
type MetaCtor func() (Meta, error)
