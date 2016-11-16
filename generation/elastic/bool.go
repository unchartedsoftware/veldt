package elastic

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism"
)

// BinaryExpression represents an must / should boolean query.
type BinaryExpression struct {
	prism.BinaryExpression
}

// Apply adds the query to the tiling job.
func (e *BinaryExpression) Get() (elastic.Query, error) {

	left, ok := e.Left.(Query)
	if !ok {
		return nil, fmt.Errorf("Left is not of type elastic.Query")
	}
	right, ok := e.Right.(Query)
	if !ok {
		return nil, fmt.Errorf("Right is not of type elastic.Query")
	}

	a, err := left.Get()
	if err != nil {
		return nil, err
	}
	b, err := right.Get()
	if err != nil {
		return nil, err
	}

	res := elastic.NewBoolQuery()
	switch e.Op {
	case prism.And:
		// AND
		res.Must(a)
		res.Must(b)
	case prism.Or:
		// OR
		res.Should(a)
		res.Should(b)
	default:
		return nil, fmt.Errorf("`%v` operator is not a valid binary operator", e.Op)
	}
	return res, nil
}

// UnaryExpression represents a must_not boolean query.
type UnaryExpression struct {
	prism.UnaryExpression
}

// Apply adds the query to the tiling job.
func (e *UnaryExpression) Get() (elastic.Query, error) {

	q, ok := e.Query.(Query)
	if !ok {
		return nil, fmt.Errorf("Left is not of type elastic.Query")
	}

	a, err := q.Get()
	if err != nil {
		return nil, err
	}

	res := elastic.NewBoolQuery()
	switch e.Op {
	case prism.Not:
		// NOT
		res.MustNot(a)
	default:
		return nil, fmt.Errorf("`%v` operator is not a valid unary operator", e.Op)
	}
	return res, nil
}
