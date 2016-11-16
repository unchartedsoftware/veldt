package prism

import (
	"fmt"
)

var (
	// registry contains all registered tile generator constructors.
	registry = make(map[string]*Pipeline)
)

// Register registers a tile generator under the provided type id string.
func Register(typeID string, p *Pipeline) {
	registry[typeID] = p
}

// GetGenerator instantiates a tile generator from a tile request.
func GetPipeline(id string) (*Pipeline, error) {
	p, ok := registry[id]
	if !ok {
		return nil, fmt.Errorf("Pipeline ID of '%s' is not recognized", id)
	}
	return p, nil
}
