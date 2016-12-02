package tile

import (
	"encoding/binary"
	"math"
)

func Encode(points []float32) []byte {
	bytes := make([]byte, len(points)*4)
	for i := 0; i < len(points); i += 2 {
		x := points[i]
		y := points[i+1]
		binary.LittleEndian.PutUint32(
			bytes[i*4:i*4+4],
			math.Float32bits(x))
		binary.LittleEndian.PutUint32(
			bytes[i*4+4:i*4+8],
			math.Float32bits(y))
	}
	return bytes
}
