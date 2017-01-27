package veldt

import (
	"fmt"
)

var (
	// registry contains all registered tile generator constructors.
	registry = make(map[string]*Pipeline)
)

// Register registers a pipeline under the provided ID string.
func Register(typeID string, p *Pipeline) {
	registry[typeID] = p
}

// GetPipeline retrieves the pipeline registered under the provided ID string.
func GetPipeline(id string) (*Pipeline, error) {
	p, ok := registry[id]
	if !ok {
		return nil, fmt.Errorf("Pipeline ID of '%s' is not recognized", id)
	}
	return p, nil
}
