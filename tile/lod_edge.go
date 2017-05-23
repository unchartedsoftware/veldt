package tile

import (
	"math"
	"sort"

	"github.com/unchartedsoftware/veldt/binning"
)

const (
	edgeStride = 6
)

// EdgeLOD takes the input edge array and sorts it by the morton code of the
// first point in the edge. It then generates an offset array which match the
// byte offsets into the edge buffer for each LOD. This is used at runtime to
// only render quadrants of the generated tile.
func EdgeLOD(data []float32, lod int) ([]float32, []int) {
	// get the edges array sorted by morton code
	edges := sortEdges(data)

	// generate codes for the sorted edges
	codes := make([]int, len(edges)/edgeStride)
	for i := 0; i < len(edges); i += edgeStride {
		sx := edges[i]   // src x
		sy := edges[i+1] // src y
		// sort based on src point
		codes[i/edgeStride] = Morton(sx, sy)
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
		offsets[j] = i * (bytesPerComponent * edgeStride)
	}
	// fill empty offsets up with next entries to ensure easy LOD
	for i := len(offsets) - 1; i >= 0; i-- {
		if offsets[i] == -1 {
			if i == len(offsets)-1 {
				offsets[i] = len(edges) * bytesPerComponent
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
	edges := make(edgeArray, len(data)/edgeStride)
	maxPixel := float32(binning.MaxTileResolution)
	for i := 0; i < len(data); i += edgeStride {
		ax := data[i]   // src x
		ay := data[i+1] // src y
		aw := data[i+2] // src weight
		bx := data[i+3] // dst x
		by := data[i+4] // dst y
		bw := data[i+5] // dst weight
		// ensure first point is within the tile
		if ax >= 0.0 && ax < maxPixel &&
			ay >= 0.0 && ay < maxPixel {
			edges[i/edgeStride] = [edgeStride]float32{ax, ay, aw, bx, by, bw}
		} else {
			edges[i/edgeStride] = [edgeStride]float32{bx, by, bw, ax, ay, aw}
		}
	}
	// sort the edges
	sort.Sort(edges)
	// convert to flat array
	res := make([]float32, len(edges)*edgeStride)
	for i, edge := range edges {
		res[i*edgeStride] = edge[0]   // src x
		res[i*edgeStride+1] = edge[1] // src y
		res[i*edgeStride+2] = edge[2] // src weight
		res[i*edgeStride+3] = edge[3] // dst x
		res[i*edgeStride+4] = edge[4] // dst y
		res[i*edgeStride+5] = edge[5] // dst weight
	}
	return res
}

type edgeArray [][edgeStride]float32 // srcX, srcY, srcWeight, dstX, dstY, dstWeight

func (e edgeArray) Len() int {
	return len(e)
}
func (e edgeArray) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
func (e edgeArray) Less(i, j int) bool {
	return Morton(e[i][0], e[i][1]) < Morton(e[j][0], e[j][1])
}
