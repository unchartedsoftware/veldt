package binning

import (
	"math"

	"github.com/unchartedsoftware/prism/util"
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

// CoordToFractionalTile converts a data coordniate to a floating point tile coordinate.
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

// CoordToTile converts a data coordniate to a tile coordinate.
func CoordToTile(coord *Coord, level uint32, bounds *Bounds) *TileCoord {
	tile := CoordToFractionalTile(coord, level, bounds)
	return &TileCoord{
		X: uint32(math.Floor(tile.X)),
		Y: uint32(math.Floor(tile.Y)),
		Z: level,
	}
}

// CoordToFractionalBin converts a data coordniate to a floating point bin coordinate.
func CoordToFractionalBin(coord *Coord, level uint32, numBins uint32, bounds *Bounds) *FractionalBinCoord {
	tile := CoordToFractionalTile(coord, level, bounds)
	fbins := float64(numBins)
	return &FractionalBinCoord{
		X: util.Fract(tile.X) * fbins,
		Y: util.Fract(tile.Y) * fbins,
	}
}

// CoordToBin converts a data coordniate to a bin coordinate.
func CoordToBin(coord *Coord, level uint32, numBins uint32, bounds *Bounds) *BinCoord {
	bin := CoordToFractionalBin(coord, level, numBins, bounds)
	return &BinCoord{
		X: uint32(math.Floor(bin.X)),
		Y: uint32(math.Floor(bin.Y)),
	}
}

// GetTileBounds returns the data coordniate bounds of the tile coordinate.
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
