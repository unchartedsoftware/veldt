package query

import (
	"encoding/json"
	"fmt"
)
/*
[
	{
		"equals": {
			"field": "surname",
			"value": "bethune"
		}
	},
	"AND"
	{
		"in": {
			"field": "hashtags",
			"values": ["dank", "420", "nugz"]
		}
	},
	"AND"
	[
		"NOT",
		{
			"exists": {
				"field": "location"
			}
		},
		"OR",
		{
			"range": {
				"field": "data",
				"lt": 123234532452
			}
		}
	]
]
*/

/*
func ParseUnaryOp(op interface{}) (string, error) {
	opString, ok := op.(string)
	if !ok {
		return "", fmt.Errorf("Unaryexpression operator `%v` is not of type `string`",
			op)
	}
	switch opString {
	case Not:
		return Not, nil
	default:
		return "", fmt.Errorf("Unaryexpression operator `%s` is not recognized",
			opString)
	}
}

func ParseBinaryOp(op interface{}) (string, error) {
	opString, ok := op.(string)
	if !ok {
		return "", fmt.Errorf("Binaryexpression operator `%v` is not of type `string`",
			op)
	}
	switch opString {
	case And:
		return And, nil
	case Or:
		return Or, nil
	default:
		return "", fmt.Errorf("Binarryexpression operator `%s` is not recognized",
			opString)
	}
}
*/

func parseExpression(args []interface{}) (Query, error) {
	exp := newExpression(args)
	return exp.Parse()
}

func parseQuery(query interface{}) (Query, error) {
	// pattern match for base queries
	return nil, nil
}

func parseToken(token interface{}) (Query, error) {
	// check if token is an expression
	exp, ok := token.([]interface{})
	if ok {
		// is expression, recursively parse it
		return parseExpression(exp)
	}
	// is query, parse it directly
	return parseQuery(exp)
}

type expression struct {
	tokens []interface{}
}

func newExpression(arr []interface{}) *expression {
	return &expression{
		tokens: arr,
	}
}

func (t *expression) pop() (interface{}, error) {
	if len(t.tokens) == 0 {
		return nil, fmt.Errorf("Expected operand missing")
	}
	token := t.tokens[0]
	t.tokens = t.tokens[1:len(t.tokens)]
	return token, nil
}


func (t *expression) popOperand() (Query, error) {
	// Pops the next operand
	// There are 4 cases to consider:
	//     - a) unary operator -> expression
	//     - b) unary operator -> query
	//     - c) expression
	//     - d) query

	// pop next token
	token, err := t.pop()
	if err != nil {
		return nil, err
	}

	// see if it is a unary operator
	op, ok := token.(string)
	if ok && isUnaryOperator(token) {
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
			Op: op,
			Query: query,
		}, nil
	}

	// parse token
	return parseToken(token)
}

func (t *expression) peek() (interface{}, error) {
	if len(t.tokens) == 0 {
		return nil, fmt.Errorf("Expected token missing")
	}
	return t.tokens[0], nil
}

func (t *expression) advance() error {
	if len(t.tokens) < 2 {
		return fmt.Errorf("Expected token missing after `%v`", t.tokens[0])
	}
	t.tokens = t.tokens[1:len(t.tokens)]
	return nil
}

func precendence(arg interface{}) int {
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
		return "", fmt.Errorf("Operator `%v` not recognized", arg)
	}
	return op, nil
}

func isBinaryOperator(arg interface{}) bool {
	str, ok := arg.(string)
	if !ok {
		return false
	}
	switch str {
	case And:
	case Or:
		return true
	}
	return false
}

func isUnaryOperator(arg interface{}) bool {
	str, ok := arg.(string)
	if !ok {
		return false
	}
	switch str {
	case Not:
		return true
	}
	return false
}

func (t *expression) parseExpressionR(lhs Query, min int) (Query, error) {

	var err error
	var op string
	var rhs Query
	var lookahead interface{}

	lookahead, err = t.peek()
	if err != nil {
		return nil, err
	}

	for isBinaryOperator(lookahead) && precendence(lookahead) >= min {

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

		lookahead, err = t.peek()
		if err != nil {
			return nil, err
		}

		for (isBinaryOperator(lookahead) && precendence(lookahead) > precendence(op)) ||
		 	(isUnaryOperator(lookahead) && precendence(lookahead) == precendence(op)) {
			rhs, err = t.parseExpressionR(rhs, precendence(lookahead))
			if err != nil {
				return nil, err
			}
			lookahead, err = t.peek()
			if err != nil {
				return nil, err
			}
		}
		lhs = &BinaryExpression{
			Left: lhs,
			Op: op,
			Right: rhs,
		}
	}
	return lhs, nil
}

func (t *expression) Parse() (Query, error) {
	lhs, err := t.popOperand()
	if err != nil {
		return nil, err
	}
	return t.parseExpressionR(lhs, 0)
}

// Parse parses the query payload into the query AST.
func Parse(bytes []byte) (Query, error) {
	var token interface{}
	json.Unmarshal(bytes, &token)
	return parseToken(token)
}
