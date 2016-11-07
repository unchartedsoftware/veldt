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
