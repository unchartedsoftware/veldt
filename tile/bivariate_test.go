package tile_test

import (
	"github.com/unchartedsoftware/veldt/tile"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/unchartedsoftware/veldt/util/test"
)

var _ = Describe("Bivariate", func() {

	var bivariate *tile.Bivariate

	BeforeEach(func() {
		bivariate = &tile.Bivariate{}
	})

	Describe("Parse", func() {
		It("should parse properties from the params argument", func() {
			params := JSON(
				`{
					"xField": "x",
					"yField": "y",
					"left": -1.0,
					"right": 1.0,
					"bottom": -1.0,
					"top": 1.0
				}`)
			err := bivariate.Parse(params)
			Expect(err).To(BeNil())
			Expect(bivariate.XField).To(Equal("x"))
			Expect(bivariate.YField).To(Equal("y"))
			Expect(bivariate.WorldBounds.Left()).To(Equal(-1.0))
			Expect(bivariate.WorldBounds.Right()).To(Equal(1.0))
			Expect(bivariate.WorldBounds.Bottom()).To(Equal(-1.0))
			Expect(bivariate.WorldBounds.Top()).To(Equal(1.0))
		})

		It("should return an error if `xField` property is not specified", func() {
			params := JSON(`{}`)
			err := bivariate.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `yField` property is not specified", func() {
			params := JSON(
				`{
					"xField": "x"
				}
				`)
			err := bivariate.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `left` property is not specified", func() {
			params := JSON(
				`{
					"xField": "x",
					"yField": "y"
				}
				`)
			err := bivariate.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `right` property is not specified", func() {
			params := JSON(
				`{
					"xField": "x",
					"yField": "y",
					"left": -1.0
				}
				`)
			err := bivariate.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `bottom` property is not specified", func() {
			params := JSON(
				`{
					"xField": "x",
					"yField": "y",
					"left": -1.0,
					"right": 1.0
				}
				`)
			err := bivariate.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `top` property is not specified", func() {
			params := JSON(
				`{
					"xField": "x",
					"yField": "y",
					"left": -1.0,
					"right": 1.0,
					"bottom": -1.0
				}
				`)
			err := bivariate.Parse(params)
			Expect(err).NotTo(BeNil())
		})
	})
})
