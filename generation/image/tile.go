package image

import (
	"fmt"
	"image"
	"image/draw"
	"io/ioutil"
	// register png decoder
	_ "image/png"
	// register jpeg decoder
	_ "image/jpeg"
	"net/http"

	"github.com/unchartedsoftware/prism"
	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/util/json"
)

type Tile struct {
	ext      string
	endpoint string
	scheme   string
}

func NewTile() prism.TileCtor {
	return func() (prism.Tile, error) {
		return &Tile{}, nil
	}
}

func (t *Tile) Parse(params map[string]interface{}) error {
	// get endpoint
	endpoint, ok := json.GetString(params, "endpoint")
	if !ok {
		return fmt.Errorf("`endpoint` parameter missing from tile")
	}
	// get scheme
	scheme, ok := json.GetString(params, "scheme")
	if !ok {
		return fmt.Errorf("`scheme` parameter missing from tile")
	}
	// get ext
	ext, ok := json.GetString(params, "ext")
	if !ok {
		return fmt.Errorf("`ext` parameter missing from tile")
	}
	t.ext = ext
	t.endpoint = endpoint
	t.scheme = scheme
	return nil
}

func (t *Tile) Create(uri string, coord *binning.TileCoord, query prism.Query) ([]byte, error) {
	// create URL
	format := "%s://%s/%s/%d/%d/%d.%s"
	url := fmt.Sprintf(format,
		t.scheme,
		t.endpoint,
		uri,
		coord.Z,
		coord.X,
		coord.Y,
		t.ext)
	// build http request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// set appropriate headers based on extention
	handleExt(t.ext, req)
	// build http request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	// check status code
	if res.StatusCode >= 400 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf(string(body))
	}

	if isImage(t.ext) {
		// decode result into bytes
		img, _, err := image.Decode(res.Body)
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
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func isImage(ext string) bool {
	return ext == "png" || ext == "jpg" || ext == "jpeg"
}

func handleExt(ext string, req *http.Request) {
	switch ext {
	case "png":
		req.Header.Set("Accept", "image/png")
	case "jpg":
		req.Header.Set("Accept", "image/jpg")
	case "jpeg":
		req.Header.Set("Accept", "image/jpeg")
	case "json":
		req.Header.Set("Accept", "application/json")
	case "bin":
		req.Header.Set("Accept", "application/octet-stream")
	}
}
