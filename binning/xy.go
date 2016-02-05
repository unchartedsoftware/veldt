package binning

import (
	"math"
)

// Bounds represents a bounding box.
type Bounds struct {
	TopLeft     *Coord
	BottomRight *Coord
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
	x := pow2 * (coord.X - bounds.TopLeft.X) / (bounds.BottomRight.X - bounds.TopLeft.X)
	y := pow2 * (coord.Y - bounds.TopLeft.Y) / (bounds.BottomRight.Y - bounds.TopLeft.Y)
	return &FractionalTileCoord{
		X: x,
		Y: y,
		Z: level,
	}
}

// CoordToTile converts a data coordinate to a tile coordinate.
func CoordToTile(coord *Coord, level uint32, bounds *Bounds) *TileCoord {
	tile := CoordToFractionalTile(coord, level, bounds)
	return &TileCoord{
		X: uint32(math.Floor(tile.X)),
		Y: uint32(math.Floor(tile.Y)),
		Z: level,
	}
}

// GetTileBounds returns the data coordinate bounds of the tile coordinate.
func GetTileBounds(tile *TileCoord, bounds *Bounds) *Bounds {
	pow2 := math.Pow(2, float64(tile.Z))
	tileXSize := (bounds.BottomRight.X - bounds.TopLeft.X) / pow2
	tileYSize := (bounds.BottomRight.Y - bounds.TopLeft.Y) / pow2
	return &Bounds{
		TopLeft: &Coord{
			X: bounds.TopLeft.X + tileXSize*float64(tile.X),
			Y: bounds.TopLeft.Y + tileYSize*float64(tile.Y),
		},
		BottomRight: &Coord{
			X: bounds.TopLeft.X + tileXSize*float64(tile.X+1),
			Y: bounds.TopLeft.Y + tileYSize*float64(tile.Y+1),
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
