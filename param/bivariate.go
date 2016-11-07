package param

import (
	"fmt"
	"math"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/util/json"
)

// Bivariate represents a bivariate tile generator.
type Bivariate struct {
	XField     string
	YField     string
	Bounds     *binning.Bounds
	MinX       float64
	MaxX       float64
	MinY       float64
	MaxY       float64
	XRange     float64
	YRange     float64
	Resolution int
	BinSizeX   float64
	BinSizeY   float64
}
