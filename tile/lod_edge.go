package tile

import (
	"math"
	"sort"
)

// EdgeLOD takes the input edge array and sorts it by the morton code of the
// first point in the edge. It then generates an offset array which match the
// byte offsets into the edge buffer for each LOD. This is used at runtime to
// only render quadrants of the generated tile.
func EdgeLOD(data []float32, lod int) ([]float32, []int) {
	// get the edges array sorted by morton code
	edges := sortEdges(data)

	// generate codes for the sorted edges
	codes := make([]int, len(edges)/4)
	for i := 0; i < len(edges); i += 4 {
		codes[i/4] = Morton(edges[i], edges[i+1])
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
		offsets[j] = i * (4 * 4)
	}
	// fill empty offsets up with next entries to ensure easy LOD
	for i := len(offsets) - 1; i >= 0; i-- {
		if offsets[i] == -1 {
			if i == len(offsets)-1 {
				offsets[i] = len(edges) * 4
			} else {
				offsets[i] = offsets[i+1]
			}
		}
	}
	return edges, offsets
}

// EncodeEdgeLOD generates the point LOD offsets and encodes them as a byte array.
func EncodeEdgeLOD(data []float32, lod int) []byte {
	// get sorted points and offsets
	edges, offsets := EdgeLOD(data, lod)
	// encode the results
	return encodeLOD(edges, offsets)
}

func sortEdges(data []float32) []float32 {
	edges := make(edgeArray, len(data)/4)
	for i := 0; i < len(data); i += 4 {
		ax := data[i]
		ay := data[i+1]
		bx := data[i+2]
		by := data[i+3]
		edges[i/4] = [4]float32{ax, ay, bx, by}
	}
	// sort the edges
	sort.Sort(edges)
	// convert to flat array
	res := make([]float32, len(edges)*4)
	for i, edge := range edges {
		res[i*4] = edge[0]
		res[i*4+1] = edge[1]
		res[i*4+2] = edge[2]
		res[i*4+3] = edge[3]
	}
	return res
}

type edgeArray [][4]float32

func (e edgeArray) Len() int {
	return len(e)
}
func (e edgeArray) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
func (e edgeArray) Less(i, j int) bool {
	return Morton(e[i][0], e[i][1]) < Morton(e[j][0], e[j][1])
}
