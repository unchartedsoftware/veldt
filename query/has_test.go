package query_test

import (
	"github.com/unchartedsoftware/veldt/query"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/unchartedsoftware/veldt/util/test"
)

var _ = Describe("Has", func() {

	var has *query.Has

	BeforeEach(func() {
		has = &query.Has{}
	})

	Describe("Parse", func() {
		It("should parse properties from the params argument", func() {
			params := JSON(
				`{
					"field": "field",
					"values": ["a", "b", "c"]
				}`)
			err := has.Parse(params)
			Expect(err).To(BeNil())
			Expect(has.Field).To(Equal("field"))
			Expect(has.Values[0]).To(Equal("a"))
			Expect(has.Values[1]).To(Equal("b"))
			Expect(has.Values[2]).To(Equal("c"))
		})

		It("should return an error if `field` property is not specified", func() {
			params := JSON(`{}`)
			err := has.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `value` property is not specified", func() {
			params := JSON(
				`{
					"field": "field"
				}`)
			err := has.Parse(params)
			Expect(err).NotTo(BeNil())
		})
	})
})
