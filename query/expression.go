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

// ExpressionParser parses a JSON query expression into its typed format.
type ExpressionParser struct {
	output         []string
	errStartIndex  int
	errEndIndex    int
	errWidth       int
	errIndent      string
	errHeaderIndex int
	errMsg         string
	err            bool
	expression     interface{}
}

// NewExpressionParser instantiates and returns a new query expression object.
func NewExpressionParser(arg interface{}) *ExpressionParser {
	return &ExpressionParser{
		output:     make([]string, 0),
		expression: arg,
	}
}

func getIndent(indent int) string {
	var strs []string
	for i := 0; i < indent; i++ {
		strs = append(strs, indentor)
	}
	return strings.Join(strs, "")
}

func incIndent(indent string) string {
	return fmt.Sprintf("%s%s", indent, indentor)
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

func (q *ExpressionParser) buffer(line string) {
	q.output = append(q.output, line)
}

func (q *ExpressionParser) size() int {
	return len(q.output)
}

// Error handling

func (q *ExpressionParser) startError(msg string, indent string) {
	q.err = true
	q.errHeaderIndex = q.size()
	q.errStartIndex = q.size() + 1
	q.errIndent = indent
	q.errMsg = msg
	q.buffer("")
}

func (q *ExpressionParser) getErrAnnotations(width int, char string) string {
	arr := make([]string, width)
	for i := 0; i < width; i++ {
		arr[i] = char
	}
	return strings.Join(arr, "")
}

func (q *ExpressionParser) getErrHeader(width int) string {
	return fmt.Sprintf("%s%s%s%s",
		redColor,
		q.errIndent,
		q.getErrAnnotations(width, "v"),
		defaultColor)
}

func (q *ExpressionParser) getErrFooter(width int) string {
	return fmt.Sprintf("%s%s%s Error: %s%s",
		redColor,
		q.errIndent,
		q.getErrAnnotations(width, "^"),
		q.errMsg,
		defaultColor)
}

func (q *ExpressionParser) getErrWidth() int {
	maxWidth := 0
	for i := q.errStartIndex; i < q.errEndIndex; i++ {
		width := (len(q.output[i]) - len(q.errIndent))
		if width > maxWidth {
			maxWidth = width
		}
	}
	return maxWidth
}

func (q *ExpressionParser) endError() {
	q.errEndIndex = q.size()
	width := q.getErrWidth()
	header := q.getErrHeader(width)
	footer := q.getErrFooter(width)
	q.output[q.errHeaderIndex] = header
	q.buffer(footer)
}

func (q *ExpressionParser) formatVal(val interface{}) string {
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

func (q *ExpressionParser) formatParams(id string, params map[string]interface{}, indent string, err error) {
	idIndent := incIndent(indent)
	paramIndent := incIndent(idIndent)
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

func (q *ExpressionParser) formatQuery(arg map[string]interface{}, indent string) {
	// pattern match for base queries
	id, params, ok := getQueryAndKey(arg)
	if !ok {
		q.startError("Empty query object", indent)
		q.buffer(fmt.Sprintf("%s%v", indent, arg))
		q.endError()
		return
	}
	_, err := GetQuery(id, params)
	if err != nil {
		q.formatParams(id, params, indent, err)
		return
	}
	q.formatParams(id, params, indent, nil)
}

func (q *ExpressionParser) formatOperator(arg interface{}, indent string) {
	op, ok := arg.(string)
	if !ok {
		q.startError("Unrecognized symbol", indent)
		q.buffer(fmt.Sprintf("%s\"%v\"", indent, arg))
		q.endError()
		return
	}
	if !isOp(op) {
		q.startError("Invalid operator", indent)
		q.buffer(fmt.Sprintf("%s\"%v\"", indent, arg))
		q.endError()
		return
	}
	q.buffer(fmt.Sprintf("%s\"%s\"", indent, op))
}

func (q *ExpressionParser) formatToken(arg interface{}, indent string) {
	obj, ok := arg.(map[string]interface{})
	if ok {
		q.formatQuery(obj, indent)
	} else {
		q.formatOperator(arg, indent)
	}
}

func (q *ExpressionParser) formatExpression(arg interface{}, n int) {
	indent := getIndent(n)
	arr, ok := arg.([]interface{})
	if ok {
		// open paren
		q.buffer(fmt.Sprintf("%s[", indent))
		// for each component
		for _, sub := range arr {
			// next line
			q.formatExpression(sub, n+1)
		}
		// close paren
		q.buffer(fmt.Sprintf("%s]", indent))
		return
	}
	q.formatToken(arg, indent)
}

func prePass(arg interface{}) (interface{}, error) {
	q := NewExpressionParser(arg)
	q.formatExpression(arg, 0)
	if q.err {
		return nil, fmt.Errorf(strings.Join(q.output, "\n"))
	}
	return q.expression, nil
}
