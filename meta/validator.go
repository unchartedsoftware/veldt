package tile

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
func (v *Validator) Validate(arg interface{}) (prism.Tile, error) {
	tile := v.validateToken(arg, 0)
	err := v.Error()
	if err != nil {
		return nil, err
	}
	return tile, nil
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

func (v *Validator) formatParams(id string, params map[string]interface{}, indent int, err error) {
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

func (v *Validator) validateMeta(arg map[string]interface{}, indent int) prism.Tile {
	// get generator id and params
	id, params, ok := v.getIDAndParams(arg)
	if !ok {
		v.StartError("Empty meta object", indent)
		v.Buffer("{", indent)
		v.Buffer("}", indent)
		v.EndError()
		return nil
	}

	//
	// TODO: FIX THIS
	tile, err := CreateMeta(id, params)
	if err != nil {
		v.formatParams(id, params, indent, err)
		return nil
	}
	// TODO: FIX THIS
	//

	v.formatParams(id, params, indent, nil)
	return tile
}

func (v *Validator) validateToken(arg interface{}, indent int) prism.Tile {
	// tile
	query, ok := arg.(map[string]interface{})
	if !ok {
		// err
		v.StartError("Unrecognized symbol", indent)
		v.Buffer(fmt.Sprintf("\"%v\"", arg), indent)
		v.EndError()
		return nil
	}
	return v.validateMeta(query, indent)
}
