package elastic

import (
	"gopkg.in/olivere/elastic.v3"
)

// Query represents an elasticsearch implementation of the prism.Query
// interface.
type Query interface {
	Get() (elastic.Query, error)
}
