package veldt

import (
	"fmt"
)

type expressionParser struct {
	pipeline *Pipeline
}

func newExpressionParser(pipeline *Pipeline) *expressionParser {
	return &expressionParser{
		pipeline: pipeline,
	}
}

func (p *expressionParser) Parse(arg interface{}) (Query, error) {
	// parse into correct AST
	return parseToken(p.pipeline, arg)
}

func parseToken(pipeline *Pipeline, token interface{}) (Query, error) {
	// check if token is an expression
	exp, ok := token.([]interface{})
	if ok {
		// is expression, recursively parse it
		return newExpression(pipeline, exp).parse()
	}
	// is query, parse it directly
	query, ok := token.(Query)
	if !ok {
		return nil, fmt.Errorf("`%v` token is unrecognized", token)
	}
	return query, nil
}

// expression parses the runtime query expression into it's runtime AST tree.
type expression struct {
	pipeline *Pipeline
	tokens   []interface{}
}

func newExpression(pipeline *Pipeline, arr []interface{}) *expression {
	return &expression{
		pipeline: pipeline,
		tokens:   arr,
	}
}

func (e *expression) parse() (Query, error) {
	lhs, err := e.popOperand()
	if err != nil {
		return nil, err
	}
	query, err := e.parseExpression(lhs, 0)
	if err != nil {
		return nil, err
	}
	return query, nil
}

func (e *expression) pop() (interface{}, error) {
	if len(e.tokens) == 0 {
		return nil, fmt.Errorf("expected operand missing")
	}
	token := e.tokens[0]
	e.tokens = e.tokens[1:len(e.tokens)]
	return token, nil
}

func (e *expression) popOperand() (Query, error) {
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
		query, err := parseToken(e.pipeline, next)
		if err != nil {
			return nil, err
		}
		// get unary expression
		unary, err := e.pipeline.GetUnary()
		if err != nil {
			return nil, err
		}
		// populate it
		err = unary.Parse(map[string]interface{}{
			"op":    op,
			"query": query,
		})
		if err != nil {
			return nil, err
		}
		return unary, nil
	}

	// parse token
	return parseToken(e.pipeline, token)
}

func (e *expression) peek() interface{} {
	if len(e.tokens) == 0 {
		return nil
	}
	return e.tokens[0]
}

func (e *expression) advance() error {
	if len(e.tokens) < 2 {
		return fmt.Errorf("expected token missing after `%v`", e.tokens[0])
	}
	e.tokens = e.tokens[1:len(e.tokens)]
	return nil
}

func (e *expression) parseExpression(lhs Query, min int) (Query, error) {

	var err error
	var op string
	var rhs Query
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
			rhs, err = e.parseExpression(rhs, precedence(lookahead))
			if err != nil {
				return nil, err
			}
			lookahead = e.peek()
		}

		// get binary expression
		binary, err := e.pipeline.GetBinary()
		if err != nil {
			return nil, err
		}
		// populate it
		err = binary.Parse(map[string]interface{}{
			"left":  lhs,
			"op":    op,
			"right": rhs,
		})
		if err != nil {
			return nil, err
		}
		lhs = binary
	}
	return lhs, nil
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
