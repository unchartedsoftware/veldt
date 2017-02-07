package veldt

import (
	"fmt"

	"github.com/unchartedsoftware/veldt/util/json"
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
	Left  Query
	Op    string
	Right Query
}

// Parse should parse through the provided JSON and populate the struct fields.
func (b *BinaryExpression) Parse(params map[string]interface{}) error {
	// left
	l, ok := params["left"]
	if !ok {
		return fmt.Errorf("`left` parameter missing from query")
	}
	left, ok := l.(Query)
	if !ok {
		return fmt.Errorf("`left` is not a query type")
	}
	b.Left = left
	// op
	op, ok := json.GetString(params, "op")
	if !ok {
		return fmt.Errorf("`op` parameter missing from query")
	}
	if op != And && op != Or {
		return fmt.Errorf("`op` parameter value is not recognized")
	}
	b.Op = op
	// right
	r, ok := params["right"]
	if !ok {
		return fmt.Errorf("`right` parameter missing from query")
	}
	right, ok := r.(Query)
	if !ok {
		return fmt.Errorf("`right` is not a query type")
	}
	b.Right = right
	return nil
}

// UnaryExpression represents a unary boolean expression.
type UnaryExpression struct {
	Query Query
	Op    string
}

// Parse should parse through the provided JSON and populate the struct fields.
func (u *UnaryExpression) Parse(params map[string]interface{}) error {
	// left
	q, ok := params["query"]
	if !ok {
		return fmt.Errorf("`query` parameter missing from query")
	}
	query, ok := q.(Query)
	if !ok {
		return fmt.Errorf("`query` is not a query type")
	}
	u.Query = query
	// op
	op, ok := json.GetString(params, "op")
	if !ok {
		return fmt.Errorf("`op` parameter missing from query")
	}
	if op != Not {
		return fmt.Errorf("`op` parameter value is not recognized")
	}
	u.Op = op
	return nil
}
