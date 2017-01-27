package tile

import (
	"fmt"
	"image"
	"image/draw"
	"io"
	"io/ioutil"
	// register png decoder
	_ "image/png"
	// register jpeg decoder
	_ "image/jpeg"
)

// DecodeImage takes an image file and encodes it into RGBA byte array format.
func DecodeImage(ext string, reader io.Reader) ([]byte, error) {
	if isImage(ext) {
		// decode result into bytes
		img, _, err := image.Decode(reader)
		if err != nil {
			return nil, err
		}
		rgba := image.NewRGBA(img.Bounds())
		if rgba.Stride != rgba.Rect.Size().X*4 {
			return nil, fmt.Errorf("unsupported stride in requested image")
		}
		draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
		return []byte(rgba.Pix), nil
	}
	// return result directly
	return ioutil.ReadAll(reader)
}

func isImage(ext string) bool {
	return ext == "png" || ext == "jpg" || ext == "jpeg"
}
