package param

import (
	"fmt"
)

var (
	// ErrMissing is the error returned if the param node is missing.
	ErrMissing = fmt.Errorf("Missing query")
)

// IsOptionalErr returns true if the error is not nil and is also not of type
// `ErrMissing`.
func IsOptionalErr(err error) bool {
	return err != nil && err != ErrMissing
}
