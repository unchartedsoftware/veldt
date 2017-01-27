package citus

import (
	"fmt"

	"github.com/unchartedsoftware/veldt"
)

// BinaryExpression represents an and/or boolean query.
type BinaryExpression struct {
	veldt.BinaryExpression
}

// NewBinaryExpression instantiates and returns a new binary expression.
func NewBinaryExpression() (veldt.Query, error) {
	return &BinaryExpression{}, nil
}

// Get adds the parameters to the query and returns the string representation.
func (e *BinaryExpression) Get(query *Query) (string, error) {

	left, ok := e.Left.(QueryString)
	if !ok {
		return "", fmt.Errorf("Left is not of type citus.Query")
	}
	right, ok := e.Right.(QueryString)
	if !ok {
		return "", fmt.Errorf("Right is not of type citus.Query")
	}

	queryStringLeft, err := left.Get(query)
	if err != nil {
		return "", err
	}
	queryStringRight, err := right.Get(query)
	if err != nil {
		return "", err
	}

	res := ""
	switch e.Op {
	case veldt.And:
		// AND
		res = fmt.Sprintf("((%s) AND (%s))", queryStringLeft, queryStringRight)
	case veldt.Or:
		// OR
		res = fmt.Sprintf("((%s) OR (%s))", queryStringLeft, queryStringRight)
	default:
		return "", fmt.Errorf("`%v` operator is not a valid binary operator", e.Op)
	}
	return res, nil
}

// UnaryExpression represents a must_not boolean query.
type UnaryExpression struct {
	veldt.UnaryExpression
}

// NewUnaryExpression instantiates and returns a new unary expression.
func NewUnaryExpression() (veldt.Query, error) {
	return &UnaryExpression{}, nil
}

// Get adds the parameters to the query and returns the string representation.
func (e *UnaryExpression) Get(query *Query) (string, error) {

	q, ok := e.Query.(QueryString)
	if !ok {
		return "", fmt.Errorf("Left is not of type citus.Query")
	}

	a, err := q.Get(query)
	if err != nil {
		return "", err
	}

	res := ""
	switch e.Op {
	case veldt.Not:
		// NOT
		res = res + fmt.Sprintf("NOT (%s)", a)
	default:
		return "", fmt.Errorf("`%v` operator is not a valid unary operator", e.Op)
	}
	return res, nil
}
