package veldt_test

import (
	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/query"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Boolean", func() {
	be := &veldt.BinaryExpression{}
	ue := &veldt.UnaryExpression{}
	be_f := &veldt.BinaryExpression{}
	ue_f := &veldt.UnaryExpression{}

	l := &query.Exists{}
	r := &query.Exists{}
	// create params
	// we use the built in `make` function to allocate the map
	params := make(map[string]interface{})
	params["left"] = l
	params["right"] = r
	params["query"] = l
	params["op"] = "AND"

	params_fail := make(map[string]interface{})

	It("should set binary ops", func() {
		ok := be.Parse(params)
		Expect(ok).To(BeNil())
		Expect(be.Right).To(Equal(r))
		Expect(be.Left).To(Equal(l))
		Expect(be.Op).To(Equal("AND"))
	})

	It("should fail on wrong binary input", func() {
		ok := be_f.Parse(params_fail)
		Expect(ok).NotTo(BeNil())
	})

	It("should set unary ops", func() {
		ok := ue.Parse(params)
		Expect(ok).To(BeNil())
		Expect(ue.Query).To(Equal(l))
		Expect(ue.Op).To(Equal("AND"))
	})

	It("should fail on wrong unary input", func() {
		ok := ue_f.Parse(params_fail)
		Expect(ok).NotTo(BeNil())
	})
})
