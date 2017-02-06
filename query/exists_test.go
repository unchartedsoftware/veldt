package query_test

import (
	"github.com/unchartedsoftware/veldt/query"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Exists", func() {

	var exists *query.Exists
	var params map[string]interface{}

	BeforeEach(func() {
		exists = &query.Exists{}
		params = make(map[string]interface{})
	})

	Describe("Parse", func() {
		It("should parse properties from the params argument", func() {
			params["field"] = "field"
			err := exists.Parse(params)
			Expect(err).To(BeNil())
			Expect(exists.Field).To(Equal("field"))
		})

		It("should return an error if `field` property is not specified", func() {
			err := exists.Parse(params)
			Expect(err).NotTo(BeNil())
		})
	})
})
