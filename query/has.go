package query

import (
	"fmt"
	"sort"
	"strings"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/util/json"
)

// Has quiery represents a query checking if the field has one or more of the
// provided values. represents an elasticsearch terms query.
type Has struct {
	Field string
	Values []interface{}
}

// NewHas instantiates and returns a has query object.
func NewHas(queries map[string]interface{}) (*Equals, error) {
}

// GetHash returns a string hash of the query.
func (q *Has) GetHash() string {
	hashes := make([]string, len(q.Values)
	for i, val := range q.Values {
		hashes[i] = fmt.Sprintf("%v", val)
	}
	return fmt.Sprintf("%s:%s",
		q.Field,
		strings.Join(hashes, ":"))
}
