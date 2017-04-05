package salt

import (
	"fmt"
	"github.com/unchartedsoftware/veldt"
)



// This file contains general boolean query expressions



// BinaryExpression is an extension of the base veldt.BinaryExpression that
// allows salt to reconstruct the configuration, so as to send it to the salt
// server.
type BinaryExpression struct {
	veldt.BinaryExpression
}

// UnaryExpression is an extension of the base veldt.UnaryExpression that
// allows salt to reconstruct the configuration, so as to send it to the salt
// server.
type UnaryExpression struct {
	veldt.UnaryExpression
}

// NewBinaryExpression instantiates and returns a new binary expression
func NewBinaryExpression() (veldt.Query, error) {
	return &BinaryExpression{}, nil
}

// NewUnaryExpression instantiates and returns a new unary expression
func NewUnaryExpression() (veldt.Query, error) {
	return &UnaryExpression{}, nil
}

// Get retrieves the configuration from a query for use by the salt server
func (e *BinaryExpression) Get () (map[string]interface{}, error) {
	left, ok := e.Left.(Query)
	if !ok {
		return nil, fmt.Errorf("`Left` is not of type salt.Query")
	}
	right, ok := e.Right.(Query)
	if !ok {
		return nil, fmt.Errorf("`Right` is not of type salt.Query")
	}

	leftConfig, err := left.Get()
	if err != nil {
		return nil, err
	}
	rightConfig, err := right.Get()
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	result["operation"] = e.Op
	result["left"] = leftConfig
	result["right"] = rightConfig

	return result, nil
}

// Get retrieves the configuration from a query for use by the salt server
func (e *UnaryExpression) Get () (map[string]interface{}, error) {
	q, ok := e.Query.(Query)
	if !ok {
		return nil, fmt.Errorf("`Operand is not of type salt.Query")
	}
	qConfig, err := q.Get()
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	result["operation"] = e.Op
	result["operand"] = qConfig

	return result, nil
}
