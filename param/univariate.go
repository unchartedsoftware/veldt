package param

import (
	"fmt"
	"math"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/util/json"
)

// Univariate represents a univariate tile generator.
type Univariate struct {
	Field      string
	Min        float64
	Max        float64
	Range      float64
	Resolution int
	BinSize    float64
}
