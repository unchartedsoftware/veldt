package binning

import (
	"math"
)

// Bounds represents a bounding box.
type Bounds struct {
	BottomLeft *Coord
	TopRight   *Coord
}

// Coord represents a point.
type Coord struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Extrema represents the min and max values for an ordinal property.
type Extrema struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// NewCoord instantiates and returns a pointer to a Coord.
func NewCoord(x, y float64) *Coord {
	return &Coord{
		X: x,
		Y: y,
	}
}

// CoordToFractionalTile converts a data coordinate to a floating point tile coordinate.
func CoordToFractionalTile(coord *Coord, level uint32, bounds *Bounds) *FractionalTileCoord {
	pow2 := math.Pow(2, float64(level))
	x := pow2 * (coord.X - bounds.BottomLeft.X) / (bounds.TopRight.X - bounds.BottomLeft.X)
	y := pow2 * (coord.Y - bounds.BottomLeft.Y) / (bounds.TopRight.Y - bounds.BottomLeft.Y)
	return &FractionalTileCoord{
		X: x,
		Y: y,
		Z: level,
	}
}

// GetTileBounds returns the data coordinate bounds of the tile coordinate.
func GetTileBounds(tile *TileCoord, bounds *Bounds) *Bounds {
	pow2 := math.Pow(2, float64(tile.Z))
	tileXSize := (bounds.TopRight.X - bounds.BottomLeft.X) / pow2
	tileYSize := (bounds.TopRight.Y - bounds.BottomLeft.Y) / pow2
	return &Bounds{
		BottomLeft: &Coord{
			X: bounds.BottomLeft.X + tileXSize*float64(tile.X),
			Y: bounds.BottomLeft.Y + tileYSize*float64(tile.Y),
		},
		TopRight: &Coord{
			X: bounds.BottomLeft.X + tileXSize*float64(tile.X+1),
			Y: bounds.BottomLeft.Y + tileYSize*float64(tile.Y+1),
		},
	}
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
