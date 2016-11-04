package query

import (
	"fmt"
	"strings"
)

const (
	indentor     = "    "
	redColor     = "\033[31m"
	defaultColor = "\033[39m"
)

func getTokenType(token interface{}) string {
	_, ok := token.([]interface{})
	if ok {
		return "exp"
	}
	op, ok := token.(string)
	if ok {
		if IsBinaryOperator(op) {
			return "binary"
		} else if IsUnaryOperator(op) {
			return "unary"
		} else {
			return "unrecognized"
		}
	}
	_, ok = token.(map[string]interface{})
	if ok {
		return "query"
	}
	return "unrecognized"
}

func nextTokenIsValid(c interface{}, n interface{}) bool {
	current := getTokenType(c)
	next := getTokenType(n)
	if current == "unrecognized" || next == "unrecognized" {
		// NOTE: consider unrecognized tokens as valid to allow the parsing to
		// continue correctly
		return true
	}
	switch current {
	case "exp":
		return next == "binary"
	case "query":
		return next == "binary"
	case "binary":
		return next == "unary" || next == "query" || next == "exp"
	case "unary":
		return next == "query" || next == "exp"
	}
	return false
}

// ExpressionValidator parses a JSON query expression into its typed format. It
// ensure all types are correct and that the syntax is valid.
type ExpressionValidator struct {
	output         []string
	errStartIndex  int
	errEndIndex    int
	errIndent      string
	errHeaderIndex int
	errMsg         string
	err            bool
}

// NewExpressionValidator instantiates and returns a new query expression object.
func NewExpressionValidator() *ExpressionValidator {
	return &ExpressionValidator{
		output: make([]string, 0),
	}
}

func getIndent(indent int) string {
	var strs []string
	for i := 0; i < indent; i++ {
		strs = append(strs, indentor)
	}
	return strings.Join(strs, "")
}

func (q *ExpressionValidator) buffer(line string) {
	q.output = append(q.output, line)
}

func (q *ExpressionValidator) size() int {
	return len(q.output)
}

// Error handling

func (q *ExpressionValidator) startError(msg string, indent string) {
	q.err = true
	q.errHeaderIndex = q.size()
	q.errStartIndex = q.size() + 1
	q.errIndent = indent
	q.errMsg = msg
	q.buffer("")
}

func (q *ExpressionValidator) getErrAnnotations(width int, char string) string {
	arr := make([]string, width)
	for i := 0; i < width; i++ {
		arr[i] = char
	}
	return strings.Join(arr, "")
}

func (q *ExpressionValidator) getErrHeader(width int) string {
	return fmt.Sprintf("%s%s%s%s",
		redColor,
		q.errIndent,
		q.getErrAnnotations(width, "v"),
		defaultColor)
}

func (q *ExpressionValidator) getErrFooter(width int) string {
	return fmt.Sprintf("%s%s%s Error: %s%s",
		redColor,
		q.errIndent,
		q.getErrAnnotations(width, "^"),
		q.errMsg,
		defaultColor)
}

func (q *ExpressionValidator) getErrWidth() int {
	maxWidth := 1
	for i := q.errStartIndex; i < q.errEndIndex; i++ {
		width := (len(q.output[i]) - len(q.errIndent))
		if width > maxWidth {
			maxWidth = width
		}
	}
	return maxWidth
}

func (q *ExpressionValidator) endError() {
	q.errEndIndex = q.size()
	width := q.getErrWidth()
	header := q.getErrHeader(width)
	footer := q.getErrFooter(width)
	q.output[q.errHeaderIndex] = header
	q.buffer(footer)
}

func (q *ExpressionValidator) formatVal(val interface{}) string {
	str, ok := val.(string)
	if ok {
		return fmt.Sprintf("\"%s\"", str)
	}
	arr, ok := val.([]interface{})
	if ok {
		vals := make([]string, len(arr))
		for i, sub := range arr {
			vals[i] = q.formatVal(sub)
		}
		return fmt.Sprintf("[ %s ]", strings.Join(vals, ", "))
	}
	return fmt.Sprintf("%v", val)
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

func (q *ExpressionValidator) formatParams(id string, params map[string]interface{}, n int, err error) {
	indent := getIndent(n)
	idIndent := getIndent(n + 1)
	paramIndent := getIndent(n + 2)
	// open bracket
	q.buffer(fmt.Sprintf("%s{", indent))
	q.buffer(fmt.Sprintf("%s\"%s\": {", idIndent, id))
	// if error, start
	if err != nil {
		q.startError(fmt.Sprintf("%v", err), paramIndent)
	}
	// values
	for key, val := range params {
		q.buffer(fmt.Sprintf("%s\"%s\": %s", paramIndent, key, q.formatVal(val)))
	}
	// if error, end
	if err != nil {
		q.endError()
	}
	// close bracket
	q.buffer(fmt.Sprintf("%s}", idIndent))
	q.buffer(fmt.Sprintf("%s}", indent))
}

func (q *ExpressionValidator) formatQuery(arg map[string]interface{}, n int) interface{} {
	indent := getIndent(n)
	// pattern match for base queries
	id, params, ok := getQueryAndKey(arg)
	if !ok {
		q.startError("Empty query object", indent)
		q.buffer(fmt.Sprintf("%s{", indent))
		q.buffer(fmt.Sprintf("%s}", indent))
		q.endError()
		return nil
	}
	query, err := GetQuery(id, params)
	if err != nil {
		q.formatParams(id, params, n, err)
		return nil
	}
	q.formatParams(id, params, n, nil)
	return query
}

func (q *ExpressionValidator) formatOperator(op string, n int) interface{} {
	indent := getIndent(n)
	if !IsBoolOperator(op) {
		q.startError("Invalid operator", indent)
		q.buffer(fmt.Sprintf("%s\"%v\"", indent, op))
		q.endError()
		return nil
	}
	q.buffer(fmt.Sprintf("%s\"%s\"", indent, op))
	return op
}

func (q *ExpressionValidator) formatExpression(exp []interface{}, n int) interface{} {
	indent := getIndent(n)
	// open paren
	q.buffer(fmt.Sprintf("%s[", indent))
	// track last token to ensure next is valid
	var last interface{}
	// for each component
	for i, sub := range exp {
		// next line
		if last != nil {
			if !nextTokenIsValid(last, sub) {
				q.startError("Unexpected token", getIndent(n+1))
				q.formatToken(sub, n+1)
				q.endError()
				last = sub
				continue
			}
		}
		exp[i] = q.formatToken(sub, n+1)
		last = sub
	}
	// close paren
	q.buffer(fmt.Sprintf("%s]", indent))
	return exp
}

func (q *ExpressionValidator) formatToken(arg interface{}, n int) interface{} {
	// expression
	exp, ok := arg.([]interface{})
	if ok {
		return q.formatExpression(exp, n)
	}
	// query
	query, ok := arg.(map[string]interface{})
	if ok {
		return q.formatQuery(query, n)
	}
	// operator
	op, ok := arg.(string)
	if ok {
		return q.formatOperator(op, n)
	}
	// err
	indent := getIndent(n)
	q.startError("Unrecognized symbol", indent)
	q.buffer(fmt.Sprintf("%s\"%v\"", indent, arg))
	q.endError()
	return arg
}

// func (q *ExpressionValidator) formatExpression(arg interface{}, n int) interface{} {
// 	indent := getIndent(n)
// 	arr, ok := arg.([]interface{})
// 	if ok {
// 		// open paren
// 		q.buffer(fmt.Sprintf("%s[", indent))
// 		var last interface{}
// 		// for each component
// 		for i, sub := range arr {
// 			// next line
// 			if last != nil {
// 				if !nextTokenIsValid(last, sub) {
// 					q.startError("Unexpected token", getIndent(n+1))
// 					q.formatExpression(sub, n+1)
// 					q.endError()
// 					last = sub
// 					continue
// 				}
// 			}
// 			arr[i] = q.formatExpression(sub, n+1)
// 			last = sub
// 		}
// 		// close paren
// 		q.buffer(fmt.Sprintf("%s]", indent))
// 		return arg
// 	}
// 	arg = q.formatToken(arg, indent)
// 	return arg
// }

func (q *ExpressionValidator) format(arg interface{}) (interface{}, error) {
	exp := q.formatToken(arg, 0)
	if q.err {
		return nil, fmt.Errorf(strings.Join(q.output, "\n"))
	}
	return exp, nil
}

func prePass(arg interface{}) (interface{}, error) {
	q := NewExpressionValidator()
	return q.format(arg)
}
