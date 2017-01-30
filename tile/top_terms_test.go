package tile_test

import (
	"github.com/unchartedsoftware/veldt/tile"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TopTerms", func() {
	tt := &tile.TopTerms{}
	tt2 := &tile.TopTerms{}

	params := make(map[string]interface{})
	params["termsCount"] = 1
	params["termsField"] = "age"

	params_fail := make(map[string]interface{})

	It("should set Field and Value", func() {
		ok := tt.Parse(params)
		Expect(ok).To(BeNil())
		Expect(tt.TermsField).To(Equal(params["termsField"]))
		Expect(tt.TermsCount).To(Equal(params["termsCount"]))
	})

	It("should fail on wrong input", func() {
		ok := tt2.Parse(params_fail)
		Expect(ok).NotTo(BeNil())
	})
})
