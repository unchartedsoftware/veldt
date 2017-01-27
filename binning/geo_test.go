package binning_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/unchartedsoftware/veldt/binning"
)

var _ = Describe("geo", func() {

	const (
		epsilon = 0.000001
	)

	var (
		bottomLeft = binning.LonLat{
			Lon: -180,
			Lat: -85.05112878,
		}
		center = binning.LonLat{
			Lon: 0,
			Lat: 0,
		}
		topRight = binning.LonLat{
			Lon: 180,
			Lat: 85.05112878,
		}
	)

	Describe("LonLatToFractionalTile", func() {
		It("should return a fractional tile coordinate", func() {
			tile := binning.LonLatToFractionalTile(&bottomLeft, 0)
			Expect(tile.X).To(BeNumerically("~", 0.0, epsilon))
			Expect(tile.Y).To(BeNumerically("~", 0.0, epsilon))

			tile = binning.LonLatToFractionalTile(&center, 1)
			Expect(tile.X).To(BeNumerically("~", 1.0, epsilon))
			Expect(tile.Y).To(BeNumerically("~", 1.0, epsilon))

			tile = binning.LonLatToFractionalTile(&topRight, 1)
			Expect(tile.X).To(BeNumerically("~", 2.0, epsilon))
			Expect(tile.Y).To(BeNumerically("~", 2.0, epsilon))
		})
	})

})
