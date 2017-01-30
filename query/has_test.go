package query_test

import (
	"github.com/unchartedsoftware/prism/query"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Has", func() {
	eq := &query.Has{}
	eq2 := &query.Has{}

	// create params
	// we use the built in `make` function to allocate the map
	params := make(map[string]interface{})
	a := make([]interface{}, 2, 2)
	a[0] = "first"
	a[1] = "second"
	params["values"] = a
	params["field"] = "age"

	params_fail := make(map[string]interface{})

	It("should set Field and Value", func() {
		ok := eq.Parse(params)
		Expect(ok).To(BeNil())
		Expect(eq.Field).To(Equal(params["field"]))
		Expect(eq.Values).NotTo(BeNil())
	})

	It("should fail on wrong input", func() {
		ok := eq2.Parse(params_fail)
		Expect(ok).NotTo(BeNil())
	})
})
