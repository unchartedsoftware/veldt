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

// IsBoolOperator returns true if the operator is a unary operator.
func IsBoolOperator(op string) bool {
	return IsBinaryOperator(op) || IsUnaryOperator(op)
}

// IsBinaryOperator returns true if the operator is a binary operator.
func IsBinaryOperator(op string) bool {
	return op == And || op == Or
}

// IsUnaryOperator returns true if the operator is a unary operator.
func IsUnaryOperator(op string) bool {
	return op == Not
}

// BinaryExpression represents a binary boolean expression.
type BinaryExpression struct {
	Left  Query
	Op    string
	Right Query
}

// Apply adds the query to the tiling job.
func (q *BinaryExpression) Apply(arg interface{}) error {
	return fmt.Errorf("Not implemented")
}

// GetHash returns the hash of the query.
func (q *BinaryExpression) GetHash() string {
	return fmt.Sprintf("%s:%s:%s",
		q.Left.GetHash(),
		q.Op,
		q.Right.GetHash())
}

// UnaryExpression represents a unary boolean expression.
type UnaryExpression struct {
	Query Query
	Op    string
}

// Apply adds the query to the tiling job.
func (q *UnaryExpression) Apply(arg interface{}) error {
	return fmt.Errorf("Not implemented")
}

// GetHash returns the hash of the query.
func (q *UnaryExpression) GetHash() string {
	return fmt.Sprintf("%s:%s",
		q.Op,
		q.Query.GetHash())
}
