package geometry_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/unchartedsoftware/veldt/geometry"
)

var _ = Describe("geomery.Coord", func() {
	x := 1.0
	y := 2.0

	It("should construct correctly", func() {
		coord := geometry.NewCoord(x, y)
		Expect(coord.X).To(Equal(x))
		Expect(coord.Y).To(Equal(y))
	})

})
