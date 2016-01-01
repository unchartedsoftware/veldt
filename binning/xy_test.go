package binning_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/unchartedsoftware/prism/binning"
)

var _ = Describe("xy", func() {

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
		extent = binning.Bounds{
			TopLeft: &binning.Coord{
				X: -1,
				Y: -1,
			},
			BottomRight: &binning.Coord{
				X: 1,
				Y: 1,
			},
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

	Describe("GetTileBounds", func() {
		It("should return bounds from a tile coordinate", func() {
			bounds := binning.GetTileBounds(&t0, &extent)
			Expect(bounds.TopLeft.X).To(BeNumerically("~", -1.0, epsilon))
			Expect(bounds.TopLeft.Y).To(BeNumerically("~", -1.0, epsilon))
			Expect(bounds.BottomRight.X).To(BeNumerically("~", 1.0, epsilon))
			Expect(bounds.BottomRight.Y).To(BeNumerically("~", 1.0, epsilon))

			bounds = binning.GetTileBounds(&t1, &extent)
			Expect(bounds.TopLeft.X).To(BeNumerically("~", -1.0, epsilon))
			Expect(bounds.TopLeft.Y).To(BeNumerically("~", -1.0, epsilon))
			Expect(bounds.BottomRight.X).To(BeNumerically("~", 0.0, epsilon))
			Expect(bounds.BottomRight.Y).To(BeNumerically("~", 0.0, epsilon))

			bounds = binning.GetTileBounds(&t2, &extent)
			Expect(bounds.TopLeft.X).To(BeNumerically("~", 0.0, epsilon))
			Expect(bounds.TopLeft.Y).To(BeNumerically("~", 0.0, epsilon))
			Expect(bounds.BottomRight.X).To(BeNumerically("~", 1.0, epsilon))
			Expect(bounds.BottomRight.Y).To(BeNumerically("~", 1.0, epsilon))
		})
		It("should support bounds where right > left", func() {
			bounds := binning.GetTileBounds(&t0, &extent)
			Expect(bounds.TopLeft.X).To(BeNumerically("~", -1.0, epsilon))
			Expect(bounds.TopLeft.Y).To(BeNumerically("~", -1.0, epsilon))
			Expect(bounds.BottomRight.X).To(BeNumerically("~", 1.0, epsilon))
			Expect(bounds.BottomRight.Y).To(BeNumerically("~", 1.0, epsilon))
		})
		It("should support bounds where bottom > top", func() {
			bounds := binning.GetTileBounds(&t0, &extent)
			Expect(bounds.TopLeft.X).To(BeNumerically("~", -1.0, epsilon))
			Expect(bounds.TopLeft.Y).To(BeNumerically("~", -1.0, epsilon))
			Expect(bounds.BottomRight.X).To(BeNumerically("~", 1.0, epsilon))
			Expect(bounds.BottomRight.Y).To(BeNumerically("~", 1.0, epsilon))
		})
		It("should support bounds where left > right", func() {
			bounds := binning.GetTileBounds(&t0, &binning.Bounds{
				TopLeft: &binning.Coord{
					X: 1,
					Y: -1,
				},
				BottomRight: &binning.Coord{
					X: -1,
					Y: 1,
				},
			})
			Expect(bounds.TopLeft.X).To(BeNumerically("~", 1.0, epsilon))
			Expect(bounds.TopLeft.Y).To(BeNumerically("~", -1.0, epsilon))
			Expect(bounds.BottomRight.X).To(BeNumerically("~", -1.0, epsilon))
			Expect(bounds.BottomRight.Y).To(BeNumerically("~", 1.0, epsilon))
		})
		It("should support bounds where top > bottom", func() {
			bounds := binning.GetTileBounds(&t0, &binning.Bounds{
				TopLeft: &binning.Coord{
					X: -1,
					Y: 1,
				},
				BottomRight: &binning.Coord{
					X: 1,
					Y: -1,
				},
			})
			Expect(bounds.TopLeft.X).To(BeNumerically("~", -1.0, epsilon))
			Expect(bounds.TopLeft.Y).To(BeNumerically("~", 1.0, epsilon))
			Expect(bounds.BottomRight.X).To(BeNumerically("~", 1.0, epsilon))
			Expect(bounds.BottomRight.Y).To(BeNumerically("~", -1.0, epsilon))
		})
	})

	Describe("CoordToFractionalTile", func() {
		It("should return a fractional tile coordinate", func() {
			tile := binning.CoordToFractionalTile(&xy0, 0, &extent)
			Expect(tile.X).To(BeNumerically("~", 0.0, epsilon))
			Expect(tile.Y).To(BeNumerically("~", 0.0, epsilon))

			tile = binning.CoordToFractionalTile(&xy1, 1, &extent)
			Expect(tile.X).To(BeNumerically("~", 1.0, epsilon))
			Expect(tile.Y).To(BeNumerically("~", 1.0, epsilon))

			tile = binning.CoordToFractionalTile(&xy2, 1, &extent)
			Expect(tile.X).To(BeNumerically("~", 2.0, epsilon))
			Expect(tile.Y).To(BeNumerically("~", 2.0, epsilon))
		})
	})

	Describe("CoordToTile", func() {
		It("should return a tile coordinate", func() {
			tile := binning.CoordToTile(&xy0, 0, &extent)
			Expect(tile.X).To(Equal(uint32(0)))
			Expect(tile.X).To(Equal(uint32(0)))

			tile = binning.CoordToTile(&xy1, 1, &extent)
			Expect(tile.X).To(Equal(uint32(1)))
			Expect(tile.X).To(Equal(uint32(1)))

			tile = binning.CoordToTile(&xy2, 1, &extent)
			Expect(tile.X).To(Equal(uint32(2)))
			Expect(tile.X).To(Equal(uint32(2)))
		})
	})

})
