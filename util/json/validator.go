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
	output          []string
	errLines        map[int]bool
	nextIndentation int
	indentation     []int
	errStartIndex   int
	errEndIndex     int
	errIndent       int
	errHeaderIndex  int
	errFooterIndex  int
	errMsg          string
	err             bool
}

// StartObject begins the buffering of an object.
func (v *Validator) StartObject() {
	v.buffer("{")
	v.nextIndentation++
}

// StartSubObject begins the buffering of a nested object to a key.
func (v *Validator) StartSubObject(key string) {
	v.buffer(fmt.Sprintf(`"%s": {`, key))
	v.nextIndentation++
}

// EndObject ends the buffering of the current object.
func (v *Validator) EndObject() {
	if v.nextIndentation > 0 {
		v.nextIndentation--
	}
	v.buffer("}")
}

// StartArray begins the buffering of an array.
func (v *Validator) StartArray() {
	v.buffer("[")
	v.nextIndentation++
}

// StartSubArray begins the buffering of an array value to a key.
func (v *Validator) StartSubArray(key string) {
	v.buffer(fmt.Sprintf(`"%s": [`, key))
	v.nextIndentation++
}

// EndArray ends the buffering of the current array.
func (v *Validator) EndArray() {
	if v.nextIndentation > 0 {
		v.nextIndentation--
	}
	v.buffer("]")
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
	length := v.Size()
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
		if v.indentation[i] != v.indentation[j] ||
			v.output[i] == "{" ||
			v.output[i] == "[" {
			formatted[i] = v.output[i]
		} else {
			formatted[i] = fmt.Sprintf("%s,", v.output[i])
		}
	}
	// return the concatenated output
	return strings.Join(formatted, "\n")
}

// StartError begins wrapping an error portion of the output buffer.
func (v *Validator) StartError(msg string) {
	v.err = true
	v.errHeaderIndex = v.Size()
	v.errStartIndex = v.Size() + 1
	v.errIndent = v.nextIndentation
	v.errMsg = msg
	v.buffer("") // header line
}

// EndError ends wrapping an error portion of the output buffer.
func (v *Validator) EndError() {
	v.errEndIndex = v.Size()
	v.buffer("") // footer line
	width := v.getErrWidth()
	header := v.getErrHeader(width)
	footer := v.getErrFooter(width)
	v.output[v.errHeaderIndex] = header
	v.output[v.errEndIndex] = footer
	// track which lines have errors
	if v.errLines == nil {
		v.errLines = make(map[int]bool)
	}
	for i := v.errHeaderIndex; i <= v.errEndIndex; i++ {
		v.errLines[i] = true
	}
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
			v.getIndentString(v.errIndent),
			v.getErrAnnotations(width, "v"),
			color.Reset)
	}
	return fmt.Sprintf("%s%s",
		v.getIndentString(v.errIndent),
		v.getErrAnnotations(width, "v"))
}

func (v *Validator) getErrFooter(width int) string {
	if color.ColorTerminal {
		return fmt.Sprintf("%s%s%s Error: %s%s",
			color.Red,
			v.getIndentString(v.errIndent),
			v.getErrAnnotations(width, "^"),
			v.errMsg,
			color.Reset)
	}
	return fmt.Sprintf("%s%s Error: %s",
		v.getIndentString(v.errIndent),
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

func (v *Validator) getIndentString(indent int) string {
	var strs []string
	for i := 0; i < indent; i++ {
		strs = append(strs, indentor)
	}
	return strings.Join(strs, "")
}

func (v *Validator) bufferKeyValue(key string, val interface{}) {
	// string
	str, ok := val.(string)
	if ok {
		v.buffer(fmt.Sprintf(`"%s": "%s"`, key, str))
		return
	}

	// array
	arr, ok := val.([]interface{})
	if ok {
		v.StartSubArray(key)
		for _, sub := range arr {
			v.bufferValue(sub)
		}
		v.EndArray()
		return
	}

	// obj
	obj, ok := val.(map[string]interface{})
	if ok {
		v.StartSubObject(key)
		for subkey, subval := range obj {
			v.bufferKeyValue(subkey, subval)
		}
		v.EndObject()
		return
	}

	// other
	v.buffer(fmt.Sprintf(`"%s": %v`, key, val))
}

// BufferKeyValue will buffer the a JSON key and it's value with correct
// indentation.
func (v *Validator) BufferKeyValue(key string, val interface{}, err error) {
	// if error, start
	if err != nil {
		v.StartError(fmt.Sprintf("%v", err))
	}
	// buffer key / val
	v.bufferKeyValue(key, val)
	// if error, end
	if err != nil {
		v.EndError()
	}
}

func (v *Validator) bufferValue(val interface{}) {
	// string
	str, ok := val.(string)
	if ok {
		v.buffer(fmt.Sprintf(`"%s"`, str))
		return
	}

	// array
	arr, ok := val.([]interface{})
	if ok {
		v.StartArray()
		for _, sub := range arr {
			v.bufferValue(sub)
		}
		v.EndArray()
		return
	}

	// obj
	obj, ok := val.(map[string]interface{})
	if ok {
		v.StartObject()
		for subkey, subval := range obj {
			v.bufferKeyValue(subkey, subval)
		}
		v.EndObject()
		return
	}

	// other
	v.buffer(fmt.Sprintf("%v", val))
}

// BufferValue will buffer the a JSON value with correct indentation.
func (v *Validator) BufferValue(val interface{}, err error) {
	// if error, start
	if err != nil {
		v.StartError(fmt.Sprintf("%v", err))
	}
	// buffer val
	v.bufferValue(val)
	// if error, end
	if err != nil {
		v.EndError()
	}
}

func (v *Validator) buffer(str string) {
	line := fmt.Sprintf("%s%s", v.getIndentString(v.nextIndentation), str)
	v.output = append(v.output, line)
	v.indentation = append(v.indentation, v.nextIndentation)
}
