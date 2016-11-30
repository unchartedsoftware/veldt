package tile

import (
	"encoding/binary"
	"math"
)

func Encode(points []float32) []byte {
	bytes := make([]byte, len(points)*8)
	for i := 0; i < len(points); i += 2 {
		x := points[i]
		y := points[i+1]
		binary.LittleEndian.PutUint32(
			bytes[i*8:i*8+4],
			math.Float32bits(x))
		binary.LittleEndian.PutUint32(
			bytes[i*8+4:i*8+8],
			math.Float32bits(y))
	}
	return bytes
}
