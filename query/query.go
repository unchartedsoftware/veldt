package query

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/util/json"
)

// Query represents a base query interface.
type Query interface {
	GetHash() string
	Apply(interface{}) error
}
