package generation

import (
	"encoding/binary"
	"math"

	"github.com/unchartedsoftware/prism/binning"
)

// Density represents a univariate density strip generator.
type Density struct {
	Univariate
}

// Parse parses the provided JSON object and populates the tiles attributes.
func (d *Density) Parse(coord *binning.TileCoord, params map[string]interface{}) error {
	return h.Univariate.Parse(coord, params)
}

// Float64ToBytes converts a []float64 to a []uint8 of equal byte size.
func (d *Density) Float64ToBytes(arr []float64) []byte {
	bits := make([]byte, len(arr)*8)
	for i, val := range arr {
		binary.LittleEndian.PutUint64(
			bits[i*8:i*8+8],
			math.Float64bits(val))
	}
	return bits[0:]
}
