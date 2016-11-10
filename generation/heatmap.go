package generation

import (
	"encoding/binary"
	"math"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/param"
)

// Heatmap represents a bivariate heatmap tile generator.
type Heatmap struct {
	Bivariate
}

// Parse parses the provided JSON object and populates the tiles attributes.
func (h *Heatmap) Parse(coord *binning.TileCoord, params map[string]interface{}) error {
	return h.Bivariate.Parse(coord, params)
}

// Float64ToBytes converts a []float64 to a []uint8 of equal byte size.
func (h *Heatmap) Float64ToBytes(arr []float64) []byte {
	bits := make([]byte, len(arr)*8)
	for i, val := range arr {
		binary.LittleEndian.PutUint64(
			bits[i*8:i*8+8],
			math.Float64bits(val))
	}
	return bits[0:]
}
