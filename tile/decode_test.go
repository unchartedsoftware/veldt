package tile_test

import (
	"bytes"
	"image"
	"image/color"
	"image/png"

	"github.com/unchartedsoftware/veldt/tile"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Decode", func() {

	var json string
	var img *image.RGBA

	BeforeEach(func() {
		// json string
		json = `{
			"a": "test",
			"b": 0.5
		}`
		// image bytes
		img = image.NewRGBA(image.Rectangle{
			image.Point{0, 0},
			image.Point{255, 255},
		})
		for x := 0; x < 16; x++ {
			for y := 0; y < 16; y++ {
				c := color.RGBA{
					uint8(x * 16),
					uint8(y * 16),
					uint8(x + y%256),
					255,
				}
				img.Set(x, y, c)
			}
		}
	})

	Describe("Decode", func() {
		It("should return the raw bytes of non-image extensions", func() {
			buffer := bytes.NewBufferString(json)
			bs, err := tile.Decode("json", buffer)
			Expect(err).To(BeNil())
			Expect(bs).To(Equal([]byte(json)))
		})
		It("should decode the raw RGBA bytes of a png / jpg / jpeg image", func() {
			buffer := &bytes.Buffer{}
			png.Encode(buffer, img)
			bs, err := tile.Decode("png", buffer)
			Expect(err).To(BeNil())
			Expect(bs).To(Equal([]byte(img.Pix)))
		})
		It("should return an error if there is no data", func() {
			buffer := &bytes.Buffer{}
			_, err := tile.Decode("json", buffer)
			Expect(err).NotTo(BeNil())
		})
	})

	Describe("DecodeImage", func() {
		It("should decode the raw RGBA bytes of a png / jpg / jpeg image", func() {
			buffer := &bytes.Buffer{}
			png.Encode(buffer, img)
			bs, err := tile.DecodeImage("png", buffer)
			Expect(err).To(BeNil())
			Expect(bs).To(Equal([]byte(img.Pix)))
		})
		It("should return an error if the image buffer is empty", func() {
			buffer := &bytes.Buffer{}
			_, err := tile.Decode("png", buffer)
			Expect(err).NotTo(BeNil())
		})
	})
})
