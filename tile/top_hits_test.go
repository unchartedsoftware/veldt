package tile_test

import (
	"github.com/unchartedsoftware/veldt/tile"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/unchartedsoftware/veldt/util/test"
)

var _ = Describe("TopHits", func() {

	var hits *tile.TopHits

	BeforeEach(func() {
		hits = &tile.TopHits{}
	})

	Describe("Parse", func() {
		It("should parse properties from the params argument", func() {
			params := JSON(
				`{
					"sortField": "field",
					"sortOrder": "desc",
					"hitsCount": 2,
					"includeFields": ["a", "b", "c"]
				}
				`)
			err := hits.Parse(params)
			Expect(err).To(BeNil())
			Expect(hits.SortField).To(Equal("field"))
			Expect(hits.SortOrder).To(Equal("desc"))
			Expect(hits.HitsCount).To(Equal(2))
			Expect(hits.IncludeFields[0]).To(Equal("a"))
			Expect(hits.IncludeFields[1]).To(Equal("b"))
			Expect(hits.IncludeFields[2]).To(Equal("c"))
		})

		It("should return an error if `sortField` property is not specified", func() {
			params := JSON(`{}`)
			err := hits.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `sortOrder` property is not specified", func() {
			params := JSON(
				`{
					"sortField": "field"
				}
				`)
			err := hits.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `sortOrder` property is not `desc` or `asc`", func() {
			params := JSON(
				`{
					"sortField": "field",
					"sortOrder": "invalid"
				}
				`)
			err := hits.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `hitsCount` property is not specified", func() {
			params := JSON(
				`{
					"sortField": "field",
					"sortOrder": "desc"
				}
				`)
			err := hits.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should gracefully handle empty `includeFields` array", func() {
			params := JSON(
				`{
					"hitsCount": 2,
					"includeFields": []
				}
				`)
			err := hits.Parse(params)
			Expect(err).To(BeNil())
		})
	})
})
