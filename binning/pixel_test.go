package binning_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/unchartedsoftware/veldt/binning"
)

var _ = Describe("pixel", func() {

	var (
		bottomLeftLonLat = binning.NewLonLat(-180, -85.05112878)
		centerLonLat = binning.NewLonLat(0, 0)
		topRightLonLat = binning.NewLonLat(180, 85.05112878)
		bottomLeftCoord = binning.NewCoord(-1, -1)
		centerCoord = binning.NewCoord(0, 0)
		topRightCoord = binning.NewCoord(1, 1)
		centerPixelCoord = binning.NewPixelCoord(
			uint64(binning.MaxPixels / 2),
			uint64(binning.MaxPixels / 2))
		topRightPixel = binning.NewPixelCoord(
			uint64(binning.MaxPixels - 1),
			uint64(binning.MaxPixels - 1))
	)

	Describe("CoordToPixelCoord", func() {
		It("should return a fractional tile coordinate", func() {

			extent := binning.Bounds{
				BottomLeft: binning.NewCoord(-1, -1),
				TopRight: binning.NewCoord(1, 1),
			}

			pixel := binning.CoordToPixelCoord(bottomLeftCoord, &extent)
			Expect(pixel.X).To(Equal(uint64(0)))
			Expect(pixel.Y).To(Equal(uint64(0)))

			pixel = binning.CoordToPixelCoord(centerCoord, &extent)
			Expect(pixel.X).To(Equal(centerPixelCoord.X))
			Expect(pixel.Y).To(Equal(centerPixelCoord.Y))

			pixel = binning.CoordToPixelCoord(topRightCoord, &extent)
			Expect(pixel.X).To(Equal(topRightPixel.X))
			Expect(pixel.Y).To(Equal(topRightPixel.Y))
		})
	})

	Describe("LonLatToPixelCoord", func() {
		It("should return a tile coordinate", func() {
			pixel := binning.LonLatToPixelCoord(bottomLeftLonLat)
			Expect(pixel.X).To(Equal(uint64(0)))
			Expect(pixel.X).To(Equal(uint64(0)))

			pixel = binning.LonLatToPixelCoord(centerLonLat)
			Expect(pixel.X).To(Equal(centerPixelCoord.X))
			Expect(pixel.Y).To(Equal(centerPixelCoord.Y))

			pixel = binning.LonLatToPixelCoord(topRightLonLat)
			Expect(pixel.X).To(Equal(topRightPixel.X))
			Expect(pixel.Y).To(Equal(topRightPixel.Y))
		})
	})

})
