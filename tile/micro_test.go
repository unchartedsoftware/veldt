package tile_test

import (
	"github.com/unchartedsoftware/veldt/tile"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Micro", func() {
	mcr := &tile.Micro{}

	params := make(map[string]interface{})
	params["lod"] = 1

	It("should set LOD", func() {
		ok := mcr.Parse(params)
		Expect(ok).To(BeNil())
		Expect(mcr.LOD).To(Equal(params["lod"]))

	})

	It("should set parse includes correctly", func() {
		result := mcr.ParseIncludes([]string{"a", "b"}, "a", "b")
		Expect(result).To(Equal([]string{"a", "b"}))
	})

})
