package json

import (
	"fmt"
	"strings"

	"github.com/unchartedsoftware/veldt/util/color"
)

const (
	indentor = "    "
)

// Validator parses a JSON query expression into its typed format. It
// ensure all types are correct and that the syntax is valid.
type Validator struct {
	output         []string
	errLines       map[int]bool
	indentation    []int
	errStartIndex  int
	errEndIndex    int
	errIndent      int
	errHeaderIndex int
	errFooterIndex int
	errMsg         string
	err            bool
}

// Buffer adds a string at the appropriate indent to the output buffer.
func (v *Validator) Buffer(str string, indent int) {
	line := fmt.Sprintf("%s%s", v.getIndent(indent), str)
	v.output = append(v.output, line)
	v.indentation = append(v.indentation, indent)
}

// Size returns the length of the current output buffer.
func (v *Validator) Size() int {
	return len(v.output)
}

// HasError returns true if an error has been encountered.
func (v *Validator) HasError() bool {
	return v.err
}

// Error returns the error if there is one.
func (v *Validator) Error() error {
	if v.err {
		return fmt.Errorf(v.String())
	}
	return nil
}

// String returns the string in the output buffer.
func (v *Validator) String() string {
	length := len(v.output)
	formatted := make([]string, length)
	// determine whether or not to append a comma on the end based on the next
	// lines indentation
	for i := 0; i < length; i++ {
		if i == length-1 {
			// last line
			formatted[i] = v.output[i]
			break
		}
		// skip any error annotation lines
		if v.errLines[i] {
			formatted[i] = v.output[i]
			continue
		}
		// get the next line that isn't an error
		j := i + 1
		for {
			// until the next line that isn't an error
			if v.errLines[j] {
				j++
				continue
			}
			break
		}
		if j > length-1 {
			// no more lines, this means the output is malformed
			break
		}
		if v.indentation[i] != v.indentation[j] {
			formatted[i] = v.output[i]
		} else {
			formatted[i] = v.output[i] + ","
		}
	}
	// return the concatenated output
	return strings.Join(formatted, "\n")
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
	v.output[v.errHeaderIndex] = header
	v.output[v.errEndIndex] = footer
	// track which lines have errors
	if v.errLines == nil {
		v.errLines = make(map[int]bool)
	}
	v.errLines[v.errHeaderIndex] = true
	v.errLines[v.errEndIndex] = true
}

func (v *Validator) getErrAnnotations(width int, char string) string {
	arr := make([]string, width)
	for i := 0; i < width; i++ {
		arr[i] = char
	}
	return strings.Join(arr, "")
}

func (v *Validator) getErrHeader(width int) string {
	if color.ColorTerminal {
		return fmt.Sprintf("%s%s%s%s",
			color.Red,
			v.getIndent(v.errIndent),
			v.getErrAnnotations(width, "v"),
			color.Reset)
	}
	return fmt.Sprintf("%s%s",
		v.getIndent(v.errIndent),
		v.getErrAnnotations(width, "v"))
}

func (v *Validator) getErrFooter(width int) string {
	if color.ColorTerminal {
		return fmt.Sprintf("%s%s%s Error: %s%s",
			color.Red,
			v.getIndent(v.errIndent),
			v.getErrAnnotations(width, "^"),
			v.errMsg,
			color.Reset)
	}
	return fmt.Sprintf("%s%s Error: %s",
		v.getIndent(v.errIndent),
		v.getErrAnnotations(width, "^"),
		v.errMsg)
}

func (v *Validator) getErrWidth() int {
	maxWidth := 1
	indentLength := v.errIndent * len(indentor)
	for i := v.errStartIndex; i < v.errEndIndex; i++ {
		width := (len(v.output[i]) - indentLength)
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

func (v *Validator) formatVal(val interface{}) string {
	str, ok := val.(string)
	if ok {
		return fmt.Sprintf(`"%s"`, str)
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

// GetIDAndParams returns the nested key and value from a JSON object of the
// form:
//     {
//         "key": {
//             "prop0": ...,
//             "prop1": ...,
//             "prop2": ...,
//         }
//     }
//
func (v *Validator) GetIDAndParams(args map[string]interface{}) (string, interface{}, error) {
	for k, v := range args {
		return k, v, nil
	}
	return "", nil, fmt.Errorf("no id found")
}

func (v *Validator) bufferKeyValue(key string, val interface{}, indent int) {
	// string
	str, ok := val.(string)
	if ok {
		v.Buffer(fmt.Sprintf(`"%s": "%s"`, key, str), indent)
		return
	}

	// array
	// TODO: split this into multiline
	arr, ok := val.([]interface{})
	if ok {
		vals := make([]string, len(arr))
		for i, sub := range arr {
			vals[i] = v.formatVal(sub)
		}
		v.Buffer(fmt.Sprintf(`"%s": [ %s ]`, key, strings.Join(vals, ", ")), indent)
		return
	}

	// obj
	obj, ok := val.(map[string]interface{})
	if ok {
		v.Buffer(fmt.Sprintf(`"%s": {`, key), indent)
		for subkey, subval := range obj {
			v.bufferKeyValue(subkey, subval, indent+1)
		}
		v.Buffer("}", indent)
		return
	}

	// other
	v.Buffer(fmt.Sprintf(`"%s": %v`, key, val), indent)
}

// BufferKeyValue will buffer the a JSON key and it's value with correct
// indentation.
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
