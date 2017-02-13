package tile

const (
	maxMorton = 256 * 256
)

var (
	mx []uint64
	my []uint64
)

// Morton returns the morton code for the provided points. Only works for values
// in the range [0.0: 256.0)]
func Morton(fx float32, fy float32) int {
	x := uint32(fx)
	y := uint32(fy)
	return int(my[y&0xFF] | mx[x&0xFF])
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
