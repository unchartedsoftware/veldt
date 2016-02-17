package store

import (
	"fmt"
)

var (
	// registry contains all registered store constructors.
	registry = make(map[string]ConnectionConstructor)
)

// Connection represents an interface for connecting to, setting, and retreiving
// values from a key-value database or in-memory storage server.
type Connection interface {
	Set(string, []byte) error
	Get(string) ([]byte, error)
	Exists(string) (bool, error)
	Close()
}

// ConnectionConstructor represents a function to instantiate a new Store
// connection.
type ConnectionConstructor func() (Connection, error)

// Register registers a meta data generator under the provided type id string.
func Register(typeID string, conn ConnectionConstructor) {
	registry[typeID] = conn
}

// GetConnection instantiates a and returns a store connection.
func GetConnection(typeID string) (Connection, error) {
	ctor, ok := registry[typeID]
	if !ok {
		return nil, fmt.Errorf("Store connection type '%s' is not recognized in request", typeID)
	}
	return ctor()
}
