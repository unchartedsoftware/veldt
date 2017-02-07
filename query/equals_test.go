package query_test

import (
	"github.com/unchartedsoftware/veldt/query"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/unchartedsoftware/veldt/util/test"
)

var _ = Describe("Parse", func() {

	var equals *query.Equals

	BeforeEach(func() {
		equals = &query.Equals{}
	})

	Describe("Parse", func() {
		It("should parse properties from the params argument", func() {
			params := JSON(
				`{
					"field": "field",
					"value": "value"
				}`)
			err := equals.Parse(params)
			Expect(err).To(BeNil())
			Expect(equals.Field).To(Equal("field"))
			Expect(equals.Value).To(Equal("value"))
		})

		It("should return an error if `field` property is not specified", func() {
			params := JSON(`{}`)
			err := equals.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `value` property is not specified", func() {
			params := JSON(
				`{
					"field": "field"
				}`)
			err := equals.Parse(params)
			Expect(err).NotTo(BeNil())
		})
	})

})
