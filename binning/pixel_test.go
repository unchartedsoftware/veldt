package binning_test

import (
	"math"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/unchartedsoftware/prism/binning"
)

var _ = Describe("pixel", func() {

	const (
		epsilon = 0.000001
	)

	var (
		maxPixels = uint64(float64(binning.MaxTileResolution) *
			math.Pow(2, float64(binning.MaxLevelSupported)))
		t0 = binning.TileCoord{
			X: 0,
			Y: 0,
			Z: 0,
		}
		t1 = binning.TileCoord{
			X: 0,
			Y: 0,
			Z: 1,
		}
		t2 = binning.TileCoord{
			X: 1,
			Y: 1,
			Z: 1,
		}
		lonlat0 = binning.LonLat{
			Lon: -180,
			Lat: 85.05112878,
		}
		lonlat1 = binning.LonLat{
			Lon: 0,
			Lat: 0,
		}
		lonlat2 = binning.LonLat{
			Lon: 180,
			Lat: -85.05112878,
		}
		xy0 = binning.Coord{
			X: -1,
			Y: -1,
		}
		xy1 = binning.Coord{
			X: 0,
			Y: 0,
		}
		xy2 = binning.Coord{
			X: 1,
			Y: 1,
		}
	)

	Describe("GetTilePixelBounds", func() {
		It("should return geo bounds from a tile coordinate", func() {
			bounds := binning.GetTilePixelBounds(&t0)
			Expect(bounds.TopLeft.X).To(Equal(uint64(0)))
			Expect(bounds.TopLeft.Y).To(Equal(uint64(0)))
			Expect(bounds.BottomRight.X).To(Equal(maxPixels - 1))
			Expect(bounds.BottomRight.Y).To(Equal(maxPixels - 1))

			bounds = binning.GetTilePixelBounds(&t1)
			Expect(bounds.TopLeft.X).To(Equal(uint64(0)))
			Expect(bounds.TopLeft.Y).To(Equal(uint64(0)))
			Expect(bounds.BottomRight.X).To(Equal(maxPixels / 2))
			Expect(bounds.BottomRight.Y).To(Equal(maxPixels / 2))

			bounds = binning.GetTilePixelBounds(&t2)
			Expect(bounds.TopLeft.X).To(Equal(maxPixels / 2))
			Expect(bounds.TopLeft.Y).To(Equal(maxPixels / 2))
			Expect(bounds.BottomRight.X).To(Equal(maxPixels - 1))
			Expect(bounds.BottomRight.Y).To(Equal(maxPixels - 1))
		})
	})

	Describe("CoordToPixelCoord", func() {
		It("should return a fractional tile coordinate", func() {

			extent := binning.Bounds{
				TopLeft: &binning.Coord{
					X: -1,
					Y: -1,
				},
				BottomRight: &binning.Coord{
					X: 1,
					Y: 1,
				},
			}

			pixel := binning.CoordToPixelCoord(&xy0, &extent)
			Expect(pixel.X).To(Equal(uint64(0)))
			Expect(pixel.Y).To(Equal(uint64(0)))

			pixel = binning.CoordToPixelCoord(&xy1, &extent)
			Expect(pixel.X).To(Equal(uint64(maxPixels / 2)))
			Expect(pixel.Y).To(Equal(uint64(maxPixels / 2)))

			pixel = binning.CoordToPixelCoord(&xy2, &extent)
			Expect(pixel.X).To(Equal(uint64(maxPixels - 1)))
			Expect(pixel.Y).To(Equal(uint64(maxPixels - 1)))
		})
	})

	Describe("LonLatToPixelCoord", func() {
		It("should return a tile coordinate", func() {
			pixel := binning.LonLatToPixelCoord(&lonlat0)
			Expect(pixel.X).To(Equal(uint64(0)))
			Expect(pixel.X).To(Equal(uint64(0)))

			pixel = binning.LonLatToPixelCoord(&lonlat1)
			Expect(pixel.X).To(Equal(uint64(maxPixels / 2)))
			Expect(pixel.Y).To(Equal(uint64(maxPixels / 2)))

			pixel = binning.LonLatToPixelCoord(&lonlat2)
			Expect(pixel.X).To(Equal(uint64(maxPixels - 1)))
			Expect(pixel.Y).To(Equal(uint64(maxPixels - 1)))
		})
	})

})
