package query

import (
	"fmt"
)

// Equals represents an equality query, checking if a field equals a provided
// value.
type Equals struct {
	Field string
	Value interface{}
}

// // NewEquals instantiates and returns an equals query object.
// func NewEquals(queries map[string]interface{}) (*Equals, error) {
// }

// GetHash returns a string hash of the query.
func (q *Equals) GetHash() string {
	return fmt.Sprintf("%s:%v",
		q.Field,
		q.Value)
}
