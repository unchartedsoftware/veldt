package prism

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

func (q *BinaryExpression) Parse(params map[string]interface{}) error {
	return fmt.Errorf("not implemented")
}

// UnaryExpression represents a unary boolean expression.
type UnaryExpression struct {
	Query Query
	Op    string
}

func (q *UnaryExpression) Parse(params map[string]interface{}) error {
	return fmt.Errorf("not implemented")
}
