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

	rng := &query.Range{}
	rng2 := &query.Range{}
	rng3 := &query.Range{}
	rng4 := &query.Range{}

	params := make(map[string]interface{})
	params["field"] = "range"
	params["gte"] = true
	params["lte"] = true

	params_fail := make(map[string]interface{})

	params_fail_upper := make(map[string]interface{})
	params_fail_upper["field"] = "range"
	params_fail_upper["gte"] = true
	params_fail_upper["gt"] = true

	params_fail_lower := make(map[string]interface{})
	params_fail_lower["field"] = "range"
	params_fail_lower["lt"] = true
	params_fail_lower["lte"] = true

	It("should set Field and Value", func() {
		ok := rng.Parse(params)
		Expect(ok).To(BeNil())
		Expect(rng.Field).To(Equal("range"))
		Expect(rng.GTE).To(Equal(true))
	})

	It("should fail on wrong input", func() {
		ok := rng2.Parse(params_fail)
		Expect(ok).NotTo(BeNil())
	})

	It("should fail if gt and gte are both set", func() {
		ok := rng3.Parse(params_fail_upper)
		Expect(ok).NotTo(BeNil())
	})

	It("should fail if lt and lte are both set", func() {
		ok := rng4.Parse(params_fail_lower)
		Expect(ok).NotTo(BeNil())
	})
})
