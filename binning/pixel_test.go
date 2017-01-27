package binning_test

import (
	"math"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/unchartedsoftware/veldt/binning"
)

var _ = Describe("pixel", func() {

	const (
		epsilon = 0.000001
	)

	var (
		maxPixels = uint64(float64(binning.MaxTileResolution) *
			math.Pow(2, float64(binning.MaxLevelSupported)))
		bottomLeftLonLat = binning.LonLat{
			Lon: -180,
			Lat: -85.05112878,
		}
		centerLonLat = binning.LonLat{
			Lon: 0,
			Lat: 0,
		}
		topRightLonLat = binning.LonLat{
			Lon: 180,
			Lat: 85.05112878,
		}
		bottomLeftCoord = binning.Coord{
			X: -1,
			Y: -1,
		}
		centerCoord = binning.Coord{
			X: 0,
			Y: 0,
		}
		topRightCoord = binning.Coord{
			X: 1,
			Y: 1,
		}
	)

	Describe("CoordToPixelCoord", func() {
		It("should return a fractional tile coordinate", func() {

			extent := binning.Bounds{
				BottomLeft: &binning.Coord{
					X: -1,
					Y: -1,
				},
				TopRight: &binning.Coord{
					X: 1,
					Y: 1,
				},
			}

			pixel := binning.CoordToPixelCoord(&bottomLeftCoord, &extent)
			Expect(pixel.X).To(Equal(uint64(0)))
			Expect(pixel.Y).To(Equal(uint64(0)))

			pixel = binning.CoordToPixelCoord(&centerCoord, &extent)
			Expect(pixel.X).To(Equal(uint64(maxPixels / 2)))
			Expect(pixel.Y).To(Equal(uint64(maxPixels / 2)))

			pixel = binning.CoordToPixelCoord(&topRightCoord, &extent)
			Expect(pixel.X).To(Equal(uint64(maxPixels - 1)))
			Expect(pixel.Y).To(Equal(uint64(maxPixels - 1)))
		})
	})

	Describe("LonLatToPixelCoord", func() {
		It("should return a tile coordinate", func() {
			pixel := binning.LonLatToPixelCoord(&bottomLeftLonLat)
			Expect(pixel.X).To(Equal(uint64(0)))
			Expect(pixel.X).To(Equal(uint64(0)))

			pixel = binning.LonLatToPixelCoord(&centerLonLat)
			Expect(pixel.X).To(Equal(uint64(maxPixels / 2)))
			Expect(pixel.Y).To(Equal(uint64(maxPixels / 2)))

			pixel = binning.LonLatToPixelCoord(&topRightLonLat)
			Expect(pixel.X).To(Equal(uint64(maxPixels - 1)))
			Expect(pixel.Y).To(Equal(uint64(maxPixels - 1)))
		})
	})

})
