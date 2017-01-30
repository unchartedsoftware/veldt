package query_test

import (
	"github.com/unchartedsoftware/prism/query"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Exists", func() {
	eq := &query.Exists{}
	eq2 := &query.Exists{}

	// create params
	// we use the built in `make` function to allocate the map
	params := make(map[string]interface{})
	params["field"] = "field"

	params_fail := make(map[string]interface{})

	It("should set Field", func() {
		ok := eq.Parse(params)
		Expect(eq.Field).To(Equal("field"))
		Expect(ok).To(BeNil())
	})

	It("should fail on missing field", func() {
		ok := eq2.Parse(params_fail)
		Expect(eq2.Field).To(Equal(""))
		Expect(ok).NotTo(BeNil())
	})
})
