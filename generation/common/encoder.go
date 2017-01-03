package common

import (
	"encoding/json"
	"math"
	"sort"

	"github.com/unchartedsoftware/prism/tile"
)

//Micro tile encoding is shared between the generators.
func EncodeMicroTileResult(hits []map[string]interface{}, points []float32, lod int) ([]byte, error) {

	if lod > 0 {
		// sort hits by morton code so they align
		sortHitsArray(hits, points)
		// sort points and get offsets
		sortedPoints, offsets := tile.LOD(points, lod)
		return json.Marshal(map[string]interface{}{
			"points":  sortedPoints,
			"offsets": offsets,
			"hits":    hits,
		})
	}

	return json.Marshal(map[string]interface{}{
		"points": points,
		"hits":   hits,
	})
}

func CastPixelResult(value interface{}) float64 {
	val64, ok := value.(float64)
	if ok {
		return val64
	}

	valint64, ok := value.(int64)
	if ok {
		val64 = float64(valint64)
	}
	return val64
}

func ExistsIn(val string, arr []string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func ToFixed(num float32, precision int) float32 {
	output := math.Pow(10, float64(precision))
	return float32(math.Floor(float64(num)*output+0.5)) / float32(output)
}

func sortHitsArray(hits []map[string]interface{}, points []float32) {
	if hits == nil {
		return
	}
	// sort hits by morton code so they align
	hitsArr := make(hitsArray, len(hits))
	for i, hit := range hits {
		// add to hits array
		hitsArr[i] = &hitWrapper{
			x:    points[i*2],
			y:    points[i*2+1],
			data: hit,
		}
	}
	sort.Sort(hitsArr)
	// copy back into same arr
	for i, hit := range hitsArr {
		hits[i] = hit.data
	}
}

type hitWrapper struct {
	x    float32
	y    float32
	data map[string]interface{}
}

type hitsArray []*hitWrapper

func (h hitsArray) Len() int {
	return len(h)
}
func (h hitsArray) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
func (h hitsArray) Less(i, j int) bool {
	return tile.Morton(h[i].x, h[i].y) < tile.Morton(h[j].x, h[j].y)
}
