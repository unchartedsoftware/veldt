package tile_test

import (
	"github.com/unchartedsoftware/veldt/tile"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Encode", func() {

	a := []float32{1, 1}

	result := []uint8{0, 0, 128, 63, 0, 0, 128, 63}

	It("should encode", func() {
		b := tile.Encode(a)
		Expect(b).To(Equal(result))
	})
})
