package citus

import (
	"fmt"

	"github.com/unchartedsoftware/prism"
)

// BinaryExpression represents an and/or boolean query.
type BinaryExpression struct {
	prism.BinaryExpression
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
	case prism.And:
		// AND
		res = fmt.Sprintf("((%s) AND (%s))", queryStringLeft, queryStringRight)
	case prism.Or:
		// OR
		res = fmt.Sprintf("((%s) OR (%s))", queryStringLeft, queryStringRight)
	default:
		return "", fmt.Errorf("`%v` operator is not a valid binary operator", e.Op)
	}
	return res, nil
}

// UnaryExpression represents a must_not boolean query.
type UnaryExpression struct {
	prism.UnaryExpression
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

	res := "NOT "
	switch e.Op {
	case prism.Not:
		// NOT
		res = res + fmt.Sprintf("(%s)", a)
	default:
		return "", fmt.Errorf("`%v` operator is not a valid unary operator", e.Op)
	}
	return res, nil
}
