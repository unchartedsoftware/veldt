package meta

import (
	"fmt"
)

// Request represents a meta data request.
type Request struct {
	Type  string `json:"type"`
	URI   string `json:"uri"`
	Store string `json:"store"`
}

// String returns the request formatted as a string.
func (r *Request) String() string {
	return fmt.Sprintf("%s/%s",
		r.Type,
		r.URI)
}
