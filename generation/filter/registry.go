package filter

import (
	"fmt"
)

// Generator represents a function which takes filter params and returns a
// filter
type Generator func(map[string]interface{}) (interface{}, bool)

// registry contains all tiling function implementations.
var (
	registry = make(map[string]Generator)
)

// Register registers a meta generator under the provided type id string.
func Register(typeID string, filter Generator) {
	registry[typeID] = filter
}

// GetGeneratorByType when given a string id will return the registered
// meta generator.
func GetGeneratorByType(typeID string) (Generator, error) {
	filter, ok := registry[typeID]
	if !ok {
		return nil, fmt.Errorf("Filter type '%s' is not recognized", typeID)
	}
	return filter, nil
}
