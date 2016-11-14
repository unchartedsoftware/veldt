package query

import (
	"fmt"
	"strings"

	"github.com/unchartedsoftware/prism/util/json"
)

// Validator parses a JSON query expression into its typed format. It
// ensure all types are correct and that the syntax is valid.
type Validator struct {
	json.Validator
}

// NewValidator instantiates and returns a new query expression object.
func NewValidator() *Validator {
	v := &Validator{}
	v.Output = make([]string, 0)
	return v
}

// Validate returns the instantiated runtime format of the provideJSON.
func (v *Validator) Validate(arg interface{}) (interface{}, error) {
	exp := v.validateToken(arg, 0)
	err := v.Error()
	if err != nil {
		return nil, err
	}
	return exp, nil
}

func (v *Validator) formatVal(val interface{}) string {
	str, ok := val.(string)
	if ok {
		return fmt.Sprintf("\"%s\"", str)
	}
	arr, ok := val.([]interface{})
	if ok {
		vals := make([]string, len(arr))
		for i, sub := range arr {
			vals[i] = v.formatVal(sub)
		}
		return fmt.Sprintf("[ %s ]", strings.Join(vals, ", "))
	}
	return fmt.Sprintf("%v", val)
}

func (v *Validator) getIDAndParams(query map[string]interface{}) (string, map[string]interface{}, bool) {
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

func (v *Validator) validateParams(id string, params map[string]interface{}, indent int, err error) {
	// open bracket
	v.Buffer("{", indent)
	v.Buffer(fmt.Sprintf("\"%s\": {", id), indent+1)
	// if error, start
	if err != nil {
		v.StartError(fmt.Sprintf("%v", err), indent+2)
	}
	// values
	for key, val := range params {
		v.Buffer(fmt.Sprintf("\"%s\": %s", key, v.formatVal(val)), indent+2)
	}
	// if error, end
	if err != nil {
		v.EndError()
	}
	// close bracket
	v.Buffer("}", indent+1)
	v.Buffer("}", indent)
}

func (v *Validator) validateQuery(arg map[string]interface{}, indent int) interface{} {
	// get query id and params
	id, params, ok := v.getIDAndParams(arg)
	if !ok {
		v.StartError("Empty query object", indent)
		v.Buffer("{", indent)
		v.Buffer("}", indent)
		v.EndError()
		return nil
	}

	//
	// TODO: FIX THIS
	// query, err := CreateQuery(id, params)
	// if err != nil {
	// 	v.validateParams(id, params, indent, err)
	// 	return nil
	// }
	// TODO: FIX THIS
	//

	v.validateParams(id, params, indent, nil)
	//return query
	return nil
}

func (v *Validator) validateOperator(op string, indent int) interface{} {
	if !IsBoolOperator(op) {
		v.StartError("Invalid operator", indent)
		v.Buffer(fmt.Sprintf("\"%v\"", op), indent)
		v.EndError()
		return nil
	}
	v.Buffer(fmt.Sprintf("\"%s\"", op), indent)
	return op
}

func (v *Validator) validateExpression(exp []interface{}, indent int) interface{} {
	// open paren
	v.Buffer("[", indent)
	// track last token to ensure next is valid
	var last interface{}
	// for each component
	for i, sub := range exp {
		// next line
		if last != nil {
			if !nextTokenIsValid(last, sub) {
				v.StartError("Unexpected token", indent+1)
				v.validateToken(sub, indent+1)
				v.EndError()
				last = sub
				continue
			}
		}
		exp[i] = v.validateToken(sub, indent+1)
		last = sub
	}
	// close paren
	v.Buffer("]", indent)
	return exp
}

func (v *Validator) validateToken(arg interface{}, indent int) interface{} {
	// expression
	exp, ok := arg.([]interface{})
	if ok {
		return v.validateExpression(exp, indent)
	}
	// query
	query, ok := arg.(map[string]interface{})
	if ok {
		return v.validateQuery(query, indent)
	}
	// operator
	op, ok := arg.(string)
	if ok {
		return v.validateOperator(op, indent)
	}
	// err
	v.StartError("Unrecognized symbol", indent)
	v.Buffer(fmt.Sprintf("\"%v\"", arg), indent)
	v.EndError()
	return arg
}

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
