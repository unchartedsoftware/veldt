package tile

import (
	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/query"
)

// Store represents an interface for connecting to, setting, and retrieving
// values from a key-value database or in-memory storage server.
type Store interface {
	Set(string, []byte) error
	Get(string) ([]byte, error)
	Exists(string) (bool, error)
	Close()
}
