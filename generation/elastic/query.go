package elastic

import (
	"gopkg.in/olivere/elastic.v3"
)

// Query represents an elasticsearch implementation of the veldt.Query
// interface.
type Query interface {
	Get() (elastic.Query, error)
}
