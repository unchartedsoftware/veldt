package query

import (
	"encoding/json"
	"fmt"
)

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

func getQueryAndKey(query map[string]interface{}) (string, map[string]interface{}, bool) {
    var key string
    var value map[string]interface{}
    found := false
    for k, v := range query {
		val, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
        key = k
        value = val
        found  = true
		break
	}
    return key, value, found
}

func parseQuery(arg interface{}) (Query, error) {
	// pattern match for base queries

    query, ok := arg.(map[string]interface{})
    if !ok {
        return nil, fmt.Errorf("`%v` is not a recognised query format", arg)
    }

    id, params, ok := getQueryAndKey(query)
    if !ok {
        return nil, fmt.Errorf("Query `%v` is empty", query)
    }

    return GetQuery(id, params)
}

func parseToken(token interface{}) (Query, error) {
	// check if token is an expression
	exp, ok := token.([]interface{})
	if ok {
    	fmt.Println("EXP:")
        for _, e := range exp {
            fmt.Printf("\t%v,\n", e)
        }
        fmt.Println()
		// is expression, recursively parse it
		return parseExpression(exp)
	}
	fmt.Println("QUERY:")
	fmt.Printf("\t%v,\n\n", token)
	// is query, parse it directly
	return parseQuery(token)
}

type expression struct {
	tokens []interface{}
    parsed []interface{}
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
    t.parsed = append(t.parsed, token)
	t.tokens = t.tokens[1:len(t.tokens)]
	return token, nil
}


func (t *expression) success(token  interface{}) {
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

// func (t *expression) peek() (interface{}, bool) {
// 	if len(t.tokens) == 0 {
// 		return nil, fmt.Errorf("Expected token missing")
// 	}
// 	return t.tokens[0], nil
// }

func (t *expression) peek() interface{} {
	if len(t.tokens) == 0 {
		return nil
	}
	return t.tokens[0]
}

func (t *expression) advance() error {
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
		return true
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

	lookahead = t.peek()

	for isBinaryOperator(lookahead) && precedence(lookahead) >= min {

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

		for (isBinaryOperator(lookahead) && precedence(lookahead) > precedence(op)) ||
		 	(isUnaryOperator(lookahead) && precedence(lookahead) == precedence(op)) {
			rhs, err = t.parseExpressionR(rhs, precedence(lookahead))
			if err != nil {
				return nil, err
			}
			lookahead = t.peek()
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
        fmt.Println("Parsed:", t.parsed[0:len(t.parsed)-1])
        fmt.Println("Failed:", t.parsed[len(t.parsed)-1])
        fmt.Println("Remaining:", t.tokens)
		return nil, err
	}
	query, err := t.parseExpressionR(lhs, 0)
    if err != nil {
        fmt.Println("Parsed:", t.parsed[0:len(t.parsed)-1])
        fmt.Println("Failed:", t.parsed[len(t.parsed)-1])
        fmt.Println("Remaining:", t.tokens)
		return nil, err
	}
    return query, nil
}

// Parse parses the query payload into the query AST.
func Parse(bytes []byte) (Query, error) {
	var token interface{}
	err := json.Unmarshal(bytes, &token)
    if err != nil {
        return nil, err
    }
	return parseToken(token)
}
