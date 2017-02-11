package query_test

import (
	"github.com/unchartedsoftware/veldt/query"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/unchartedsoftware/veldt/util/test"
)

var _ = Describe("Range", func() {

	var rang *query.Range

	BeforeEach(func() {
		rang = &query.Range{}
	})

	Describe("Parse", func() {
		It("should parse properties from the params argument", func() {
			params := JSON(
				`{
					"field": "field",
					"gte": 0.0,
					"lt": 256.0
				}`)
			err := rang.Parse(params)
			Expect(err).To(BeNil())
			Expect(rang.Field).To(Equal("field"))
			Expect(rang.GTE).To(Equal(0.0))
			Expect(rang.LT).To(Equal(256.0))
		})

		It("should return an error if `field` property is not specified", func() {
			params := JSON(`{}`)
			err := rang.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if both `gte` and `gt` are specified", func() {
			params := JSON(
				`{
					"field": "field",
					"gte": 0.0,
					"gt": 0.0
				}`)
			err := rang.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if both `lte` and `lt` are specified", func() {
			params := JSON(
				`{
					"field": "field",
					"gte": 0.0,
					"lte": 256.0,
					"lt": 256.0
				}`)
			err := rang.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `gte`, `gt`, `lte`, and `lt` are not specified", func() {
			params := JSON(
				`{
					"field": "field"
				}`)
			err := rang.Parse(params)
			Expect(err).NotTo(BeNil())
		})
	})

})
