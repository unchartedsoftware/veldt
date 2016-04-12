package param_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/unchartedsoftware/prism/generation/elastic/param"
	"gopkg.in/olivere/elastic.v3"
)

var _ = Describe("bool_query", func() {

	const ()

	var (
		BooleanQuery1 param.BoolQuery
		BooleanQuery2 param.BoolQuery
	)

	BeforeEach(func() {
		termsField := "feature.author"
		terms := [...]string{"Penn", "Teller"}
		tq := elastic.NewTermsQuery(termsField, terms)

		rangeField := "feature.numeric"
		from := 0.5
		to := 1.0
		rq := elastic.NewRangeQuery(rangeField).From(from).To(to)
		BooleanQuery1.Query = elastic.NewBoolQuery().Must(tq, rq)

		BooleanQuery2.Query = elastic.NewBoolQuery().Must(tq)
	})

	Describe("GetHash", func() {
		It("Should hash correctly for multiple queries in must clause", func() {
			hash := BooleanQuery1.GetHash()
			Expect(hash).To(Equal("feature.author:[Penn Teller]:feature.numeric:0.5:1"))
		})

		It("Should hash correctly for single query in must clause", func() {
			hash := BooleanQuery2.GetHash()
			Expect(hash).To(Equal("feature.author:[Penn Teller]"))
		})
	})
})
