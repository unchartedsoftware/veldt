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
func Morton(fx float32, fy float32) uint64 {
	x := uint32(fx)
	y := uint32(fy)
	return (my[y&0xFF] | mx[x&0xFF]) // + (my[(y>>8)&0xFF]|mx[(x>>8)&0xFF])*0x10000
}

// LOD takes the input point array and sorts it by morton code. It then
// generates an offset array which match the byte offsets into the point buffer
// for each LOD.
func LOD(data []float32, lod int) ([]float32, []int) {
	arr := newPointArray(data)

	partitions := math.Pow(4, float64(lod))
	paritionStride := maxMorton / int(partitions)

	// set offsets
	offsets := make([]int, int(partitions))
	// init offsets as -1
	for i := range offsets {
		offsets[i] = -1
	}
	// set the offsets to the least byte in the array
	for i := len(arr.codes) - 1; i >= 0; i-- {
		code := arr.codes[i]
		j := code / paritionStride
		offsets[j] = i * 8
	}
	// fill empty offsets up with next entries to ensure easy LOD
	for i := len(offsets) - 1; i >= 0; i-- {
		if offsets[i] == -1 {
			if i == len(offsets)-1 {
				offsets[i] = len(arr.points) * 8
			} else {
				offsets[i] = offsets[i+1]
			}
		}
	}
	// convert to point array
	points := make([]float32, len(arr.points)*2)
	for i, point := range arr.points {
		points[i*2] = point[0]
		points[i*2+1] = point[1]
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
	// TODO: we only use 0 -> 65536, reduce the length of this array
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

type pointArray struct {
	points [][]float32
	codes  []int
}

func newPointArray(data []float32) *pointArray {
	points := make([][]float32, len(data)/2)
	for i := 0; i < len(data); i += 2 {
		x := data[i]
		y := data[i+1]
		points[i/2] = []float32{x, y}
	}
	arr := &pointArray{
		points: points,
	}
	// sort the points
	sort.Sort(arr)
	// now generate codes for the sorted points
	codes := make([]int, len(points))
	for i, p := range points {
		codes[i] = int(Morton(p[0], p[1]))
	}
	arr.codes = codes
	return arr
}

func (p pointArray) Len() int {
	return len(p.points)
}
func (p pointArray) Swap(i, j int) {
	p.points[i], p.points[j] = p.points[j], p.points[i]
}
func (p pointArray) Less(i, j int) bool {
	return Morton(p.points[i][0], p.points[i][1]) < Morton(p.points[j][0], p.points[j][1])
}
