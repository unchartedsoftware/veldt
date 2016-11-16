package json

import (
	"fmt"
	"strings"
)

const (
	indentor     = "    "
	redColor     = "\033[31m"
	defaultColor = "\033[39m"
)

// Validator parses a JSON query expression into its typed format. It
// ensure all types are correct and that the syntax is valid.
type Validator struct {
	Output         []string
	errStartIndex  int
	errEndIndex    int
	errIndent      int
	errHeaderIndex int
	errFooterIndex int
	errMsg         string
	err            bool
}

// Validate returns the instantiated runtime format of the provideJSON.
func (v *Validator) Validate(arg interface{}) (interface{}, error) {
	return nil, fmt.Errorf("Not implemented")
}

// Buffer adds a string at the appropriate indent to the output buffer.
func (v *Validator) Buffer(str string, indent int) {
	line := fmt.Sprintf("%s%s", v.getIndent(indent), str)
	v.Output = append(v.Output, line)
}

// Size returns the length of the current output buffer.
func (v *Validator) Size() int {
	return len(v.Output)
}

// HasError begins wrapping an error portion of the output buffer.
func (v *Validator) Error() error {
	if v.err {
		return fmt.Errorf(strings.Join(v.Output, "\n"))
	}
	return nil
}

// StartError begins wrapping an error portion of the output buffer.
func (v *Validator) StartError(msg string, indent int) {
	v.err = true
	v.errHeaderIndex = v.Size()
	v.errStartIndex = v.Size() + 1
	v.errIndent = indent
	v.errMsg = msg
	v.Buffer("", 0) // header line
}

// EndError ends wrapping an error portion of the output buffer.
func (v *Validator) EndError() {
	v.errEndIndex = v.Size()
	v.Buffer("", 0) // footer line
	width := v.getErrWidth()
	header := v.getErrHeader(width)
	footer := v.getErrFooter(width)
	v.Output[v.errHeaderIndex] = header
	v.Output[v.errEndIndex] = footer
}

func (v *Validator) getErrAnnotations(width int, char string) string {
	arr := make([]string, width)
	for i := 0; i < width; i++ {
		arr[i] = char
	}
	return strings.Join(arr, "")
}

func (v *Validator) getErrHeader(width int) string {
	return fmt.Sprintf("%s%s%s%s",
		redColor,
		v.getIndent(v.errIndent),
		v.getErrAnnotations(width, "v"),
		defaultColor)
}

func (v *Validator) getErrFooter(width int) string {
	return fmt.Sprintf("%s%s%s Error: %s%s",
		redColor,
		v.getIndent(v.errIndent),
		v.getErrAnnotations(width, "^"),
		v.errMsg,
		defaultColor)
}

func (v *Validator) getErrWidth() int {
	maxWidth := 1
	for i := v.errStartIndex; i < v.errEndIndex; i++ {
		width := (len(v.Output[i]) - v.errIndent)
		if width > maxWidth {
			maxWidth = width
		}
	}
	return maxWidth
}

func (v *Validator) getIndent(indent int) string {
	var strs []string
	for i := 0; i < indent; i++ {
		strs = append(strs, indentor)
	}
	return strings.Join(strs, "")
}

func (v *Validator) FormatVal(val interface{}) string {
	str, ok := val.(string)
	if ok {
		return fmt.Sprintf("\"%s\"", str)
	}
	arr, ok := val.([]interface{})
	if ok {
		vals := make([]string, len(arr))
		for i, sub := range arr {
			vals[i] = v.FormatVal(sub)
		}
		return fmt.Sprintf("[ %s ]", strings.Join(vals, ", "))
	}
	return fmt.Sprintf("%v", val)
}

func (v *Validator) GetIDAndParams(args map[string]interface{}) (string, map[string]interface{}, error) {
	var key string
	var value map[string]interface{}
	found := false
	for k, v := range args {
		val, ok := v.(map[string]interface{})
		if !ok {
			return k, nil, fmt.Errorf("`%v` does not contain any attributes", k)
		}
		key = k
		value = val
		found = true
		break
	}
	if !found {
		return "", nil, fmt.Errorf("no id found")
	}
	return key, value, nil
}

func (v *Validator) bufferKeyValue(key string, val interface{}, indent int) {
	// string
	str, ok := val.(string)
	if ok {
		v.Buffer(fmt.Sprintf("\"%s\": %s", key, str), indent)
		return
	}

	// array
	// TODO: split this into multiline
	arr, ok := val.([]interface{})
	if ok {
		vals := make([]string, len(arr))
		for i, sub := range arr {
			vals[i] = v.FormatVal(sub)
		}
		v.Buffer(fmt.Sprintf("[ %s ]", strings.Join(vals, ", ")), indent)
		return
	}

	// obj
	obj, ok := val.(map[string]interface{})
	if ok {
		v.Buffer(fmt.Sprintf("\"%s\": {", key), indent)
		for subkey, subval := range obj {
			v.bufferKeyValue(subkey, subval, indent+1)
		}
		v.Buffer("}", indent)
	}

	// other
	v.Buffer(fmt.Sprintf("\"%s\": %v", key, val), indent)
}

func (v *Validator) BufferKeyValue(key string, val interface{}, indent int, err error) {
	// if error, start
	if err != nil {
		v.StartError(fmt.Sprintf("%v", err), indent)
	}
	// buffer key / val
	v.bufferKeyValue(key, val, indent)
	// if error, end
	if err != nil {
		v.EndError()
	}
}
