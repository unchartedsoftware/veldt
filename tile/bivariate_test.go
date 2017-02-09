package tile_test

import (
	"github.com/unchartedsoftware/veldt/binning"
	"github.com/unchartedsoftware/veldt/tile"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/unchartedsoftware/veldt/util/test"
)

var _ = Describe("Bivariate", func() {

	var bivariate *tile.Bivariate

	BeforeEach(func() {
		bivariate = &tile.Bivariate{}
	})

	Describe("Parse", func() {
		It("should parse properties from the params argument", func() {
			params := JSON(
				`{
					"xField": "x",
					"yField": "y",
					"left": -1.0,
					"right": 1.0,
					"bottom": -1.0,
					"top": 1.0,
					"resolution": 256
				}`)
			err := bivariate.Parse(params)
			Expect(err).To(BeNil())
			Expect(bivariate.XField).To(Equal("x"))
			Expect(bivariate.YField).To(Equal("y"))
			Expect(bivariate.Resolution).To(Equal(256))
		})

		It("should return an error if `xField` property is not specified", func() {
			params := JSON(`{}`)
			err := bivariate.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `yField` property is not specified", func() {
			params := JSON(
				`{
					"xField": "x"
				}
				`)
			err := bivariate.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `left` property is not specified", func() {
			params := JSON(
				`{
					"xField": "x",
					"yField": "y"
				}
				`)
			err := bivariate.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `right` property is not specified", func() {
			params := JSON(
				`{
					"xField": "x",
					"yField": "y",
					"left": -1.0
				}
				`)
			err := bivariate.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `bottom` property is not specified", func() {
			params := JSON(
				`{
					"xField": "x",
					"yField": "y",
					"left": -1.0,
					"right": 1.0
				}
				`)
			err := bivariate.Parse(params)
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `top` property is not specified", func() {
			params := JSON(
				`{
					"xField": "x",
					"yField": "y",
					"left": -1.0,
					"right": 1.0,
					"bottom": -1.0
				}
				`)
			err := bivariate.Parse(params)
			Expect(err).NotTo(BeNil())
		})
	})

	Describe("TileBounds", func() {
		It("should return the tile bounds for the provided tile coord", func() {
			params := JSON(
				`{
					"xField": "x",
					"yField": "y",
					"left": -1.0,
					"right": 1.0,
					"bottom": -1.0,
					"top": 1.0,
					"resolution": 256
				}`)
			coord := &binning.TileCoord{
				Z: 0,
				X: 0,
				Y: 0,
			}
			err := bivariate.Parse(params)
			Expect(err).To(BeNil())
			bounds := bivariate.TileBounds(coord)
			Expect(bounds.Left).To(Equal(-1.0))
			Expect(bounds.Right).To(Equal(1.0))
			Expect(bounds.Bottom).To(Equal(-1.0))
			Expect(bounds.Top).To(Equal(1.0))
		})
	})

	Describe("BinSizeX", func() {
		It("should return the size of a bin over the x axis", func() {
			params := JSON(
				`{
					"xField": "x",
					"yField": "y",
					"left": -1.0,
					"right": 1.0,
					"bottom": -1.0,
					"top": 1.0,
					"resolution": 256
				}`)
			coord := &binning.TileCoord{
				Z: 0,
				X: 0,
				Y: 0,
			}
			err := bivariate.Parse(params)
			Expect(err).To(BeNil())
			size := bivariate.BinSizeX(coord)
			Expect(size).To(Equal(2.0 / 256.0))
		})
	})

	Describe("BinSizeY", func() {
		It("should return the tile bounds for the provided tile coord", func() {
			params := JSON(
				`{
					"xField": "x",
					"yField": "y",
					"left": -1.0,
					"right": 1.0,
					"bottom": -1.0,
					"top": 1.0,
					"resolution": 256
				}`)
			coord := &binning.TileCoord{
				Z: 0,
				X: 0,
				Y: 0,
			}
			err := bivariate.Parse(params)
			Expect(err).To(BeNil())
			size := bivariate.BinSizeY(coord)
			Expect(size).To(Equal(2.0 / 256.0))
		})
	})

	Describe("GetXBin", func() {
		It("should return the x bin for the provided tile coord for left < right", func() {
			params := JSON(
				`{
					"xField": "x",
					"yField": "y",
					"left": -1.0,
					"right": 1.0,
					"bottom": -1.0,
					"top": 1.0,
					"resolution": 256
				}`)
			coord := &binning.TileCoord{
				Z: 0,
				X: 0,
				Y: 0,
			}
			err := bivariate.Parse(params)
			Expect(err).To(BeNil())
			binA := bivariate.GetXBin(coord, -1.0)
			binB := bivariate.GetXBin(coord, -0.5)
			binC := bivariate.GetXBin(coord, 0.0)
			binD := bivariate.GetXBin(coord, 0.5)
			binE := bivariate.GetXBin(coord, 1.0)
			Expect(binA).To(Equal(0))
			Expect(binB).To(Equal(64))
			Expect(binC).To(Equal(128))
			Expect(binD).To(Equal(192))
			Expect(binE).To(Equal(255))
		})

		It("should return the x bin for the provided tile coord for left > right", func() {
			params := JSON(
				`{
					"xField": "x",
					"yField": "y",
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
			err := bivariate.Parse(params)
			Expect(err).To(BeNil())
			binA := bivariate.GetXBin(coord, 1.0)
			binB := bivariate.GetXBin(coord, 0.5)
			binC := bivariate.GetXBin(coord, 0.0)
			binD := bivariate.GetXBin(coord, -0.5)
			binE := bivariate.GetXBin(coord, -1.0)
			Expect(binA).To(Equal(0))
			Expect(binB).To(Equal(63))
			Expect(binC).To(Equal(127))
			Expect(binD).To(Equal(191))
			Expect(binE).To(Equal(255))
		})
	})

	Describe("GetYBin", func() {
		It("should return the y bin for the provided tile coord for bottom < top", func() {
			params := JSON(
				`{
					"xField": "x",
					"yField": "y",
					"left": -1.0,
					"right": 1.0,
					"bottom": -1.0,
					"top": 1.0,
					"resolution": 256
				}`)
			coord := &binning.TileCoord{
				Z: 0,
				X: 0,
				Y: 0,
			}
			err := bivariate.Parse(params)
			Expect(err).To(BeNil())
			binA := bivariate.GetYBin(coord, -1.0)
			binB := bivariate.GetYBin(coord, -0.5)
			binC := bivariate.GetYBin(coord, 0.0)
			binD := bivariate.GetYBin(coord, 0.5)
			binE := bivariate.GetYBin(coord, 1.0)
			Expect(binA).To(Equal(0))
			Expect(binB).To(Equal(64))
			Expect(binC).To(Equal(128))
			Expect(binD).To(Equal(192))
			Expect(binE).To(Equal(255))
		})

		It("should return the y bin for the provided tile coord for bottom > top", func() {
			params := JSON(
				`{
					"xField": "x",
					"yField": "y",
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
			err := bivariate.Parse(params)
			Expect(err).To(BeNil())
			binA := bivariate.GetYBin(coord, 1.0)
			binB := bivariate.GetYBin(coord, 0.5)
			binC := bivariate.GetYBin(coord, 0.0)
			binD := bivariate.GetYBin(coord, -0.5)
			binE := bivariate.GetYBin(coord, -1.0)
			Expect(binA).To(Equal(0))
			Expect(binB).To(Equal(63))
			Expect(binC).To(Equal(127))
			Expect(binD).To(Equal(191))
			Expect(binE).To(Equal(255))
		})
	})

	Describe("GetX", func() {
		It("should return the x coordinate for the provided tile coord for left < right", func() {
			params := JSON(
				`{
					"xField": "x",
					"yField": "y",
					"left": -1.0,
					"right": 1.0,
					"bottom": -1.0,
					"top": 1.0,
					"resolution": 256
				}`)
			coord := &binning.TileCoord{
				Z: 0,
				X: 0,
				Y: 0,
			}
			err := bivariate.Parse(params)
			Expect(err).To(BeNil())
			binA := bivariate.GetX(coord, -1.0)
			binB := bivariate.GetX(coord, -0.5)
			binC := bivariate.GetX(coord, 0.0)
			binD := bivariate.GetX(coord, 0.5)
			binE := bivariate.GetX(coord, 1.0)
			Expect(binA).To(Equal(0.0))
			Expect(binB).To(Equal(64.0))
			Expect(binC).To(Equal(128.0))
			Expect(binD).To(Equal(192.0))
			Expect(binE).To(Equal(256.0))
		})

		It("should return the x coordinate for the provided tile coord for left > right", func() {
			params := JSON(
				`{
					"xField": "x",
					"yField": "y",
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
			err := bivariate.Parse(params)
			Expect(err).To(BeNil())
			binA := bivariate.GetX(coord, 1.0)
			binB := bivariate.GetX(coord, 0.5)
			binC := bivariate.GetX(coord, 0.0)
			binD := bivariate.GetX(coord, -0.5)
			binE := bivariate.GetX(coord, -1.0)
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
					"xField": "x",
					"yField": "y",
					"left": -1.0,
					"right": 1.0,
					"bottom": -1.0,
					"top": 1.0,
					"resolution": 256
				}`)
			coord := &binning.TileCoord{
				Z: 0,
				X: 0,
				Y: 0,
			}
			err := bivariate.Parse(params)
			Expect(err).To(BeNil())
			binA := bivariate.GetY(coord, -1.0)
			binB := bivariate.GetY(coord, -0.5)
			binC := bivariate.GetY(coord, 0.0)
			binD := bivariate.GetY(coord, 0.5)
			binE := bivariate.GetY(coord, 1.0)
			Expect(binA).To(Equal(0.0))
			Expect(binB).To(Equal(64.0))
			Expect(binC).To(Equal(128.0))
			Expect(binD).To(Equal(192.0))
			Expect(binE).To(Equal(256.0))
		})

		It("should return the y coordinate for the provided tile coord for bottom > top", func() {
			params := JSON(
				`{
					"xField": "x",
					"yField": "y",
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
			err := bivariate.Parse(params)
			Expect(err).To(BeNil())
			binA := bivariate.GetY(coord, 1.0)
			binB := bivariate.GetY(coord, 0.5)
			binC := bivariate.GetY(coord, 0.0)
			binD := bivariate.GetY(coord, -0.5)
			binE := bivariate.GetY(coord, -1.0)
			Expect(binA).To(Equal(0.0))
			Expect(binB).To(Equal(64.0))
			Expect(binC).To(Equal(128.0))
			Expect(binD).To(Equal(192.0))
			Expect(binE).To(Equal(256.0))
		})
	})

	Describe("GetXY", func() {
		It("should return the x and y coordinate from the hit map for the provided tile coord", func() {
			params := JSON(
				`{
					"xField": "x",
					"yField": "y",
					"left": -1.0,
					"right": 1.0,
					"bottom": -1.0,
					"top": 1.0,
					"resolution": 256
				}`)
			coord := &binning.TileCoord{
				Z: 0,
				X: 0,
				Y: 0,
			}
			hitA := JSON(`{ "x": -1.0, "y": -1.0 }`)
			hitB := JSON(`{ "x": 0.0, "y": 0.0 }`)
			hitC := make(map[string]interface{})
			hitC["x"] = int64(1)
			hitC["y"] = int64(1)
			err := bivariate.Parse(params)
			Expect(err).To(BeNil())
			xA, yA, ok := bivariate.GetXY(coord, hitA)
			Expect(ok).To(Equal(true))
			Expect(xA).To(Equal(0.0))
			Expect(yA).To(Equal(0.0))
			xB, yB, ok := bivariate.GetXY(coord, hitB)
			Expect(ok).To(Equal(true))
			Expect(xB).To(Equal(128.0))
			Expect(yB).To(Equal(128.0))
			xC, yC, ok := bivariate.GetXY(coord, hitC)
			Expect(ok).To(Equal(true))
			Expect(xC).To(Equal(256.0))
			Expect(yC).To(Equal(256.0))
		})
		It("should return false if the `xField` does not exist in the hit", func() {
			params := JSON(
				`{
					"xField": "x",
					"yField": "y",
					"left": -1.0,
					"right": 1.0,
					"bottom": -1.0,
					"top": 1.0,
					"resolution": 256
				}`)
			coord := &binning.TileCoord{
				Z: 0,
				X: 0,
				Y: 0,
			}
			hit := JSON(`{ "y": -1.0 }`)
			err := bivariate.Parse(params)
			Expect(err).To(BeNil())
			_, _, ok := bivariate.GetXY(coord, hit)
			Expect(ok).To(Equal(false))
		})
		It("should return false if the `yField` does not exist in the hit", func() {
			params := JSON(
				`{
					"xField": "x",
					"yField": "y",
					"left": -1.0,
					"right": 1.0,
					"bottom": -1.0,
					"top": 1.0,
					"resolution": 256
				}`)
			coord := &binning.TileCoord{
				Z: 0,
				X: 0,
				Y: 0,
			}
			hit := JSON(`{ "x": -1.0 }`)
			err := bivariate.Parse(params)
			Expect(err).To(BeNil())
			_, _, ok := bivariate.GetXY(coord, hit)
			Expect(ok).To(Equal(false))
		})
		It("should return false if the x or y value is not numeric", func() {
			params := JSON(
				`{
					"xField": "x",
					"yField": "y",
					"left": -1.0,
					"right": 1.0,
					"bottom": -1.0,
					"top": 1.0,
					"resolution": 256
				}`)
			coord := &binning.TileCoord{
				Z: 0,
				X: 0,
				Y: 0,
			}
			hit := JSON(`{ "x": "string", "y": "string" }`)
			err := bivariate.Parse(params)
			Expect(err).To(BeNil())
			_, _, ok := bivariate.GetXY(coord, hit)
			Expect(ok).To(Equal(false))
		})
	})
})
