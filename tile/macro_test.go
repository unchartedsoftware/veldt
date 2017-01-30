package tile_test

import (
	"github.com/unchartedsoftware/prism/tile"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Macro", func() {
	eq := &tile.Macro{}
	eq2 := &tile.Macro{}

	// create params
	// we use the built in `make` function to allocate the map
	params := make(map[string]interface{})
	params["lod"] = 1

	params_2 := make([]float32, 2)
	params_2[0] = float32(1.0)
	params_2[1] = float32(1.0)

	result := make([]uint8, 8, 8)
	result[0] = 0
	result[1] = 0
	result[2] = 128
	result[3] = 63
	result[4] = 0
	result[5] = 0
	result[6] = 128
	result[7] = 63

	It("should set LOD field", func() {
		eq.Parse(params)
		Expect(eq.LOD).To(Equal(0))
	})

	It("should be nil on wrong input", func() {
		rslt, _ := eq2.Encode(params_2)
		Expect(rslt).To(Equal(result))
	})
})
