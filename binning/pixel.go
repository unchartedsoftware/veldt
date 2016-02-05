package binning

import (
	"math"

	"github.com/unchartedsoftware/prism/util"
)

const (
	// MaxLevelSupported represents the maximum zoom level supported by the pixel coordinate system.
	MaxLevelSupported = float64(24)
	// MaxTileResolution represents the maximum bin resolution of a tile
	MaxTileResolution = float64(256)
)

var (
	// MaxPixels represents the maximum value of the pixel coordinates
	MaxPixels = MaxTileResolution * math.Pow(2, MaxLevelSupported)
)

// PixelBounds represents a bounding box in pixel coordinates.
type PixelBounds struct {
	TopLeft     *PixelCoord
	BottomRight *PixelCoord
}

// PixelCoord represents a point in pixel coordinates.
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
	// Converting to range from [0:1] where 0,0 is top left
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
	// Converting to range from [0:1] where 0,0 is top left
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

// GetTilePixelBounds returns the pixel coordinate bounds of the tile coordinate.
func GetTilePixelBounds(tile *TileCoord) *PixelBounds {
	pow2 := math.Pow(2, float64(tile.Z))
	// Converting to range from [0:1] where 0,0 is top left
	xMin := float64(tile.X) / pow2
	xMax := float64(tile.X+1) / pow2
	yMin := float64(tile.Y) / pow2
	yMax := float64(tile.Y+1) / pow2
	return &PixelBounds{
		TopLeft: &PixelCoord{
			X: uint64(math.Min(MaxPixels-1, util.Round(xMin*MaxPixels))),
			Y: uint64(math.Min(MaxPixels-1, util.Round(yMin*MaxPixels))),
		},
		BottomRight: &PixelCoord{
			X: uint64(math.Min(MaxPixels-1, util.Round(xMax*MaxPixels))),
			Y: uint64(math.Min(MaxPixels-1, util.Round(yMax*MaxPixels))),
		},
	}
}
