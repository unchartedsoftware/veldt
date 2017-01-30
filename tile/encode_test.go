package tile_test

import (
	"github.com/unchartedsoftware/prism/tile"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Encode", func() {

	a := make([]float32, 2)
	a[0] = 1
	a[1] = 1
	// <[]uint8 | len:8, cap:8>: [0, 0, 128, 63, 0, 0, 128, 63]
	result := make([]uint8, 8, 8)
	result[0] = 0
	result[1] = 0
	result[2] = 128
	result[3] = 63
	result[4] = 0
	result[5] = 0
	result[6] = 128
	result[7] = 63

	It("should set Field and Value", func() {
		b := tile.Encode(a)
		Expect(b).To(Equal(result))
		//	Expect(eq.Field).To(Equal(params["field"]))
		//	Expect(eq.Values).NotTo(BeNil())
	})
})
