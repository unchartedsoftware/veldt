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
		found = true
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

func (t *expression) success(token interface{}) {
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

func isBinaryOperator(arg interface{}) (bool, error) {
	str, ok := arg.(string)
	if !ok {
		return false, nil
	}
	switch str {
	case And:
		return true, nil
	case Or:
		return true, nil
	case Not:
		return false, nil
	}
	return false, fmt.Errorf("Operator not recognized")
}

func isUnaryOperator(arg interface{}) (bool, error) {
	str, ok := arg.(string)
	if !ok {
		return false, nil
	}
	switch str {
	case Not:
		return true, nil
	case And:
		return false, nil
	case Or:
		return false, nil
	}
	return false, fmt.Errorf("Operator not recognized")
}

func (t *expression) parseExpressionR(lhs Query, min int) (Query, error) {

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

func (t *expression) Parse() (Query, error) {
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

/*
func getIndent(indent int) string {
	strs := make([]string, 0)
	for i := 0; i < indent; i++ {
		strs = append(strs, indentor)
	}
	return strings.Join(strs, "")
}

var (
	indentor = "    "
	output   []string

	errStartIndex  = 0
	errEndIndex    = 0
	errWidth       = 0
	errIndent      = ""
	errHeaderIndex = 0
	errMsg         = ""
	redColor       = "\033[31m"
	defaultColor   = "\033[39m"
)

// Two passes

// Pre-pass
//  - log to string slice
//  - instantiate all queries, on errors, add appropriate annotations in string

// Order of Operations pass
//  - take the instantiate tree, and add the correct logical flow

func formatVal(val interface{}) string {
	str, ok := val.(string)
	if ok {
		return fmt.Sprintf("\"%s\"", str)
	}
	arr, ok := val.([]interface{})
	if ok {
		vals := make([]string, len(arr))
		for i, sub := range arr {
			vals[i] = formatVal(sub)
		}
		return fmt.Sprintf("[ %s ]", strings.Join(vals, ", "))
	}

	return fmt.Sprintf("%v", val)
}

func startError(msg string, indent string) {
	errHeaderIndex = len(output)
	errStartIndex = len(output) + 1
	errIndent = indent
	errMsg = msg
	output = append(output, "")
}

func endError() {
	errEndIndex = len(output)
	maxWidth := 0
	for i := errStartIndex; i < errEndIndex; i++ {
		width := (len(output[i]) - len(errIndent))
		if width > maxWidth {
			maxWidth = width
		}
	}
	errHeader := make([]string, maxWidth)
	errFooter := make([]string, maxWidth)
	for i := 0; i < maxWidth; i++ {
		errHeader[i] = "v"
		errFooter[i] = "^"
	}
	output[errHeaderIndex] = fmt.Sprintf("%s%s%s%s", redColor, errIndent, strings.Join(errHeader, ""), defaultColor)
	output = append(output, fmt.Sprintf("%s%s%s Error: %s%s", redColor, errIndent, strings.Join(errFooter, ""), errMsg, defaultColor))
}

func formatParams(id string, params map[string]interface{}, indent string, err error) {
	idIndent := fmt.Sprintf("%s%s", indent, indentor)
	paramIndent := fmt.Sprintf("%s%s%s", indent, indentor, indentor)

	// open bracket
	output = append(output, fmt.Sprintf("%s{", indent))
	output = append(output, fmt.Sprintf("%s\"%s\": {", idIndent, id))
	if err != nil {
		startError(fmt.Sprintf("%v", err,), paramIndent)
	}
	// values
	for key, val := range params {
		// next line
		output = append(output, fmt.Sprintf("%s\"%s\": %s",
			paramIndent,
			key,
			formatVal(val)))
	}
	if err != nil {
		endError()
	}
	// close bracket
	output = append(output, fmt.Sprintf("%s}", idIndent))
	output = append(output, fmt.Sprintf("%s}", indent))
}

func formatQuery(arg map[string]interface{}, indent string) {
	// pattern match for base queries
	id, params, ok := getQueryAndKey(arg)
	if !ok {
		startError("Empty query object", indent)
		output = append(output, fmt.Sprintf("%s%v", indent, arg))
		endError()
		return
	}
	_, err := GetQuery(id, params)
	if err != nil {
		// TODO: move this in deeper so that it doesn't wrap thje outer query object and only the params
		//startError(fmt.Sprintf("%v", err), indent)
		formatParams(id, params, indent, err)
		//endError()
		return
	}

	formatParams(id, params, indent, nil)
}

func isOp(op string) bool {
	switch op {
	case And:
		return true
	case Or:
		return true
	case Not:
		return true
	}
	return false
}

func formatOperator(arg interface{}, indent string) {
	op, ok := arg.(string)
	if !ok {
		startError("Unrecognized symbol", indent)
		output = append(output, fmt.Sprintf("%s\"%v\"", indent, arg))
		endError()
		return
	}
	if !isOp(op) {
		startError("Invalid operator", indent)
		output = append(output, fmt.Sprintf("%s\"%v\"", indent, arg))
		endError()
		return
	}
	output = append(output, fmt.Sprintf("%s\"%s\"", indent, op))
}

func formatToken(arg interface{}, indent string) {
	obj, ok := arg.(map[string]interface{})
	if ok {
		formatQuery(obj, indent)
	} else {
		formatOperator(arg, indent)
	}
}

func formatExpression(arg interface{}, n int) {
	indent := getIndent(n)
	arr, ok := arg.([]interface{})
	if ok {
		// open paren
		open := fmt.Sprintf("%s[", indent)
		output = append(output, open)
		// for each component
		for _, sub := range arr {
			// next line
			formatExpression(sub, n+1)
		}
		// close paren
		close := fmt.Sprintf("%s]", indent)
		output = append(output, close)
		return
	}
	formatToken(arg, indent)
}

func prePass(arg interface{}) string {
	output = make([]string, 0)
	formatExpression(arg, 0)
	return strings.Join(output, "\n")
}
*/

// Parse parses the query payload into the query AST.
func Parse(bytes []byte) (Query, error) {
	// unmarshal the query
	var token interface{}
	err := json.Unmarshal(bytes, &token)
	if err != nil {
		return nil, err
	}
	// run a pre-pass to check for valid query syntax
	_, err = prePass(token)
	if err != nil {
		return nil, err
	}
	// parse into correct AST
	return parseToken(token)
}
