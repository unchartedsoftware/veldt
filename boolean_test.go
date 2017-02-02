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

	left := &query.Exists{}
	right := &query.Exists{}

	params := make(map[string]interface{})
	params["left"] = left
	params["right"] = right
	params["query"] = left
	params["op"] = "AND"

	params_fail := make(map[string]interface{})

	It("should set binary ops", func() {
		ok := be.Parse(params)
		Expect(ok).To(BeNil())
		Expect(be.Right).To(Equal(params["right"]))
		Expect(be.Left).To(Equal(params["left"]))
		Expect(be.Op).To(Equal(params["op"]))
	})

	It("should fail on wrong binary input", func() {
		ok := be_f.Parse(params_fail)
		Expect(ok).NotTo(BeNil())
	})

	It("should set unary ops", func() {
		ok := ue.Parse(params)
		Expect(ok).To(BeNil())
		Expect(ue.Query).To(Equal(params["left"]))
		Expect(ue.Op).To(Equal("AND"))
	})

	It("should fail on wrong unary input", func() {
		ok := ue_f.Parse(params_fail)
		Expect(ok).NotTo(BeNil())
	})
})
