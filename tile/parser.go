package tile

import (
	"encoding/json"
	"fmt"
)

// Parse parses the query payload into the query AST.
func Parse(bytes []byte) (Query, error) {
	// unmarshal the query
	var token interface{}
	err := json.Unmarshal(bytes, &token)
	if err != nil {
		return nil, fmt.Errorf("Invalid JSON: %v", err)
	}
	// validate the JSON into ot's runtime query components
	validator := NewValidator()
	exp, err := validator.Validate(token)
	if err != nil {
		return nil, err
	}
	// parse into correct AST
	return parseToken(exp)
}
