package elastic

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/veldt"
)

// BinaryExpression represents an must / should boolean query.
type BinaryExpression struct {
	veldt.BinaryExpression
}

// NewBinaryExpression instantiates and returns a new binary expression.
func NewBinaryExpression() (veldt.Query, error) {
	return &BinaryExpression{}, nil
}

// Parse implements parse explicitly on an elastic.BinaryExpression, so that
// elastic.BinaryExpression will implement veldt.Query, which is necessary for
// multi-level binary consolidation
func (b *BinaryExpression) Parse(params map[string]interface{}) error {
	return b.BinaryExpression.Parse(params)
}

// clauses gets the clauses of this binary expression for construction of
// an elastic query.  It consolidates clauses across multiple levels of a
// query so as to simplify the elastic query.
func (e *BinaryExpression) clauses () ([]elastic.Query, error) {
	left, ok := e.Left.(Query)
	if !ok {
		return nil, fmt.Errorf("`Left` is not of type elastic.Query")
	}
	right, ok := e.Right.(Query)
	if !ok {
		return nil, fmt.Errorf("`Right` is not of type elastic.Query")
	}

	clauses := make([]elastic.Query, 0, 0)

	// Get left-hand clauses
	beLeft, leftIsBool := e.Left.(*BinaryExpression)
	if leftIsBool && beLeft.Op == e.Op {
		aClauses, err := beLeft.clauses()
		if err != nil {
			return nil, err
		}
		clauses = append(clauses, aClauses...)
	} else {
		a, err := left.Get()
		if err != nil {
			return nil, err
		}
		clauses = append(clauses, a)
	}

	// Get right-hand clauses
	beRight, rightIsBool := e.Right.(*BinaryExpression)
	if rightIsBool && beRight.Op == e.Op {
		bClauses, err := beRight.clauses()
		if err != nil {
			return nil, err
		}
		clauses = append(clauses, bClauses...)
	} else {
		b, err := right.Get()
		if err != nil {
			return nil, err
		}
		clauses = append(clauses, b)
	}

	return clauses, nil
}

// Get returns the appropriate elasticsearch query for the binary expression.
func (e *BinaryExpression) Get() (elastic.Query, error) {
	clauses, err := e.clauses()
	if err != nil {
		return nil, err
	}
	
	res := elastic.NewBoolQuery()
	switch e.Op {
	case veldt.And:
		// AND
		for _, clause := range(clauses) {
			res.Must(clause)
		}
	case veldt.Or:
		// OR
		for _, clause := range(clauses) {
			res.Should(clause)
		}
	default:
			return nil, fmt.Errorf("`%v` operator is not a valid binary operator", e.Op)
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

// Get returns the appropriate elasticsearch query for the unary expression.
func (e *UnaryExpression) Get() (elastic.Query, error) {

	q, ok := e.Query.(Query)
	if !ok {
		return nil, fmt.Errorf("`Query` is not of type elastic.Query")
	}

	a, err := q.Get()
	if err != nil {
		return nil, err
	}

	res := elastic.NewBoolQuery()
	switch e.Op {
	case veldt.Not:
		// NOT
		res.MustNot(a)
	default:
		return nil, fmt.Errorf("`%v` operator is not a valid unary operator", e.Op)
	}
	return res, nil
}
