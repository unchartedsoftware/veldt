package query_test

import (
	"github.com/unchartedsoftware/prism/query"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Exists", func() {
	ex := &query.Exists{}
	ex2 := &query.Exists{}

	params := make(map[string]interface{})
	params["field"] = "field"

	params_fail := make(map[string]interface{})

	It("should set Field", func() {
		ok := ex.Parse(params)
		Expect(ok).To(BeNil())
		Expect(ex.Field).To(Equal("field"))
	})

	It("should fail on missing field", func() {
		ok := ex2.Parse(params_fail)
		Expect(ok).NotTo(BeNil())
		Expect(ex2.Field).To(Equal(""))
	})
})
