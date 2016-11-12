package tile

import (
	"github.com/unchartedsoftware/prism"
)

// Parse parses the query payload into the query AST.
func Parse(arg interface{}) (prism.Query, error) {
	// validate the JSON into ot's runtime query components
	validator := NewValidator()
	return validator.Validate(arg)
}
