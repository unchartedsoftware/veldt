package query

import (
	"fmt"
)

// Exists represents an exists query checking if a field is not null.
type Exists struct {
	Field string
}

// GetHash returns a string hash of the query.
func (q *Exists) GetHash() string {
	return fmt.Sprintf("%s", q.Field)
}
