package query

import (
	//"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/query"
)

// BinaryExpression represents an must / should boolean query.
type BinaryExpression struct {
	query.BinaryExpression
}

// Apply adds the query to the tiling job.
func (q *BinaryExpression) Apply(arg interface{}) error {
	query, ok := arg.(*elastic.BoolQuery)
	if !ok {
		return fmt.Errorf("`%v` is not of type *elastic.BoolQuery", arg)
	}
	binary := elastic.NewBoolQuery()
	switch q.Op {
	case query.And:
		// AND
		left := elastic.NewBoolQuery()
		right := elastic.NewBoolQuery()
		err := q.Left.Apply(left)
		if err != nil {
			return err
		}
		err := q.Right.Apply(right)
		if err != nil {
			return err
		}
		binary.Must(left)
		binary.Must(right)
	case query.Or:
		// OR
		left := elastic.NewBoolQuery()
		right := elastic.NewBoolQuery()
		err := q.Left.Apply(left)
		if err != nil {
			return err
		}
		err := q.Right.Apply(right)
		if err != nil {
			return err
		}
		binary.Should(left)
		binary.Should(right)
	default:
		return fmt.Errorf("`%v` operator is not a valid binary operator", q.Op)
	}
	query.Must(binary)
	return nil
}

// UnaryExpression represents a must_not boolean query.
type UnaryExpression struct {
	query.UnaryExpression
}

// Apply adds the query to the tiling job.
func (q *BinaryExpression) Apply(arg interface{}) error {
	query, ok := arg.(*elastic.BoolQuery)
	if !ok {
		return fmt.Errorf("`%v` is not of type *elastic.BoolQuery", arg)
	}
	unary := elastic.NewBoolQuery()
	switch q.Op {
	case query.Not:
		// NOT
		sub := elastic.NewBoolQuery()
		err := q.Left.Apply(sub)
		if err != nil {
			return err
		}
		unary.MustNot(sub)
	default:
		return fmt.Errorf("`%v` operator is not a valid unary operator", q.Op)
	}
	query.Must(unary)
	return nil
}
