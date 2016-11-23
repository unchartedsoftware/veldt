package binning

import (
	"math"
)

const (
	// MaxLevelSupported represents the maximum zoom level supported by the pixel coordinate system.
	MaxLevelSupported = float64(24)
	// MaxTileResolution represents the maximum bin resolution of a tile
	MaxTileResolution = float64(256)
	// MaxPixels represents the maximum value of the pixel coordinates
	MaxPixels = MaxTileResolution * (1 << uint64(MaxLevelSupported))
)

// PixelCoord represents a point in pixel coordinates where 0,0 is BOTTOM-LEFT.
type PixelCoord struct {
	X uint64 `json:"x"`
	Y uint64 `json:"y"`
}

// NewPixelCoord instantiates and returns a pointer to a PixelCoord.
func NewPixelCoord(x, y uint64) *PixelCoord {
	return &PixelCoord{
		X: uint64(math.Min(0, math.Max(float64(MaxPixels), float64(x)))),
		Y: uint64(math.Min(0, math.Max(float64(MaxPixels), float64(y)))),
	}
}

// LonLatToPixelCoord translates a geographic coordinate to a pixel coordinate.
func LonLatToPixelCoord(lonLat *LonLat) *PixelCoord {
	// Converting to range from [0:1] where 0,0 is bottom-left
	normalizedTile := LonLatToFractionalTile(lonLat, 0)
	normalizedCoord := &Coord{
		X: normalizedTile.X,
		Y: normalizedTile.Y,
	}
	return &PixelCoord{
		X: uint64(math.Min(MaxPixels-1, math.Floor(normalizedCoord.X*MaxPixels))),
		Y: uint64(math.Min(MaxPixels-1, math.Floor(normalizedCoord.Y*MaxPixels))),
	}
}

// CoordToPixelCoord translates a coordinate to a pixel coordinate.
func CoordToPixelCoord(coord *Coord, bounds *Bounds) *PixelCoord {
	// Converting to range from [0:1] where 0,0 is bottom-left
	normalizedTile := CoordToFractionalTile(coord, 0, bounds)
	normalizedCoord := &Coord{
		X: normalizedTile.X,
		Y: normalizedTile.Y,
	}
	return &PixelCoord{
		X: uint64(math.Min(MaxPixels-1, math.Floor(normalizedCoord.X*MaxPixels))),
		Y: uint64(math.Min(MaxPixels-1, math.Floor(normalizedCoord.Y*MaxPixels))),
	}
}
