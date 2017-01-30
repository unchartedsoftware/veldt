package tile_test

import (
	"github.com/unchartedsoftware/veldt/tile"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Macro", func() {
	eq := &tile.Macro{}
	eq2 := &tile.Macro{}

	params := make(map[string]interface{})
	params["lod"] = 1

	params_2 := []float32{1.0, 1.0}

	result := []uint8{0, 0, 128, 63, 0, 0, 128, 63}

	It("should set LOD field", func() {
		eq.Parse(params)
		Expect(eq.LOD).To(Equal(params["lod"]))
	})

	It("should be nil on wrong input", func() {
		rslt, _ := eq2.Encode(params_2)
		Expect(rslt).To(Equal(result))
	})
})
