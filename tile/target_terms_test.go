package tile_test

import (
	"github.com/unchartedsoftware/veldt/tile"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TargetTerms", func() {
	//eq := &tile.TargetTerms{}
	eq2 := &tile.TargetTerms{}

	// create params
	// we use the built in `make` function to allocate the map
	params := make(map[string]interface{})

	a := make([]string, 2, 2)
	a[0] = "one"
	a[1] = "two"

	params["termsField"] = "age"
	params["terms"] = a

	params_fail := make(map[string]interface{})

	/*It("should set Field and Value", func() {
		ok := eq.Parse(params)
		Expect(ok).To(BeNil())
		Expect(eq.TermsField).To(Equal(params["termsField"]))
		Expect(eq.Terms).To(Equal(params["terms"]))
	})*/

	It("should fail on wrong input", func() {
		ok := eq2.Parse(params_fail)
		Expect(ok).NotTo(BeNil())
	})
})
