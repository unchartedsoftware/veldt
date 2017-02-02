package tile_test

import (
	"github.com/unchartedsoftware/veldt/tile"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lod", func() {

	a := []float32{1.0, 1.0}

	res := []uint8{8, 0, 0, 0, 4, 0, 0, 0, 0, 0, 128, 63, 0, 0, 128, 63, 0, 0, 0, 0}
	res_int := []int{0}
	res_float := []float32{1, 1}

	It("should set Field and Value", func() {
		floats := tile.EncodeLOD(a, int(0))
		Expect(floats).To(Equal(res))
	})

	It("should set LOD", func() {
		floats, ints := tile.LOD(a, int(0))
		Expect(floats).To(Equal(res_float))
		Expect(ints).To(Equal(res_int))
	})

})
