package binning

import (
	"math"
)

const (
	minLon           = -180.0
	maxLon           = 180.0
	minLat           = -85.05112878
	maxLat           = 85.05112878
	degreesToRadians = math.Pi / maxLon // Factor for changing degrees to radians
	radiansToDegrees = maxLon / math.Pi // Factor for changing radians to degrees
)

// LonLat represents a geographic point.
type LonLat struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

// NewLonLat instantiates and returns a pointer to a LonLat.
func NewLonLat(lon, lat float64) *LonLat {
	return &LonLat{
		Lon: math.Min(maxLon, math.Max(minLon, lon)),
		Lat: math.Min(maxLat, math.Max(minLat, lat)),
	}
}

// LonLatToFractionalTile converts a geographic coordinate into a floating point tile coordinate.
func LonLatToFractionalTile(lonLat *LonLat, level uint32) *FractionalTileCoord {
	latR := lonLat.Lat * degreesToRadians
	pow2 := math.Pow(2, float64(level))
	x := (lonLat.Lon + maxLon) / (maxLon * 2) * pow2
	y := (pow2 * (1 - math.Log(math.Tan(latR)+1/math.Cos(latR))/math.Pi) / 2)
	return &FractionalTileCoord{
		X: x,
		Y: pow2 - y,
		Z: level,
	}
}
