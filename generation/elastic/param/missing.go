package param

import (
	"fmt"
)

const (
	// MissingPrefix represents the string prefix of a missing parameter error.
	MissingPrefix = "Missing"
)

// IsOptionalErr returns true if the error is not nil and also does not contain
// the prefix of a missing property error.
func IsOptionalErr(err error) bool {
	return err != nil && fmt.Sprintf("%v", err)[0:len(MissingPrefix)] != MissingPrefix
}
