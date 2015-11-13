package binning

import (
	"math"
)

// Bounds represents a bounding box
type Bounds struct {
	BottomLeft *Coord
	TopRight *Coord
}

// Coord represents a point
type Coord struct {
	X float64
	Y float64
}

// CoordToFractionalTile converts a data coordniate to a floating point tile coordinate
func CoordToFractionalTile( coord *Coord, level uint32, bounds *Bounds ) *FractionalTileCoord {
    pow2 := math.Pow( 2, float64( level ) )
    x := pow2 * ( coord.X - bounds.BottomLeft.X ) / ( bounds.TopRight.X - bounds.BottomLeft.X )
    y := pow2 * ( coord.Y - bounds.TopRight.Y ) / ( bounds.BottomLeft.Y - bounds.TopRight.Y )
    return &FractionalTileCoord{
        X: x,
        Y: y,
	    Z: level,
    }
}

// CoordToTile converts a data coordniate to a tile coordinate
func CoordToTile( coord *Coord, level uint32, bounds *Bounds ) *TileCoord {
	tile := CoordToFractionalTile( coord, level, bounds )
	return &TileCoord{
		X: uint32( math.Floor( tile.X ) ),
		Y: uint32( math.Floor( tile.Y ) ),
		Z: level,
	}
}

// CoordToFractionalBin converts a data coordniate to a floating point bin coordinate
func CoordToFractionalBin( coord *Coord, level uint32, numBins uint32, bounds *Bounds ) *FractionalBinCoord {
	tile := CoordToFractionalTile( coord, level, bounds )
	fbins := float64( numBins )
    return &FractionalBinCoord{
        X: fract( tile.X ) * fbins,
        Y: fract( tile.Y ) * fbins - 1,
    }
}

// CoordToBin converts a data coordniate to a bin coordinate
func CoordToBin( coord *Coord, level uint32, numBins uint32, bounds *Bounds ) *BinCoord {
	bin := CoordToFractionalBin( coord, level, numBins, bounds )
	return &BinCoord{
		X: uint32( math.Floor( bin.X ) ),
		Y: uint32( math.Floor( bin.Y ) ),
	}
}

// GetTileBounds returns the data coordniate bounds of the tile coordinate
func GetTileBounds( tile *TileCoord, bounds *Bounds ) *Bounds {
    pow2 := math.Pow( 2, float64( tile.Z ) )
    tileXSize := float64( bounds.TopRight.X - bounds.BottomLeft.X ) / pow2
    tileYSize := float64( bounds.TopRight.Y - bounds.BottomLeft.Y ) / pow2
	return &Bounds{
		BottomLeft: &Coord{
			X: bounds.BottomLeft.X + tileXSize * float64( tile.X ),
			Y: bounds.BottomLeft.Y + tileYSize * float64( tile.Y ),
		},
		TopRight: &Coord{
			X: bounds.BottomLeft.X + tileXSize * float64( tile.X + 1 ),
			Y: bounds.BottomLeft.Y + tileYSize * float64( tile.Y + 1 ),
		},
    }
}
