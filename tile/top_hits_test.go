package tile_test

import (
	"github.com/unchartedsoftware/veldt/tile"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TopHits", func() {
	th := &tile.TopHits{}
	th2 := &tile.TopHits{}

	params := make(map[string]interface{})
	a := []string{"first", "second"}

	params["sortField"] = "field"
	params["sortOrder"] = "order"
	params["hitsCount"] = 2
	params["includeFields"] = a

	params_fail := make(map[string]interface{})

	It("should set fields", func() {
		ok := th.Parse(params)
		Expect(ok).To(BeNil())
		Expect(th.SortField).To(Equal(params["sortField"]))
		Expect(th.SortOrder).To(Equal(params["sortOrder"]))
		Expect(th.HitsCount).To(Equal(params["hitsCount"]))
		//Expect(eq.IncludeFields).To(Equal(a))
	})

	It("should fail on wrong input", func() {
		ok := th2.Parse(params_fail)
		Expect(ok).NotTo(BeNil())
	})
})
