package tile

import (
	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/query"
)

// Meta represents an interface for generating tile data.
type Meta interface {
	Create(string) ([]byte, error)
	Parse(map[string]interface{}) error
}
