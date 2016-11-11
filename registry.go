package tile

import (
	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/query"
)

var (
	// registry contains all registered tile generator constructors.
	registry = make(map[string]*Pipeline)
)

// Register registers a tile generator under the provided type id string.
func Register(typeID string, pipeline *Pipeline) {
	registry[typeID] = pipeline
}

// GetGenerator instantiates a tile generator from a tile request.
func GetPipeline(id string) (*Pipeline, error) {
	pipeline, ok := registry[id]
	if !ok {
		return nil, fmt.Errorf("Pipeline ID of '%s' is not recognized", id)
	}
	return pipeline, nil
}
