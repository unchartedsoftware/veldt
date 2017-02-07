package veldt_test

import (
	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/query"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BinaryExpression", func() {

	var binary *veldt.BinaryExpression
	var params map[string]interface{}

	BeforeEach(func() {
		binary = &veldt.BinaryExpression{}
		params = make(map[string]interface{})
	})

	Describe("Parse", func() {
		It("should parse properties from the params argument", func() {
			params["left"] = &query.Exists{}
			params["op"] = "AND"
			params["right"] = &query.Exists{}
			err := binary.Parse(params)
			Expect(err).To(BeNil())
			Expect(binary.Left).To(Equal(params["left"]))
			Expect(binary.Op).To(Equal(params["op"]))
			Expect(binary.Right).To(Equal(params["right"]))
		})

		It("should return an error if `left` property is not specified", func() {
			err := binary.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `left` property is not of type `Query`", func() {
			params["left"] = 5.0
			err := binary.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `op` property is not specified", func() {
			params["left"] = &query.Exists{}
			err := binary.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `op` property is not recognized", func() {
			params["left"] = &query.Exists{}
			params["op"] = "INVALID"
			err := binary.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `right` property is not specified", func() {
			params["left"] = &query.Exists{}
			params["op"] = "AND"
			err := binary.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `right` property is not of type `Query`", func() {
			params["left"] = &query.Exists{}
			params["op"] = "AND"
			params["right"] = 5.0
			err := binary.Parse(params)
			Expect(err).NotTo(BeNil())
		})
	})
})

var _ = Describe("UnaryExpression", func() {

	var unary *veldt.UnaryExpression
	var params map[string]interface{}

	BeforeEach(func() {
		unary = &veldt.UnaryExpression{}
		params = make(map[string]interface{})
	})

	Describe("Parse", func() {
		It("should parse properties from the params argument", func() {
			params["op"] = "NOT"
			params["query"] = &query.Exists{}
			err := unary.Parse(params)
			Expect(err).To(BeNil())
			Expect(unary.Op).To(Equal(params["op"]))
			Expect(unary.Query).To(Equal(params["query"]))
		})

		It("should return an error if `query` property is not specified", func() {
			err := unary.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `query` property is not of type `Query`", func() {
			params["query"] = 5.0
			err := unary.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `op` property is not specified", func() {
			params["query"] = &query.Exists{}
			err := unary.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `op` property is not recognized", func() {
			params["query"] = &query.Exists{}
			params["op"] = "INVALID"
			err := unary.Parse(params)
			Expect(err).NotTo(BeNil())
		})
	})
})
