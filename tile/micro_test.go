package tile_test

import (
	"github.com/unchartedsoftware/veldt/tile"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Micro", func() {
	eq := &tile.Micro{}

	// create params
	// we use the built in `make` function to allocate the map
	params := make(map[string]interface{})
	params["lod"] = 1

	It("should set LOD", func() {
		ok := eq.Parse(params)
		Expect(ok).To(BeNil())
		Expect(eq.LOD).To(Equal(params["lod"]))

	})

	It("should set LOD 2", func() {
		result := eq.ParseIncludes([]string{"a", "b"}, "a", "b")
		Expect(result).To(Equal([]string{"a", "b"}))
	})

})
