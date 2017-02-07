package query_test

import (
	"github.com/unchartedsoftware/veldt/query"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/unchartedsoftware/veldt/util/test"
)

var _ = Describe("Exists", func() {

	var exists *query.Exists

	BeforeEach(func() {
		exists = &query.Exists{}
	})

	Describe("Parse", func() {
		It("should parse properties from the params argument", func() {
			params := JSON(
				`{
					"field": "field"
				}`)
			err := exists.Parse(params)
			Expect(err).To(BeNil())
			Expect(exists.Field).To(Equal("field"))
		})

		It("should return an error if `field` property is not specified", func() {
			params := JSON(`{}`)
			err := exists.Parse(params)
			Expect(err).NotTo(BeNil())
		})
	})
})
