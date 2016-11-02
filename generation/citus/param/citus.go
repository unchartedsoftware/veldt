package param

import (
	"fmt"
	"sort"
	"strings"

	"github.com/unchartedsoftware/prism/tile"
	"github.com/unchartedsoftware/prism/util/json"
)

// Citus represents params specifically tailored to citus queries.
type Citus struct {
	Types []string
}

// NewCitus instantiates and returns a new macro/micro parameter object.
func NewCitus(tileReq *tile.Request) (*Citus, error) {
	params := json.GetChildOrEmpty(tileReq.Params, "citus")
	types, _ := json.GetStringArray(params, "types")
	sort.Strings(types)
	return &Citus{
		Types: types,
	}, nil
}

// GetHash returns a string hash of the parameter state.
func (p *Citus) GetHash() string {
	return fmt.Sprintf("%s",
		strings.Join(p.Types, ":"))
}
