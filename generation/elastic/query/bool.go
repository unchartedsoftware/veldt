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
func (q *BinaryExpression) Get() elastic.Query {
	query := elastic.NewBoolQuery()
	switch q.Op {
	case query.And:
		// AND
		query.Must(q.Left.Get())
		query.Must(q.Right.Get())
	case query.Or:
		// OR
		query.Should(q.Left.Get())
		query.Should(q.Right.Get())
	default:
		return fmt.Errorf("`%v` operator is not a valid binary operator", q.Op)
	}
	return query
}

// UnaryExpression represents a must_not boolean query.
type UnaryExpression struct {
	query.UnaryExpression
}

// Apply adds the query to the tiling job.
func (q *UnaryExpression) Apply(arg interface{}) error {
	query := elastic.NewBoolQuery()
	switch q.Op {
	case query.Not:
		// NOT
		unary.MustNot(q.Query.Get())
	default:
		return fmt.Errorf("`%v` operator is not a valid unary operator", q.Op)
	}
	return query
}
