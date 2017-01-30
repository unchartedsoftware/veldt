package query_test

import (
	"github.com/unchartedsoftware/veldt/query"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parse", func() {
	eq := &query.Equals{}
	eq2 := &query.Equals{}

	// create params
	params := make(map[string]interface{})
	params["value"] = 35
	params["field"] = "age"

	params_fail := make(map[string]interface{})

	It("should set Field and Value", func() {
		ok := eq.Parse(params)
		Expect(ok).To(BeNil())
		Expect(eq.Field).To(Equal(params["field"]))
		Expect(eq.Value).To(Equal(params["value"]))
	})

	It("should fail on wrong input", func() {
		ok := eq2.Parse(params_fail)
		Expect(ok).NotTo(Equal(nil))
		Expect(eq2.Field).To(Equal(""))
		Expect(eq2.Value).To(BeNil())
	})
})
