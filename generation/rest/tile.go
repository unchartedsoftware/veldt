package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
	"github.com/unchartedsoftware/veldt/tile"
	"github.com/unchartedsoftware/veldt/util/json"
)

// Tile represents a REST tile type.
type Tile struct {
	ext      string
	endpoint string
	scheme   string
}

// NewTile instantiates and returns a new REST tile.
func NewTile() veldt.TileCtor {
	return func() (veldt.Tile, error) {
		return &Tile{}, nil
	}
}

// Parse parses the provided JSON object and populates the tiles attributes.
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
	// do we pad the coords?
	t.ext = ext
	t.endpoint = endpoint
	t.scheme = scheme
	return nil
}

// Create generates a tile from the provided URI, tile coordinate and query
// parameters.
func (t *Tile) Create(uri string, coord *binning.TileCoord, query veldt.Query) ([]byte, error) {
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
	// set appropriate headers based on extension
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
	return tile.DecodeImage(t.ext, res.Body)
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
