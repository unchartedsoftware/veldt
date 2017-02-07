package tile_test

import (
	"math"

	"github.com/unchartedsoftware/veldt/tile"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/unchartedsoftware/veldt/util/test"
)

var _ = Describe("Macro", func() {

	var macro *tile.Macro
	var input []float32

	BeforeEach(func() {
		macro = &tile.Macro{}
		input = []float32{
			137.24, 7.07, 224.49, 123.95,
			124.51, 148.33, 72.40, 22.15,
			160.13, 77.59, 128.77, 183.32,
			65.36, 36.25, 107.91, 250.01,
			96.05, 198.40, 66.70, 73.39,
		}
	})

	Describe("Parse", func() {
		It("should parse properties from the params argument", func() {
			params := JSON(
				`{
					"lod": 4
				}`)
			err := macro.Parse(params)
			Expect(err).To(BeNil())
			Expect(macro.LOD).To(Equal(4))
		})
	})

	Describe("Encode", func() {
		It("should not encode offset information if LOD == 0", func() {
			bytesPerFloat := 4
			params := JSON(
				`{
					"lod": 0
				}`)
			err := macro.Parse(params)
			Expect(err).To(BeNil())
			bytes, err := macro.Encode(input)
			Expect(err).To(BeNil())
			Expect(len(bytes)).To(Equal(len(input) * bytesPerFloat))
		})
		It("should encode offset information if LOD > 0", func() {
			bytesPerFloat := 4
			bytesPerOffset := 4
			params := JSON(
				`{
					"lod": 4
				}`)
			err := macro.Parse(params)
			Expect(err).To(BeNil())
			bytes, err := macro.Encode(input)
			pointCount := bytesPerFloat
			offsetCount := bytesPerOffset
			pointsBytes := len(input) * bytesPerFloat
			offsetBytes := int(math.Pow(4, 4)) * bytesPerOffset // lod = 4
			totalBytes := pointCount + offsetCount + pointsBytes + offsetBytes
			Expect(err).To(BeNil())
			Expect(len(bytes)).To(Equal(totalBytes))
		})
	})
})
