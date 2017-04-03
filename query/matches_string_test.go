package query_test

import (
	"github.com/unchartedsoftware/veldt/query"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/unchartedsoftware/veldt/util/test"
)

var _ = Describe("MatchesString", func() {

	var matchesString *query.MatchesString

	BeforeEach(func() {
		matchesString = &query.MatchesString{}
	})

	Describe("Parse", func() {
		It("should parse properties from the params argument", func() {
			params := JSON(
				`{
					"match": "string-query",
					"fields": ["a", "b", "c"]
				}`)
			err := matchesString.Parse(params)
			Expect(err).To(BeNil())
			Expect(matchesString.Match).To(Equal("string-query"))
			Expect(matchesString.Fields[0]).To(Equal("a"))
			Expect(matchesString.Fields[1]).To(Equal("b"))
			Expect(matchesString.Fields[2]).To(Equal("c"))
		})

		It("should return an error if `match` property is not specified", func() {
			params := JSON(`{}`)
			err := matchesString.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `fields` property is not specified", func() {
			params := JSON(
				`{
					"match": "stuff to match"
				}`)
			err := matchesString.Parse(params)
			Expect(err).NotTo(BeNil())
		})
	})
})
