package tile

import (
	"encoding/binary"
	"math"
	"sort"
)

var (
	mx []uint64
	my []uint64
)

func init() {
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

func morton(fx float32, fy float32) uint64 {
	x := uint32(fx)
	y := uint32(fy)
	return (my[y&0xFF] | mx[x&0xFF]) + (my[(y>>8)&0xFF]|mx[(x>>8)&0xFF])*0x10000
}

type point []float32

type pointArray struct {
	points []point
	codes  []int
}

func newPointArray(data []float32) *pointArray {
	points := make([]point, len(data)/2)
	for i := 0; i < len(data); i += 2 {
		x := data[i]
		y := data[i+1]
		points[i/2] = point{x, y}
	}
	arr := &pointArray{
		points: points,
	}
	// sort the points
	sort.Sort(arr)
	// now generate codes for the sorted points
	codes := make([]int, len(points))
	for i, p := range points {
		codes[i] = int(morton(p[0], p[1]))
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
	return morton(p.points[i][0], p.points[i][1]) < morton(p.points[j][0], p.points[j][1])
}

func EncodeLOD(data []float32, lod int) []byte {

	arr := newPointArray(data)

	max := 65536
	partitions := math.Pow(4, float64(lod))
	paritionStride := max / int(partitions)

	// set offsets
	offsets := make([]int, int(partitions))
	for i := range offsets {
		offsets[i] = -1
	}

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

	// encode points
	pointBytes := make([]byte, len(arr.points)*8)
	for i, p := range arr.points {
		binary.LittleEndian.PutUint32(
			pointBytes[i*8:i*8+4],
			math.Float32bits(p[0]))
		binary.LittleEndian.PutUint32(
			pointBytes[i*8+4:i*8+8],
			math.Float32bits(p[1]))
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
