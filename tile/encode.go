package tile

import (
	"encoding/binary"
	"math"
)

// EncodeFloat32 takes a []float32 values and returns an encoded byte array in
// little endian format.
func EncodeFloat32(data []float32) []byte {
	bytes := make([]byte, len(data)*4)
	for i, datum := range data {
		fixed := toFixed(datum, 2) // 2 decimal places
		binary.LittleEndian.PutUint32(
			bytes[i*4:i*4+4],
			math.Float32bits(fixed))
	}
	return bytes
}

// EncodeInt takes a []int values and returns an encoded byte array in little
// endian format.
func EncodeInt(data []int) []byte {
	bytes := make([]byte, len(data)*4)
	for i, datum := range data {
		binary.LittleEndian.PutUint32(
			bytes[i*4:i*4+4],
			uint32(datum))
	}
	return bytes
}

// encodeLOD generates the point LOD offsets and encodes them as a byte array.
func encodeLOD(data []float32, offsets []int) []byte {
	// encode data
	dataBytes := EncodeFloat32(data)
	// encode offsets
	offsetBytes := EncodeInt(offsets)
	// data length
	dataLength := make([]byte, 4)
	binary.LittleEndian.PutUint32(
		dataLength,
		uint32(len(dataBytes)))
	// offset length
	offsetLength := make([]byte, 4)
	binary.LittleEndian.PutUint32(
		offsetLength,
		uint32(len(offsetBytes)))
	// combine the buffers
	a := len(dataLength)
	b := len(offsetLength)
	c := len(dataBytes)
	d := len(offsetBytes)
	// create buffer
	bytes := make([]byte, a+b+c+d)
	// copy into buffer
	copy(bytes[0:a], dataLength)
	copy(bytes[a:a+b], offsetLength)
	copy(bytes[a+b:a+b+c], dataBytes)
	copy(bytes[a+b+c:a+b+c+d], offsetBytes)
	// return buffer
	return bytes
}

func toFixed(num float32, precision int) float32 {
	output := math.Pow(10, float64(precision))
	return float32(math.Floor(float64(num)*output+0.5)) / float32(output)
}
