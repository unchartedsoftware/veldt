package generation

import (
	"encoding/binary"
	"math"

	"github.com/unchartedsoftware/prism/param"
	"github.com/unchartedsoftware/prism/binning"
)

// Heatmap represents a bivariate heatmap tile generator.
type Heatmap struct {
	Binning *param.Bivariate
}

// SetHeatmapParams sets the params for the specific generator.
func SetHeatmapParams(arg interface{}, coord *binning.TileCoord, params map[string]interface{}) error {
	return SetBivariateParams(arg, coord, params)
}

// Float64ToBytes converts a []float64 to a []uint8 of equal byte size.
func (t *Heatmap) Float64ToBytes(arr []float64) []byte {
	bits := make([]byte, len(arr)*8)
	for i, val := range arr {
		binary.LittleEndian.PutUint64(
			bits[i*8:i*8+8],
			math.Float64bits(val))
	}
	return bits[0:]
}
