package tile_test

import (
	"github.com/unchartedsoftware/veldt/tile"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/unchartedsoftware/veldt/util/test"
)

var _ = Describe("Macro", func() {

	var macro *tile.Macro

	BeforeEach(func() {
		macro = &tile.Macro{}
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
})
