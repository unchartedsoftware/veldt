package tile_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/unchartedsoftware/veldt/util/test"

	"github.com/unchartedsoftware/veldt/tile"
)

var _ = Describe("TargetTerms", func() {

	var terms *tile.TargetTerms

	BeforeEach(func() {
		terms = &tile.TargetTerms{}
	})

	Describe("Parse", func() {
		It("should parse properties from the params argument", func() {
			params := JSON(
				`{
					"termsField": "age",
					"terms": ["one", "two"]
				}`)
			err := terms.Parse(params)
			Expect(err).To(BeNil())
			Expect(terms.TermsField).To(Equal(params["termsField"]))
			Expect(terms.Terms[0]).To(Equal("one"))
			Expect(terms.Terms[1]).To(Equal("two"))
		})

		It("should return an error if `termsField` property is not specified", func() {
			params := JSON(`{}`)
			err := terms.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `terms` property is not specified or empty", func() {
			paramsA := JSON(`{}`)
			paramsB := JSON(
				`{
					"terms": []
				}`)
			errA := terms.Parse(paramsA)
			Expect(errA).NotTo(BeNil())
			errB := terms.Parse(paramsB)
			Expect(errB).NotTo(BeNil())
		})
	})
})
