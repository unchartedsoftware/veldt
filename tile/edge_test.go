package tile_test

import (
	"github.com/unchartedsoftware/veldt/binning"
	"github.com/unchartedsoftware/veldt/tile"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/unchartedsoftware/veldt/util/test"
)

var _ = Describe("Edge", func() {

	var edge *tile.Edge

	BeforeEach(func() {
		edge = &tile.Edge{}
	})

	Describe("Parse", func() {
		It("should parse properties from the params argument", func() {
			params := JSON(
				`{
					"srcXField": "sx",
					"srcYField": "sy",
					"dstXField": "dx",
					"dstYField": "dy",
					"left": -1.0,
					"right": 1.0,
					"bottom": -1.0,
					"top": 1.0
				}`)
			err := edge.Parse(params)
			Expect(err).To(BeNil())
			Expect(edge.SrcXField).To(Equal("sx"))
			Expect(edge.SrcYField).To(Equal("sy"))
			Expect(edge.DstXField).To(Equal("dx"))
			Expect(edge.DstYField).To(Equal("dy"))
		})

		It("should return an error if `srcXField` property is not specified", func() {
			params := JSON(`{}`)
			err := edge.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `srcYField` property is not specified", func() {
			params := JSON(
				`{
					"srcXField": "sx"
				}
				`)
			err := edge.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `dstXField` property is not specified", func() {
			params := JSON(
				`{
					"srcXField": "sx",
					"srcYField": "sy"
				}
				`)
			err := edge.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `dstYField` property is not specified", func() {
			params := JSON(
				`{
					"srcXField": "sx",
					"srcYField": "sy",
					"dstXField": "dx"
				}
				`)
			err := edge.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `left` property is not specified", func() {
			params := JSON(
				`{
					"srcXField": "sx",
					"srcYField": "sy",
					"dstXField": "dx",
					"dstYField": "dy"
				}
				`)
			err := edge.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `right` property is not specified", func() {
			params := JSON(
				`{
					"srcXField": "sx",
					"srcYField": "sy",
					"dstXField": "dx",
					"dstYField": "dy",
					"left": -1.0
				}
				`)
			err := edge.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `bottom` property is not specified", func() {
			params := JSON(
				`{
					"srcXField": "sx",
					"srcYField": "sy",
					"dstXField": "dx",
					"dstYField": "dy",
					"left": -1.0,
					"right": 1.0
				}
				`)
			err := edge.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `top` property is not specified", func() {
			params := JSON(
				`{
					"srcXField": "sx",
					"srcYField": "sy",
					"dstXField": "dx",
					"dstYField": "dy",
					"left": -1.0,
					"right": 1.0,
					"bottom": -1.0
				}
				`)
			err := edge.Parse(params)
			Expect(err).NotTo(BeNil())
		})
	})

	Describe("TileBounds", func() {
		It("should return the tile bounds for the provided tile coord", func() {
			params := JSON(
				`{
					"srcXField": "sx",
					"srcYField": "sy",
					"dstXField": "dx",
					"dstYField": "dy",
					"left": -1.0,
					"right": 1.0,
					"bottom": -1.0,
					"top": 1.0
				}`)
			coord := &binning.TileCoord{
				Z: 0,
				X: 0,
				Y: 0,
			}
			err := edge.Parse(params)
			Expect(err).To(BeNil())
			bounds := edge.TileBounds(coord)
			Expect(bounds.Left).To(Equal(-1.0))
			Expect(bounds.Right).To(Equal(1.0))
			Expect(bounds.Bottom).To(Equal(-1.0))
			Expect(bounds.Top).To(Equal(1.0))
		})
	})

	Describe("GetX", func() {
		It("should return the x coordinate for the provided tile coord for left < right", func() {
			params := JSON(
				`{
					"srcXField": "sx",
					"srcYField": "sy",
					"dstXField": "dx",
					"dstYField": "dy",
					"left": -1.0,
					"right": 1.0,
					"bottom": -1.0,
					"top": 1.0
				}`)
			coord := &binning.TileCoord{
				Z: 0,
				X: 0,
				Y: 0,
			}
			err := edge.Parse(params)
			Expect(err).To(BeNil())
			binA := edge.GetX(coord, -1.0)
			binB := edge.GetX(coord, -0.5)
			binC := edge.GetX(coord, 0.0)
			binD := edge.GetX(coord, 0.5)
			binE := edge.GetX(coord, 1.0)
			Expect(binA).To(Equal(0.0))
			Expect(binB).To(Equal(64.0))
			Expect(binC).To(Equal(128.0))
			Expect(binD).To(Equal(192.0))
			Expect(binE).To(Equal(256.0))
		})

		It("should return the x coordinate for the provided tile coord for left > right", func() {
			params := JSON(
				`{
					"srcXField": "sx",
					"srcYField": "sy",
					"dstXField": "dx",
					"dstYField": "dy",
					"left": 1.0,
					"right": -1.0,
					"bottom": 1.0,
					"top": -1.0,
					"resolution": 256
				}`)
			coord := &binning.TileCoord{
				Z: 0,
				X: 0,
				Y: 0,
			}
			err := edge.Parse(params)
			Expect(err).To(BeNil())
			binA := edge.GetX(coord, 1.0)
			binB := edge.GetX(coord, 0.5)
			binC := edge.GetX(coord, 0.0)
			binD := edge.GetX(coord, -0.5)
			binE := edge.GetX(coord, -1.0)
			Expect(binA).To(Equal(0.0))
			Expect(binB).To(Equal(64.0))
			Expect(binC).To(Equal(128.0))
			Expect(binD).To(Equal(192.0))
			Expect(binE).To(Equal(256.0))
		})
	})

	Describe("GetY", func() {
		It("should return the y coordinate for the provided tile coord for bottom < top", func() {
			params := JSON(
				`{
					"srcXField": "sx",
					"srcYField": "sy",
					"dstXField": "dx",
					"dstYField": "dy",
					"left": -1.0,
					"right": 1.0,
					"bottom": -1.0,
					"top": 1.0
				}`)
			coord := &binning.TileCoord{
				Z: 0,
				X: 0,
				Y: 0,
			}
			err := edge.Parse(params)
			Expect(err).To(BeNil())
			binA := edge.GetY(coord, -1.0)
			binB := edge.GetY(coord, -0.5)
			binC := edge.GetY(coord, 0.0)
			binD := edge.GetY(coord, 0.5)
			binE := edge.GetY(coord, 1.0)
			Expect(binA).To(Equal(0.0))
			Expect(binB).To(Equal(64.0))
			Expect(binC).To(Equal(128.0))
			Expect(binD).To(Equal(192.0))
			Expect(binE).To(Equal(256.0))
		})

		It("should return the y coordinate for the provided tile coord for bottom > top", func() {
			params := JSON(
				`{
					"srcXField": "sx",
					"srcYField": "sy",
					"dstXField": "dx",
					"dstYField": "dy",
					"left": 1.0,
					"right": -1.0,
					"bottom": 1.0,
					"top": -1.0,
					"resolution": 256
				}`)
			coord := &binning.TileCoord{
				Z: 0,
				X: 0,
				Y: 0,
			}
			err := edge.Parse(params)
			Expect(err).To(BeNil())
			binA := edge.GetY(coord, 1.0)
			binB := edge.GetY(coord, 0.5)
			binC := edge.GetY(coord, 0.0)
			binD := edge.GetY(coord, -0.5)
			binE := edge.GetY(coord, -1.0)
			Expect(binA).To(Equal(0.0))
			Expect(binB).To(Equal(64.0))
			Expect(binC).To(Equal(128.0))
			Expect(binD).To(Equal(192.0))
			Expect(binE).To(Equal(256.0))
		})
	})

	Describe("GetSrcXY", func() {
		It("should return the x and y coordinate from the hit map for the provided tile coord", func() {
			params := JSON(
				`{
					"srcXField": "sx",
					"srcYField": "sy",
					"dstXField": "dx",
					"dstYField": "dy",
					"left": -1.0,
					"right": 1.0,
					"bottom": -1.0,
					"top": 1.0
				}`)
			coord := &binning.TileCoord{
				Z: 0,
				X: 0,
				Y: 0,
			}
			hitA := JSON(`{ "sx": -1.0, "sy": -1.0 }`)
			hitB := JSON(`{ "sx": 0.0, "sy": 0.0 }`)
			hitC := JSON(`{ "sx": 1.0, "sy": 1.0 }`)
			err := edge.Parse(params)
			Expect(err).To(BeNil())
			xA, yA, ok := edge.GetSrcXY(coord, hitA)
			Expect(ok).To(Equal(true))
			Expect(xA).To(Equal(0.0))
			Expect(yA).To(Equal(0.0))
			xB, yB, ok := edge.GetSrcXY(coord, hitB)
			Expect(ok).To(Equal(true))
			Expect(xB).To(Equal(128.0))
			Expect(yB).To(Equal(128.0))
			xC, yC, ok := edge.GetSrcXY(coord, hitC)
			Expect(ok).To(Equal(true))
			Expect(xC).To(Equal(256.0))
			Expect(yC).To(Equal(256.0))
		})
		It("should return false if the `srcXField` does not exist in the hit", func() {
			params := JSON(
				`{
					"srcXField": "sx",
					"srcYField": "sy",
					"dstXField": "dx",
					"dstYField": "dy",
					"left": -1.0,
					"right": 1.0,
					"bottom": -1.0,
					"top": 1.0
				}`)
			coord := &binning.TileCoord{
				Z: 0,
				X: 0,
				Y: 0,
			}
			hit := JSON(`{ "sy": -1.0 }`)
			err := edge.Parse(params)
			Expect(err).To(BeNil())
			_, _, ok := edge.GetSrcXY(coord, hit)
			Expect(ok).To(Equal(false))
		})
		It("should return false if the `srcYField` does not exist in the hit", func() {
			params := JSON(
				`{
					"srcXField": "sx",
					"srcYField": "sy",
					"dstXField": "dx",
					"dstYField": "dy",
					"left": -1.0,
					"right": 1.0,
					"bottom": -1.0,
					"top": 1.0
				}`)
			coord := &binning.TileCoord{
				Z: 0,
				X: 0,
				Y: 0,
			}
			hit := JSON(`{ "sx": -1.0 }`)
			err := edge.Parse(params)
			Expect(err).To(BeNil())
			_, _, ok := edge.GetSrcXY(coord, hit)
			Expect(ok).To(Equal(false))
		})
		It("should return false if the x or y value is not numeric", func() {
			params := JSON(
				`{
					"srcXField": "sx",
					"srcYField": "sy",
					"dstXField": "dx",
					"dstYField": "dy",
					"left": -1.0,
					"right": 1.0,
					"bottom": -1.0,
					"top": 1.0
				}`)
			coord := &binning.TileCoord{
				Z: 0,
				X: 0,
				Y: 0,
			}
			hit := JSON(`{ "sx": "string", "sy": "string" }`)
			err := edge.Parse(params)
			Expect(err).To(BeNil())
			_, _, ok := edge.GetSrcXY(coord, hit)
			Expect(ok).To(Equal(false))
		})
	})

	Describe("GetDstXY", func() {
		It("should return the x and y coordinate from the hit map for the provided tile coord", func() {
			params := JSON(
				`{
					"srcXField": "sx",
					"srcYField": "sy",
					"dstXField": "dx",
					"dstYField": "dy",
					"left": -1.0,
					"right": 1.0,
					"bottom": -1.0,
					"top": 1.0
				}`)
			coord := &binning.TileCoord{
				Z: 0,
				X: 0,
				Y: 0,
			}
			hitA := JSON(`{ "dx": -1.0, "dy": -1.0 }`)
			hitB := JSON(`{ "dx": 0.0, "dy": 0.0 }`)
			hitC := JSON(`{ "dx": 1.0, "dy": 1.0 }`)
			err := edge.Parse(params)
			Expect(err).To(BeNil())
			xA, yA, ok := edge.GetDstXY(coord, hitA)
			Expect(ok).To(Equal(true))
			Expect(xA).To(Equal(0.0))
			Expect(yA).To(Equal(0.0))
			xB, yB, ok := edge.GetDstXY(coord, hitB)
			Expect(ok).To(Equal(true))
			Expect(xB).To(Equal(128.0))
			Expect(yB).To(Equal(128.0))
			xC, yC, ok := edge.GetDstXY(coord, hitC)
			Expect(ok).To(Equal(true))
			Expect(xC).To(Equal(256.0))
			Expect(yC).To(Equal(256.0))
		})
		It("should return false if the `srcXField` does not exist in the hit", func() {
			params := JSON(
				`{
					"srcXField": "sx",
					"srcYField": "sy",
					"dstXField": "dx",
					"dstYField": "dy",
					"left": -1.0,
					"right": 1.0,
					"bottom": -1.0,
					"top": 1.0
				}`)
			coord := &binning.TileCoord{
				Z: 0,
				X: 0,
				Y: 0,
			}
			hit := JSON(`{ "dy": -1.0 }`)
			err := edge.Parse(params)
			Expect(err).To(BeNil())
			_, _, ok := edge.GetDstXY(coord, hit)
			Expect(ok).To(Equal(false))
		})
		It("should return false if the `srcYField` does not exist in the hit", func() {
			params := JSON(
				`{
					"srcXField": "sx",
					"srcYField": "sy",
					"dstXField": "dx",
					"dstYField": "dy",
					"left": -1.0,
					"right": 1.0,
					"bottom": -1.0,
					"top": 1.0
				}`)
			coord := &binning.TileCoord{
				Z: 0,
				X: 0,
				Y: 0,
			}
			hit := JSON(`{ "dx": -1.0 }`)
			err := edge.Parse(params)
			Expect(err).To(BeNil())
			_, _, ok := edge.GetDstXY(coord, hit)
			Expect(ok).To(Equal(false))
		})
		It("should return false if the x or y value is not numeric", func() {
			params := JSON(
				`{
					"srcXField": "sx",
					"srcYField": "sy",
					"dstXField": "dx",
					"dstYField": "dy",
					"left": -1.0,
					"right": 1.0,
					"bottom": -1.0,
					"top": 1.0
				}`)
			coord := &binning.TileCoord{
				Z: 0,
				X: 0,
				Y: 0,
			}
			hit := JSON(`{ "dx": "string", "dy": "string" }`)
			err := edge.Parse(params)
			Expect(err).To(BeNil())
			_, _, ok := edge.GetDstXY(coord, hit)
			Expect(ok).To(Equal(false))
		})
	})
})
