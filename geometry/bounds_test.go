package geometry_test

import (
	"math"

	"github.com/unchartedsoftware/veldt/geometry"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/unchartedsoftware/veldt/util/test"
)

var _ = Describe("Bounds", func() {

	var bounds *geometry.Bounds

	BeforeEach(func() {
		bounds = geometry.NewBounds(-1, 1, -1, 1)
	})

	Describe("Parse", func() {
		It("should parse properties from the params argument", func() {
			params := JSON(
				`{
					"left": -2.0,
					"right": 2.0,
					"bottom": -2.0,
					"top": 2.0
				}`)
			err := bounds.Parse(params)
			Expect(err).To(BeNil())
			Expect(bounds.Left).To(Equal(-2.0))
			Expect(bounds.Right).To(Equal(2.0))
			Expect(bounds.Bottom).To(Equal(-2.0))
			Expect(bounds.Top).To(Equal(2.0))
		})

		It("should return an error if `left` property is not specified", func() {
			params := JSON(`{}`)
			err := bounds.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `right` property is not specified", func() {
			params := JSON(
				`{
					"left": -1.0
				}
				`)
			err := bounds.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `bottom` property is not specified", func() {
			params := JSON(
				`{
					"left": -1.0,
					"right": 1.0
				}
				`)
			err := bounds.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `top` property is not specified", func() {
			params := JSON(
				`{
					"left": -1.0,
					"right": 1.0,
					"bottom": -1.0
				}
				`)
			err := bounds.Parse(params)
			Expect(err).NotTo(BeNil())
		})
	})

	Describe("MinX", func() {
		It("should return the minimum x value", func() {
			Expect(bounds.MinX()).To(Equal(math.Min(bounds.Left, bounds.Right)))
		})
	})

	Describe("MaxX", func() {
		It("should return the maximum x value", func() {
			Expect(bounds.MaxX()).To(Equal(math.Max(bounds.Left, bounds.Right)))
		})
	})

	Describe("MinY", func() {
		It("should return the minimum y value", func() {
			Expect(bounds.MinY()).To(Equal(math.Min(bounds.Top, bounds.Bottom)))
		})
	})

	Describe("MaxY", func() {
		It("should return the maximum y value", func() {
			Expect(bounds.MaxY()).To(Equal(math.Max(bounds.Top, bounds.Bottom)))
		})
	})

	Describe("RangeX", func() {
		It("should return the absolute distance between left and right", func() {
			Expect(bounds.RangeX()).To(Equal(math.Abs(bounds.Left - bounds.Right)))
		})
	})

	Describe("MaxY", func() {
		It("should return the absolute distance between top and bottom", func() {
			Expect(bounds.RangeY()).To(Equal(math.Abs(bounds.Top - bounds.Bottom)))
		})
	})

})
