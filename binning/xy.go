package binning

import (
	"math"

	"github.com/unchartedsoftware/veldt/geometry"
)

// Extrema represents the min and max values for an ordinal property.
type Extrema struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// CoordToFractionalTile converts a data coordinate to a floating point tile coordinate.
func CoordToFractionalTile(coord *geometry.Coord, level uint32, bounds *geometry.Bounds) *FractionalTileCoord {
	pow2 := math.Pow(2, float64(level))
	x := pow2 * (coord.X - bounds.Left) / (bounds.Right - bounds.Left)
	y := pow2 * (coord.Y - bounds.Bottom) / (bounds.Top - bounds.Bottom)
	return &FractionalTileCoord{
		X: x,
		Y: y,
		Z: level,
	}
}

// GetTileBounds returns the data coordinate bounds of the tile coordinate.
func GetTileBounds(tile *TileCoord, bounds *geometry.Bounds) *geometry.Bounds {
	pow2 := math.Pow(2, float64(tile.Z))
	tileXSize := (bounds.Right - bounds.Left) / pow2
	tileYSize := (bounds.Top - bounds.Bottom) / pow2
	return geometry.NewBounds(
		bounds.Left+tileXSize*float64(tile.X),
		bounds.Left+tileXSize*float64(tile.X+1),
		bounds.Bottom+tileYSize*float64(tile.Y),
		bounds.Bottom+tileYSize*float64(tile.Y+1))
}
