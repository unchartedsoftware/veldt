package tile_test

import (
	"github.com/unchartedsoftware/veldt/tile"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Encode", func() {

	var input []float32
	var bytes []byte

	BeforeEach(func() {
		input = []float32{
			137.24, 7.07, 224.49, 123.95,
			124.51, 148.33, 72.40, 22.15,
			160.13, 77.59, 128.77, 183.32,
			65.36, 36.25, 107.91, 250.01,
			96.05, 198.40, 66.70, 73.39,
		}
		bytes = []byte{
			113, 61, 9, 67, 113, 61, 226, 64, 113, 125, 96, 67, 102,
			230, 247, 66, 31, 5, 249, 66, 123, 84, 20, 67, 205, 204,
			144, 66, 51, 51, 177, 65, 72, 33, 32, 67, 20, 46, 155, 66,
			31, 197, 0, 67, 236, 81, 55, 67, 82, 184, 130, 66, 0, 0, 17,
			66, 236, 209, 215, 66, 143, 2, 122, 67, 154, 25, 192, 66,
			102, 102, 70, 67, 102, 102, 133, 66, 174, 199, 146, 66,
		}
	})

	Describe("Encode", func() {
		It("should encode provided []float32 into a []byte", func() {
			bs := tile.Encode(input)
			Expect(bs).To(Equal(bytes))
		})
	})
})
