package binning_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/unchartedsoftware/prism/binning"
)

var _ = Describe("geo", func() {

	const (
		epsilon = 0.000001
	)

	var (
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
	)

	Describe("GetTileGeoBounds", func() {
		It("should return geo bounds from a tile coordinate", func() {
			bounds := binning.GetTileGeoBounds(&t0)
			Expect(bounds.TopLeft.Lon).To(BeNumerically("~", -180.0, epsilon))
			Expect(bounds.BottomRight.Lon).To(BeNumerically("~", 180.0, epsilon))
			Expect(bounds.TopLeft.Lat).To(BeNumerically("~", 85.05112878, epsilon))
			Expect(bounds.BottomRight.Lat).To(BeNumerically("~", -85.05112878, epsilon))

			bounds = binning.GetTileGeoBounds(&t1)
			Expect(bounds.TopLeft.Lon).To(BeNumerically("~", -180.0, epsilon))
			Expect(bounds.BottomRight.Lon).To(BeNumerically("~", 0.0, epsilon))
			Expect(bounds.TopLeft.Lat).To(BeNumerically("~", 85.05112878, epsilon))
			Expect(bounds.BottomRight.Lat).To(BeNumerically("~", 0.0, epsilon))

			bounds = binning.GetTileGeoBounds(&t2)
			Expect(bounds.TopLeft.Lon).To(BeNumerically("~", 0.0, epsilon))
			Expect(bounds.BottomRight.Lon).To(BeNumerically("~", 180.0, epsilon))
			Expect(bounds.TopLeft.Lat).To(BeNumerically("~", 0.0, epsilon))
			Expect(bounds.BottomRight.Lat).To(BeNumerically("~", -85.05112878, epsilon))
		})
	})

	Describe("LonLatToFractionalTile", func() {
		It("should return a fractional tile coordinate", func() {
			tile := binning.LonLatToFractionalTile(&lonlat0, 0)
			Expect(tile.X).To(BeNumerically("~", 0.0, epsilon))
			Expect(tile.Y).To(BeNumerically("~", 0.0, epsilon))

			tile = binning.LonLatToFractionalTile(&lonlat1, 1)
			Expect(tile.X).To(BeNumerically("~", 1.0, epsilon))
			Expect(tile.Y).To(BeNumerically("~", 1.0, epsilon))

			tile = binning.LonLatToFractionalTile(&lonlat2, 1)
			Expect(tile.X).To(BeNumerically("~", 2.0, epsilon))
			Expect(tile.Y).To(BeNumerically("~", 2.0, epsilon))
		})
	})

	Describe("LonLatToTile", func() {
		It("should return a tile coordinate", func() {
			tile := binning.LonLatToTile(&lonlat0, 0)
			Expect(tile.X).To(Equal(uint32(0)))
			Expect(tile.X).To(Equal(uint32(0)))

			tile = binning.LonLatToTile(&lonlat1, 1)
			Expect(tile.X).To(Equal(uint32(1)))
			Expect(tile.X).To(Equal(uint32(1)))

			tile = binning.LonLatToTile(&lonlat2, 1)
			Expect(tile.X).To(Equal(uint32(2)))
			Expect(tile.X).To(Equal(uint32(2)))
		})
	})

})
