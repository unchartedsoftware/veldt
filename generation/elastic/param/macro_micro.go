package param

import (
	"fmt"

	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/util/json"
)

const (
	defaultThreshold = 1000
)

// MacroMicro represents params for macro/micro data.
type MacroMicro struct {
	Threshold int64
}

// NewMacroMicro instantiates and returns a new macro/micro parameter object.
func NewMacroMicro(tileReq *tile.Request) (*MacroMicro, error) {
	params := json.GetChildOrEmpty(tileReq.Params, "macro_micro")
	threshold := int64(json.GetNumberDefault(params, defaultThreshold, "threshold"))
	return &MacroMicro{
		Threshold: threshold,
	}, nil
}

// GetHash returns a string hash of the parameter state.
func (p *MacroMicro) GetHash() string {
	return fmt.Sprintf("%d", p.Threshold)
}
