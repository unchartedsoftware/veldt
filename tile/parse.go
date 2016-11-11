package tile

import (
	"encoding/json"
	"fmt"
)

// Parse parses the query payload into the query AST.
func Parse(arg interface{}) (Query, error) {
	// validate the JSON into ot's runtime query components
	validator := NewValidator()
	return validator.Validate(arg)
}
