package tile_test

import (
	"github.com/unchartedsoftware/prism/tile"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TopHits", func() {
	eq := &tile.TopHits{}
	eq2 := &tile.TopHits{}

	// create params
	// we use the built in `make` function to allocate the map
	params := make(map[string]interface{})
	a := make([]string, 2)
	a[0] = "first"
	a[1] = "second"
	params["sortField"] = "field"
	params["sortOrder"] = "order"
	params["hitsCount"] = 2.0
	params["includeFields"] = a

	params_fail := make(map[string]interface{})

	It("should set fields", func() {
		ok := eq.Parse(params)
		Expect(ok).To(BeNil())
		Expect(eq.SortField).To(Equal(params["sortField"]))
		Expect(eq.SortOrder).To(Equal(params["sortOrder"]))
		Expect(eq.HitsCount).To(Equal(2))
		//Expect(eq.IncludeFields).To(Equal(a))
	})

	It("should fail on wrong input", func() {
		ok := eq2.Parse(params_fail)
		Expect(ok).NotTo(BeNil())
	})
})
