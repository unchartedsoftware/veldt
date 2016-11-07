package query

import (
	"fmt"
)

var (
	// registry contains all registered meta data generator constructors.
	registry = make(map[string]Constructor)
)

// Register registers a meta data generator under the provided type id string.
func Register(typeID string, ctor Constructor) {
	registry[typeID] = ctor
}

// GetQuery instantiates a meta data generator from a meta data request.
func GetQuery(typeID string, args map[string]interface{}) (Query, error) {
	ctor, ok := registry[typeID]
	if !ok {
		return nil, fmt.Errorf("query `%s` has not been registered",
			typeID)
	}
	return ctor(args)
}
