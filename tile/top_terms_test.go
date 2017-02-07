package tile_test

import (
	"github.com/unchartedsoftware/veldt/tile"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/unchartedsoftware/veldt/util/test"
)

var _ = Describe("TopTerms", func() {

	var terms *tile.TopTerms

	BeforeEach(func() {
		terms = &tile.TopTerms{}
	})

	Describe("Parse", func() {
		It("should parse properties from the params argument", func() {
			params := JSON(
				`{
					"termsCount": 1,
					"termsField": "age"
				}
				`)
			ok := terms.Parse(params)
			Expect(ok).To(BeNil())
			Expect(terms.TermsCount).To(Equal(1))
			Expect(terms.TermsField).To(Equal("age"))
		})

		It("should return an error if `termsField` property is not specified", func() {
			params := JSON(`{}`)
			err := terms.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `termsCount` property is not specified", func() {
			params := JSON(
				`{
					"termsField": "age"
				}
				`)
			err := terms.Parse(params)
			Expect(err).NotTo(BeNil())
		})
	})
})
