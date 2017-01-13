package tile

import (
	"encoding/binary"
	"math"
	"sort"
)

const (
	maxMorton = 256 * 256
)

var (
	mx []uint64
	my []uint64
)

// Point represents a tile point with components in the range [0.0: 256.0)
type Point []float32

// Morton returns the morton code for the provided points. Only works for values
// in the range [0.0: 256.0)]
func Morton(fx float32, fy float32) int {
	x := uint32(fx)
	y := uint32(fy)
	return int(my[y&0xFF] | mx[x&0xFF])
}

// LOD takes the input point array and sorts it by morton code. It then
// generates an offset array which match the byte offsets into the point buffer
// for each LOD.
func LOD(data []float32, lod int) ([]float32, []int) {
	// get the points array sorted by morton code
	points := sortPoints(data)

	// generate codes for the sorted points
	codes := make([]int, len(points)/2)
	for i := 0; i < len(points); i += 2 {
		codes[i/2] = Morton(points[i], points[i+1])
	}

	// calc number of partitions and partition stride
	partitions := math.Pow(4, float64(lod))
	paritionStride := maxMorton / int(partitions)

	// set offsets
	offsets := make([]int, int(partitions))
	// init offsets as -1
	for i := range offsets {
		offsets[i] = -1
	}
	// set the offsets to the least byte in the array
	for i := len(codes) - 1; i >= 0; i-- {
		code := codes[i]
		j := code / paritionStride
		offsets[j] = i * 8
	}
	// fill empty offsets up with next entries to ensure easy LOD
	for i := len(offsets) - 1; i >= 0; i-- {
		if offsets[i] == -1 {
			if i == len(offsets)-1 {
				offsets[i] = len(points) * 4
			} else {
				offsets[i] = offsets[i+1]
			}
		}
	}
	return points, offsets
}

// EncodeLOD generates the point LOD offsets and encodes them as a byte array.
func EncodeLOD(data []float32, lod int) []byte {

	// get sorted points and offsets
	points, offsets := LOD(data, lod)

	// encode points
	pointBytes := make([]byte, len(points)*4)
	for i, point := range points {
		binary.LittleEndian.PutUint32(
			pointBytes[i*4:i*4+4],
			math.Float32bits(point))
	}

	// encode offsets
	offsetBytes := make([]byte, len(offsets)*4)
	for i, offset := range offsets {
		binary.LittleEndian.PutUint32(
			offsetBytes[i*4:i*4+4],
			uint32(offset))
	}

	// point length
	pointLength := make([]byte, 4)
	binary.LittleEndian.PutUint32(
		pointLength,
		uint32(len(pointBytes)))

	// offset length
	offsetLength := make([]byte, 4)
	binary.LittleEndian.PutUint32(
		offsetLength,
		uint32(len(offsetBytes)))

	a := len(pointLength)
	b := len(offsetLength)
	c := len(pointBytes)
	d := len(offsetBytes)

	bytes := make([]byte, a+b+c+d)

	copy(bytes[0:a], pointLength)
	copy(bytes[a:a+b], offsetLength)
	copy(bytes[a+b:a+b+c], pointBytes)
	copy(bytes[a+b+c:a+b+c+d], offsetBytes)

	return bytes
}

func init() {
	// init the morton code lookups.
	mx = []uint64{0, 1}
	my = []uint64{0, 2}
	for i := 4; i < 0xFFFF; i <<= 2 {
		l := len(mx)
		for j := 0; j < l; j++ {
			mx = append(mx, mx[j]|uint64(i))
			my = append(my, (mx[j]|uint64(i))<<1)
		}
	}
}

func sortPoints(data []float32) []float32 {
	points := make(pointArray, len(data)/2)
	for i := 0; i < len(data); i += 2 {
		x := data[i]
		y := data[i+1]
		points[i/2] = [2]float32{x, y}
	}
	// sort the points
	sort.Sort(points)
	// convert to flat array
	res := make([]float32, len(points)*2)
	for i, point := range points {
		res[i*2] = point[0]
		res[i*2+1] = point[1]
	}
	return res
}

type pointArray [][2]float32

func (p pointArray) Len() int {
	return len(p)
}
func (p pointArray) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p pointArray) Less(i, j int) bool {
	return Morton(p[i][0], p[i][1]) < Morton(p[j][0], p[j][1])
}
