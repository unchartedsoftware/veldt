package binning

import (
	"math"
)

const (
	degreesToRadians = math.Pi / 180.0 // Factor for changing degrees to radians
	radiansToDegrees = 180.0 / math.Pi // Factor for changing radians to degrees
)

// GeoBounds represents a geographical bounding box.
type GeoBounds struct {
	BottomLeft *LonLat
	TopRight   *LonLat
}

// LonLat represents a geographic point.
type LonLat struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

// NewLonLat instantiates and returns a pointer to a LonLat.
func NewLonLat(lon, lat float64) *LonLat {
	return &LonLat{
		Lon: math.Min(180, math.Max(-180, lon)),
		Lat: math.Min(85.0, math.Max(-85.0, lat)),
	}
}

// LonLatToFractionalTile converts a geographic coordinate into a floating point tile coordinate.
func LonLatToFractionalTile(lonLat *LonLat, level uint32) *FractionalTileCoord {
	latR := lonLat.Lat * degreesToRadians
	pow2 := math.Pow(2, float64(level))
	x := (lonLat.Lon + 180.0) / 360.0 * pow2
	y := (pow2 * (1 - math.Log(math.Tan(latR)+1/math.Cos(latR))/math.Pi) / 2)
	return &FractionalTileCoord{
		X: x,
		Y: pow2 - y,
		Z: level,
	}
}
