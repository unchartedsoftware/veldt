package meta

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

func init() {
	spew.Config.SortKeys = true
}

// Request represents a meta data request.
type Request struct {
	Type  string
	Param  param.Params
	URI   string
	Store string
}

// GetHash returns a unique hash for the request.
func (r *Request) GetHash() string {
	return fmt.Sprintf("%s:%s", "meta", spew.Dump(r))
}
