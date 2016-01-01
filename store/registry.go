package store

import (
	"fmt"
)

var (
	// registry contains all registered store constructors.
	registry = make(map[string]ConnectionConstructor)
)

// Request represents a storage request.
type Request struct {
	Type     string `json:"type"`
	Endpoint string `json:"endpoint"`
}

// String returns the request formatted as a string.
func (r *Request) String() string {
	return fmt.Sprintf("%s/%s/",
		r.Endpoint,
		r.Type)
}

// Connection represents an interface for connecting to, setting, and retreiving
// values from a key-value database or in-memory storage server.
type Connection interface {
	Set(string, []byte) error
	Get(string) ([]byte, error)
	Exists(string) (bool, error)
	Close()
}

// ConnectionConstructor represents a function to instantiate a new generator
// from a meta data request.
type ConnectionConstructor func(*Request) (Connection, error)

// Register registers a meta data generator under the provided type id string.
func Register(typeID string, conn ConnectionConstructor) {
	registry[typeID] = conn
}

// GetConnection instantiates a and returns a store connection.
func GetConnection(req *Request) (Connection, error) {
	ctor, ok := registry[req.Type]
	if !ok {
		return nil, fmt.Errorf("Store connection type '%s' is not recognized in request %s", req.Type, req.String())
	}
	return ctor(req)
}
