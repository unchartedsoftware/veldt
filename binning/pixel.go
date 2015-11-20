package binning

import (
	"math"

	"github.com/unchartedsoftware/prism/util"
)

// PixelBounds represents a bounding box in pixel coordinates.
type PixelBounds struct {
	TopLeft *PixelCoord
	BottomRight *PixelCoord
}

// PixelCoord represents a point in pixel coordinates.
type PixelCoord struct {
	X uint64 `json:"x"`
	Y uint64 `json:"y"`
}

// LonLatToPixelCoord translates a geographic coordinate to a pixel coordinate.
func LonLatToPixelCoord( lonLat *LonLat, level uint32, tileResolution uint32 ) *PixelCoord {
    pow2 := math.Pow( 2, float64( level ) )
    // Converting to range from [0:1] where 0,0 is top left
    normalizedTile := LonLatToFractionalTile( lonLat, 0 )
    normalizedCoord := &Coord{
        X: normalizedTile.X,
        Y: normalizedTile.Y,
    }
    return &PixelCoord{
        X: uint64( math.Floor( normalizedCoord.X * float64( tileResolution ) * pow2 ) ),
        Y: uint64( math.Floor( normalizedCoord.Y * float64( tileResolution ) * pow2 ) ),
    }
}

// CoordToPixelCoord translates a coordinate to a pixel coordinate.
func CoordToPixelCoord( coord *Coord, level uint32, bounds *Bounds, tileResolution uint32 ) *PixelCoord {
    pow2 := math.Pow( 2, float64( level ) )
    // Converting to range from [0:1] where 0,0 is top left
    normalizedTile := CoordToFractionalTile( coord, 0, bounds )
    normalizedCoord := &Coord{
        X: normalizedTile.X,
        Y: normalizedTile.Y,
    }
    return &PixelCoord{
        X: uint64( math.Floor( normalizedCoord.X * float64( tileResolution ) * pow2 ) ),
        Y: uint64( math.Floor( normalizedCoord.Y * float64( tileResolution ) * pow2 ) ),
    }
}

// GetTilePixelBounds returns the pixel coordniate bounds of the tile coordinate.
func GetTilePixelBounds( tile *TileCoord, level uint32, tileResolution uint32 ) *PixelBounds {
	sourcePow2 := math.Pow( 2, float64( tile.Z ) )
	// Converting to range from [0:1] where 0,0 is top left
	xMin := float64( tile.X ) / sourcePow2
	xMax := float64( tile.X + 1 ) / sourcePow2
	yMin := float64( tile.Y ) / sourcePow2
	yMax := float64( tile.Y + 1 ) / sourcePow2
	// Projecting to pixel coords
	destPow2 := util.Round( math.Pow( 2, float64( level ) ) * float64( tileResolution ) )
	return &PixelBounds{
		TopLeft: &PixelCoord{
			X: uint64( util.Round( xMin * destPow2 ) ),
			Y: uint64( util.Round( yMin * destPow2 ) ),
		},
		BottomRight: &PixelCoord{
			X: uint64( util.Round( xMax * destPow2 ) ),
			Y: uint64( util.Round( yMax * destPow2 ) ),
		},
	}
}
