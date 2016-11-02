package query

import (
	"fmt"
)

const (
	// And represents an AND binary operator
	And = "AND"

	// Or represents an OR binary operator
	Or = "OR"

	// Not represents a NOT unary operator.
	Not = "NOT"
)

// BinaryExpression represents a binary boolean expression.
type BinaryExpression struct {
	A Query
	B Query
	Op string
}

// GetHash returns the hash of the query.
func (q *BinaryExpression) GetHash() string {
	return fmt.Sprintf("%s:%s:%s",
		q.A.GetHash(),
		q.Op,
		q.B.GetHash())
}

// UnaryExpression represents a unary boolean expression.
type UnaryExpression struct {
	Q Query
	Op string
}

// GetHash returns the hash of the query.
func (q *UnaryExpression) GetHash() string {
	return fmt.Sprintf("%s:%s",
		q.Op,
		q.Q.GetHash())
}
