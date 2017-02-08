package geometry_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/unchartedsoftware/veldt/geometry"
)

var _ = Describe("geometry.Bounds", func() {

	left := -180.0
	right := 180.0
	bottom := -90.0
	top := 90.0
	bounds := geometry.NewBounds(left, right, bottom, top)

	Describe("Constructors", func() {
		It("should return sides correctly", func() {
			Expect(bounds.Left()).To(Equal(left))
			Expect(bounds.Right()).To(Equal(right))
			Expect(bounds.Bottom()).To(Equal(bottom))
			Expect(bounds.Top()).To(Equal(top))
		})

		It("should construct correctly from a Rectangle", func() {
			b := geometry.NewBoundsFromRectangle(
				&geometry.Rectangle{
					BottomLeft: geometry.NewCoord(left, bottom),
					TopRight:   geometry.NewCoord(right, top),
				},
			)
			Expect(b.Left()).To(Equal(left))
			Expect(b.Right()).To(Equal(right))
			Expect(b.Bottom()).To(Equal(bottom))
			Expect(b.Top()).To(Equal(top))
		})

		It("should contruct via parsing param map", func() {
			params := make(map[string]interface{})
			params["left"] = left
			params["right"] = right
			params["bottom"] = bottom
			params["top"] = top
			b, _ := geometry.NewBoundsByParse(params)
			Expect(b.Left()).To(Equal(left))
			Expect(b.Right()).To(Equal(right))
			Expect(b.Bottom()).To(Equal(bottom))
			Expect(b.Top()).To(Equal(top))
		})
	})

	Describe("Parsing validation", func() {
		It("should return error for any missing side", func() {
			params := make(map[string]interface{})

			_, err := geometry.NewBoundsByParse(params)
			Expect(err).To(HaveOccurred())

			params["left"] = left
			_, err2 := geometry.NewBoundsByParse(params)
			Expect(err2).To(HaveOccurred())

			params["right"] = right
			_, err3 := geometry.NewBoundsByParse(params)
			Expect(err3).To(HaveOccurred())

			params["bottom"] = bottom
			_, err4 := geometry.NewBoundsByParse(params)
			Expect(err4).To(HaveOccurred())

		})
	})

	Describe("Calculated values", func() {
		It("should return an equivalent set of corners", func() {
			corners := bounds.Corners()
			Expect(corners.BottomLeft.X).To(Equal(left))
			Expect(corners.BottomLeft.Y).To(Equal(bottom))
			Expect(corners.TopRight.X).To(Equal(right))
			Expect(corners.TopRight.Y).To(Equal(top))
		})
		It("should have equivalent direct corner accessors", func() {
			Expect(bounds.BottomLeft().X).To(Equal(left))
			Expect(bounds.BottomLeft().Y).To(Equal(bottom))
			Expect(bounds.TopRight().X).To(Equal(right))
			Expect(bounds.TopRight().Y).To(Equal(top))
		})
		It("should calculate the mathematical extrema", func() {
			b := geometry.NewBounds(1.0, -1.0, 1.0, -1.0)
			Expect(b.MinX()).To(Equal(-1.0))
			Expect(b.MaxX()).To(Equal(1.0))
			Expect(b.MinY()).To(Equal(-1.0))
			Expect(b.MaxY()).To(Equal(1.0))

			// Ranges
			Expect(b.RangeX()).To(Equal(2.0))
			Expect(b.RangeY()).To(Equal(2.0))
		})
	})

})
