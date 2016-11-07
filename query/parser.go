package query

import (
	"encoding/json"
	"fmt"
)

// Parser parses the runtime query expression into it's runtime AST tree.
type Parser struct {
	tokens []interface{}
}

// NewParser instantiates and returns a new expression parser object.
func NewParser(arr []interface{}) *Parser {
	return &Parser{
		tokens: arr,
	}
}

// Parse parses the query payload into the query AST.
func Parse(bytes []byte) (Query, error) {
	// unmarshal the query
	var token interface{}
	err := json.Unmarshal(bytes, &token)
	if err != nil {
		return nil, fmt.Errorf("Invalid JSON: %v", err)
	}
	// validate the JSON into ot's runtime query components
	validator := NewValidator()
	exp, err := validator.Validate(token)
	if err != nil {
		return nil, err
	}
	// parse into correct AST
	return parseToken(exp)
}

func (t *Parser) pop() (interface{}, error) {
	if len(t.tokens) == 0 {
		return nil, fmt.Errorf("Expected operand missing")
	}
	token := t.tokens[0]
	t.tokens = t.tokens[1:len(t.tokens)]
	return token, nil
}

func (t *Parser) popOperand() (Query, error) {
	// pops the next operand
	//     cases to consider:
	//         - a) unary operator -> expression
	//         - b) unary operator -> query
	//         - c) expression
	//         - d) query

	// pop next token
	token, err := t.pop()
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
		next, err := t.pop()
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
			Query: query,
		}, nil
	}

	// parse token
	return parseToken(token)
}

func (t *Parser) peek() interface{} {
	if len(t.tokens) == 0 {
		return nil
	}
	return t.tokens[0]
}

func (t *Parser) advance() error {
	if len(t.tokens) < 2 {
		return fmt.Errorf("Expected token missing after `%v`", t.tokens[0])
	}
	t.tokens = t.tokens[1:len(t.tokens)]
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

func (t *Parser) parseExpressionR(lhs Query, min int) (Query, error) {

	var err error
	var op string
	var rhs Query
	var lookahead interface{}
	var isBinary, isUnary bool

	lookahead = t.peek()

	isBinary, err = isBinaryOperator(lookahead)
	if err != nil {
		return nil, err
	}

	for isBinary && precedence(lookahead) >= min {

		op, err = toOperator(lookahead)
		if err != nil {
			return nil, err
		}

		err = t.advance()
		if err != nil {
			return nil, err
		}

		rhs, err = t.popOperand()
		if err != nil {
			return nil, err
		}

		lookahead = t.peek()

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
			rhs, err = t.parseExpressionR(rhs, precedence(lookahead))
			if err != nil {
				return nil, err
			}
			lookahead = t.peek()
		}
		lhs = &BinaryExpression{
			Left:  lhs,
			Op:    op,
			Right: rhs,
		}
	}
	return lhs, nil
}

// Parse parses the expression.
func (t *Parser) parse() (Query, error) {
	lhs, err := t.popOperand()
	if err != nil {
		return nil, err
	}
	query, err := t.parseExpressionR(lhs, 0)
	if err != nil {
		return nil, err
	}
	return query, nil
}

func parseExpression(args []interface{}) (Query, error) {
	exp := NewParser(args)
	return exp.parse()
}

func parseToken(token interface{}) (Query, error) {
	// check if token is an expression
	exp, ok := token.([]interface{})
	if ok {
		// is expression, recursively parse it
		return parseExpression(exp)
	}
	// is query, parse it directly
	query, ok := token.(Query)
	if !ok {
		return nil, fmt.Errorf("`%v` token is unrecognized", token)
	}
	return query, nil
}
