package query

import (
	"encoding/json"
	"fmt"

	"github.com/unchartedsoftware/prism"
)

// Parse parses the query payload into the query AST.
func Parse(arg interface{}) (prism.Query, error) {
	// validate the JSON into ot's runtime query components
	validator := NewValidator()
	token, err := validator.Validate(arg)
	if err != nil {
		return nil, err
	}
	// parse into correct AST
	return parseToken(token)
}

func parseExpression(args []interface{}) (prism.Query, error) {
	exp := NewParser(args)
	return exp.parse()
}

func parseToken(token interface{}) (prism.Query, error) {
	// check if token is an expression
	exp, ok := token.([]interface{})
	if ok {
		// is expression, recursively parse it
		return parseExpression(exp)
	}
	// is query, parse it directly
	query, ok := token.(prism.Query)
	if !ok {
		return nil, fmt.Errorf("`%v` token is unrecognized", token)
	}
	return query, nil
}

// expression parses the runtime query expression into it's runtime AST tree.
type expression struct {
	tokens []interface{}
}

func newExpression(arr []interface{}) *expression {
	return &expression{
		tokens: arr,
	}
}

func (e *expression) pop() (interface{}, error) {
	if len(e.tokens) == 0 {
		return nil, fmt.Errorf("Expected operand missing")
	}
	token := e.tokens[0]
	e.tokens = e.tokens[1:len(e.tokens)]
	return token, nil
}

func (e *expression) popOperand() (prism.Query, error) {
	// pops the next operand
	//     cases to consider:
	//         - a) unary operator -> expression
	//         - b) unary operator -> query
	//         - c) expression
	//         - d) query

	// pop next token
	token, err := e.pop()
	if err != nil {
		return nil, err
	}

	// see if it is a unary operator
	op, ok := token.(string)

	isUnary, err := isUnaryOperator(token)
	if err != nil {
		return nil, err
	}

	if ok && isUnary {
		// get next token
		next, err := e.pop()
		if err != nil {
			return nil, err
		}
		// parse token
		query, err := parseToken(next)
		if err != nil {
			return nil, err
		}
		// return unary expression
		return &UnaryExpression{
			Op:    op,
			prism.Query: query,
		}, nil
	}

	// parse token
	return parseToken(token)
}

func (e *expression) peek() interface{} {
	if len(e.tokens) == 0 {
		return nil
	}
	return e.tokens[0]
}

func (e *expression) advance() error {
	if len(e.tokens) < 2 {
		return fmt.Errorf("Expected token missing after `%v`", e.tokens[0])
	}
	e.tokens = e.tokens[1:len(e.tokens)]
	return nil
}

func precedence(arg interface{}) int {
	op, _ := toOperator(arg)
	switch op {
	case And:
		return 2
	case Or:
		return 1
	case Not:
		return 3
	}
	return 0
}

func toOperator(arg interface{}) (string, error) {
	op, ok := arg.(string)
	if !ok {
		return "", fmt.Errorf("`%v` operator is not of type string", arg)
	}
	return op, nil
}

func isBinaryOperator(arg interface{}) (bool, error) {
	op, ok := arg.(string)
	if !ok {
		return false, nil
	}
	switch op {
	case And:
		return true, nil
	case Or:
		return true, nil
	case Not:
		return false, nil
	}
	return false, fmt.Errorf("`%v` operator not recognized", op)
}

func isUnaryOperator(arg interface{}) (bool, error) {
	op, ok := arg.(string)
	if !ok {
		return false, nil
	}
	switch op {
	case Not:
		return true, nil
	case And:
		return false, nil
	case Or:
		return false, nil
	}
	return false, fmt.Errorf("`%v` operator not recognized", op)
}

func (e *expression) parseExpressionR(lhs prism.Query, min int) (prism.Query, error) {

	var err error
	var op string
	var rhs prism.Query
	var lookahead interface{}
	var isBinary, isUnary bool

	lookahead = e.peek()

	isBinary, err = isBinaryOperator(lookahead)
	if err != nil {
		return nil, err
	}

	for isBinary && precedence(lookahead) >= min {

		op, err = toOperator(lookahead)
		if err != nil {
			return nil, err
		}

		err = e.advance()
		if err != nil {
			return nil, err
		}

		rhs, err = e.popOperand()
		if err != nil {
			return nil, err
		}

		lookahead = e.peek()

		isBinary, err = isBinaryOperator(lookahead)
		if err != nil {
			return nil, err
		}

		isUnary, err = isUnaryOperator(lookahead)
		if err != nil {
			return nil, err
		}

		for (isBinary && precedence(lookahead) > precedence(op)) ||
			(isUnary && precedence(lookahead) == precedence(op)) {
			rhs, err = e.parseExpressionR(rhs, precedence(lookahead))
			if err != nil {
				return nil, err
			}
			lookahead = e.peek()
		}
		lhs = &BinaryExpression{
			Left:  lhs,
			Op:    op,
			Right: rhs,
		}
	}
	return lhs, nil
}

func (e *expression) parse() (prism.Query, error) {
	lhs, err := e.popOperand()
	if err != nil {
		return nil, err
	}
	query, err := e.parseExpressionR(lhs, 0)
	if err != nil {
		return nil, err
	}
	return query, nil
}
