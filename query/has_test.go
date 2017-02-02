package query_test

import (
	"github.com/unchartedsoftware/prism/query"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Has", func() {
	has := &query.Has{}
	has2 := &query.Has{}

	params := make(map[string]interface{})
	a := make([]interface{}, 2, 2)
	a[0] = "first"
	a[1] = "second"
	params["values"] = a
	params["field"] = "age"

	params_fail := make(map[string]interface{})

	It("should set Field and Value", func() {
		ok := has.Parse(params)
		Expect(ok).To(BeNil())
		Expect(has.Field).To(Equal(params["field"]))
		Expect(has.Values).ToEqual(params["values"])
	})

	It("should fail on wrong input", func() {
		ok := has2.Parse(params_fail)
		Expect(ok).NotTo(BeNil())
	})
})
