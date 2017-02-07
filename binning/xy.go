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
	x := pow2 * (coord.X - bounds.BottomLeft().X) / (bounds.TopRight().X - bounds.BottomLeft().X)
	y := pow2 * (coord.Y - bounds.BottomLeft().Y) / (bounds.TopRight().Y - bounds.BottomLeft().Y)
	return &FractionalTileCoord{
		X: x,
		Y: y,
		Z: level,
	}
}

// GetTileBounds returns the data coordinate bounds of the tile coordinate.
func GetTileBounds(tile *TileCoord, bounds *geometry.Bounds) *geometry.Bounds {
	pow2 := math.Pow(2, float64(tile.Z))
	corners := bounds.Corners()
	tileXSize := (corners.TopRight.X - corners.BottomLeft.X) / pow2
	tileYSize := (corners.TopRight.Y - corners.BottomLeft.Y) / pow2
	return geometry.NewBoundsFromRectangle(
		&geometry.Rectangle{
			BottomLeft: geometry.NewCoord(
				bounds.BottomLeft().X+tileXSize*float64(tile.X),
				bounds.BottomLeft().Y+tileYSize*float64(tile.Y),
			),
			TopRight: geometry.NewCoord(
				bounds.BottomLeft().X+tileXSize*float64(tile.X+1),
				bounds.BottomLeft().Y+tileYSize*float64(tile.Y+1),
			),
		},
	)
}

// GetTileExtrema returns the data coordinate bounds of the tile coordinate.
func GetTileExtrema(coord uint32, level uint32, extrema *Extrema) *Extrema {
	pow2 := math.Pow(2, float64(level))
	interval := (extrema.Max - extrema.Min) / pow2
	return &Extrema{
		Min: extrema.Min + interval*float64(coord),
		Max: extrema.Min + interval*float64(coord+1),
	}
}
