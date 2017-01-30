package tile_test

import (
	"github.com/unchartedsoftware/prism/tile"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TopTerms", func() {
	eq := &tile.TopTerms{}
	eq2 := &tile.TopTerms{}

	// create params
	// we use the built in `make` function to allocate the map
	params := make(map[string]interface{})

	params["termsCount"] = 1.0
	params["termsField"] = "age"

	params_fail := make(map[string]interface{})

	It("should set Field and Value", func() {
		ok := eq.Parse(params)
		Expect(ok).To(BeNil())
		Expect(eq.TermsField).To(Equal(params["termsField"]))
		Expect(eq.TermsCount).To(Equal(1))
	})

	It("should fail on wrong input", func() {
		ok := eq2.Parse(params_fail)
		Expect(ok).NotTo(BeNil())
	})
})
