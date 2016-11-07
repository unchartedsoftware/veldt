package generation

import (
	"encoding/binary"
	"math"

	"github.com/unchartedsoftware/prism/binning"
)

// Density represents a univariate density strip generator.
type Density struct {
	Binning *param.Univariate
}

// SetDensityParams sets the tiling params on a tiling generator.
func SetDensityParams(arg interface{}, coord *binning.TileCoord, params map[string]interface{}) error {
	return SetUnivariateParams(arg, coord, params)
}

// Float64ToBytes converts a []float64 to a []uint8 of equal byte size.
func (t *Density) Float64ToBytes(arr []float64) []byte {
	bits := make([]byte, len(arr)*8)
	for i, val := range arr {
		binary.LittleEndian.PutUint64(
			bits[i*8:i*8+8],
			math.Float64bits(val))
	}
	return bits[0:]
}
