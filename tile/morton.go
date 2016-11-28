package tile

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

func Morton(x uint32, y uint32) uint64 {
	return (my[y&0xFF] | mx[x&0xFF]) +
		(my[(y>>8)&0xFF]|mx[(x>>8)&0xFF])*0x10000 +
		(my[(y>>16)&0xFF]|mx[(x>>16)&0xFF])*0x100000000
}
