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

	encd := make([]map[string]interface{},2)
	hit1 := make(map[string]interface{})
	hit2 := make(map[string]interface{})
	encd[0] = hit1
	encd[1] = hit2
	encd_rslt := []byte(`{"hits":null,"offsets":[0,8,8,8],"points":[1,1]}`)

	It("should set LOD", func() {
		ok := mcr.Parse(params)
		Expect(ok).To(BeNil())
		Expect(mcr.LOD).To(Equal(params["lod"]))

	})

	It("should set parse includes correctly", func() {
		result := mcr.ParseIncludes([]string{"a", "b"}, "a", "b")
		Expect(result).To(Equal([]string{"a", "b"}))
	})

	It("should set parse includes correctly", func() {
		result, ok := mcr.Encode(encd, []float32{1,1})
		Expect(ok).To(BeNil());
		Expect(result).To(Equal(encd_rslt))
	})


})
