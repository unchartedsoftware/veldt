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

// DecodeImage takes an image file and decodes it into RGBA byte array format.
func DecodeImage(ext string, reader io.Reader) ([]byte, error) {
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

// Decode takes a io.Reader and decodes the data based on the provided
// extension.
func Decode(ext string, reader io.Reader) ([]byte, error) {
	if isImage(ext) {
		return DecodeImage(ext, reader)
	}
	// return result directly
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	if len(bytes) == 0 {
		return nil, fmt.Errorf("cannot decode empty reader")
	}
	return bytes, nil
}

func isImage(ext string) bool {
	return ext == "png" || ext == "jpg" || ext == "jpeg"
}
