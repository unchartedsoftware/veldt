package query_test

import (
	"github.com/unchartedsoftware/veldt/query"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Range", func() {
	rng := &query.Range{}
	rng2 := &query.Range{}
	rng3 := &query.Range{}
	rng4 := &query.Range{}

	params := make(map[string]interface{})
	params["field"] = "range"
	params["gte"] = true
	params["lte"] = true

	params_fail := make(map[string]interface{})

	params_fail_upper := make(map[string]interface{})
	params_fail_upper["field"] = "range"
	params_fail_upper["gte"] = true
	params_fail_upper["gt"] = true

	params_fail_lower := make(map[string]interface{})
	params_fail_lower["field"] = "range"
	params_fail_lower["lt"] = true
	params_fail_lower["lte"] = true

	It("should set Field and Value", func() {
		ok := rng.Parse(params)
		Expect(ok).To(BeNil())
		Expect(rng.Field).To(Equal("range"))
		Expect(rng.GTE).To(Equal(true))
	})

	It("should fail on wrong input", func() {
		ok := rng2.Parse(params_fail)
		Expect(ok).NotTo(BeNil())
	})

	It("should fail if gt and gte are both set", func() {
		ok := rng3.Parse(params_fail_upper)
		Expect(ok).NotTo(BeNil())
	})

	It("should fail if lt and lte are both set", func() {
		ok := rng4.Parse(params_fail_lower)
		Expect(ok).NotTo(BeNil())
	})
})
