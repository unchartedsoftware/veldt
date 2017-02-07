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

	Describe("Corners", func() {
		It("should return an equivalent set of corners", func() {

			corners := bounds.Corners()
			Expect(corners.BottomLeft.X).To(Equal(left))
			Expect(corners.BottomLeft.Y).To(Equal(bottom))
			Expect(corners.TopRight.X).To(Equal(right))
			Expect(corners.TopRight.Y).To(Equal(top))
		})
	})

})
