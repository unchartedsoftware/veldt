package file

import (
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
	"github.com/unchartedsoftware/veldt/tile"
	"github.com/unchartedsoftware/veldt/util/json"
)

// Tile represents a filesystem tile type.
type Tile struct {
	path      string
	ext       string
	padCoords bool
}

// NewTile instantiates and returns a new filesystem tile.
func NewTile() veldt.TileCtor {
	return func() (veldt.Tile, error) {
		return &Tile{}, nil
	}
}

// Parse parses the provided JSON object and populates the tiles attributes.
func (t *Tile) Parse(params map[string]interface{}) error {
	// get path
	path, ok := json.GetString(params, "path")
	if !ok {
		return fmt.Errorf("`path` parameter missing from tile")
	}
	// get ext
	ext, ok := json.GetString(params, "ext")
	if !ok {
		return fmt.Errorf("`ext` parameter missing from tile")
	}
	// do we pad the coords?
	padcoords := json.GetBoolDefault(params, true, "padcoords")
	// set attributes
	t.path = path
	t.ext = ext
	t.padCoords = padcoords
	return nil
}

// Create generates a tile from the provided URI, tile coordinate and query
// parameters.
func (t *Tile) Create(uri string, coord *binning.TileCoord, query veldt.Query) ([]byte, error) {
	// get the format string for the path
	format := t.getFormat(coord)
	// get the filepath
	filepath := fmt.Sprintf(format,
		t.path,
		coord.Z,
		coord.X,
		coord.Y,
		t.ext,
	)
	// open the file
	file, err := os.Open(filepath)
	if err != nil {
		// don't return an error if the tile doesn't exist
		if err == os.ErrNotExist {
			return []byte{}, nil
		}
		return nil, err
	}
	defer file.Close()
	// decode file
	return tile.DecodeImage(t.ext, file)
}

func (t *Tile) getFormat(coord *binning.TileCoord) string {
	if t.padCoords {
		digits := strconv.Itoa(int(math.Floor(math.Log10(float64(int(1)<<coord.Z)))) + 1)
		return "%s/%02d/%0" + digits + "d/%0" + digits + "d.%s"
	}
	return "%s/%d/%d/%d.%s"
}
