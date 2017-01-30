package query_test

import (
	"github.com/unchartedsoftware/prism/query"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parse", func() {
	eq := &query.Equals{}
	eq2 := &query.Equals{}

	// create params
	// we use the built in `make` function to allocate the map
	params := make(map[string]interface{})
	params["value"] = 35
	params["field"] = "age"

	params_fail := make(map[string]interface{})

	It("should set Field and Value", func() {
		ok := eq.Parse(params)
		Expect(eq.Field).To(Equal("age"))
		Expect(eq.Value).To(Equal(35))
		Expect(ok).To(BeNil())
	})

	It("should fail on wrong input", func() {
		ok := eq2.Parse(params_fail)
		Expect(eq2.Field).To(Equal(""))
		Expect(eq2.Value).To(BeNil())
		Expect(ok).NotTo(Equal(nil))
	})
})
