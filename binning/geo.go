package binning

import (
	"math"

	"github.com/unchartedsoftware/prism/util"
)

const (
	degreesToRadians = math.Pi / 180.0 // Factor for changing degrees to radians
	radiansToDegrees = 180.0 / math.Pi // Factor for changing radians to degrees
)

// GeoBounds represents a geographical bounding box.
type GeoBounds struct {
	TopLeft     *LonLat
	BottomRight *LonLat
}

// LonLat represents a geographic point.
type LonLat struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

func tileToLon(x float64, level uint32) float64 {
	pow2 := math.Pow(2, float64(level))
	return x/pow2*360.0 - 180.0
}

func tileToLat(y float64, level uint32) float64 {
	pow2 := math.Pow(2, float64(level))
	n := math.Pi - (2.0*math.Pi*y)/pow2
	return math.Atan(math.Sinh(n)) * radiansToDegrees
}

// LonLatToFractionalTile converts a geographic coordinate into a floating point tile coordinate.
func LonLatToFractionalTile(lonLat *LonLat, level uint32) *FractionalTileCoord {
	latR := lonLat.Lat * degreesToRadians
	pow2 := math.Pow(2, float64(level))
	x := (lonLat.Lon + 180.0) / 360.0 * pow2
	y := (pow2 * (1 - math.Log(math.Tan(latR)+1/math.Cos(latR))/math.Pi) / 2)
	return &FractionalTileCoord{
		X: x,
		Y: y,
		Z: level,
	}
}

// LonLatToTile converts a geographic coordinate into tile coordinate.
func LonLatToTile(lonlat *LonLat, level uint32) *TileCoord {
	tile := LonLatToFractionalTile(lonlat, level)
	return &TileCoord{
		X: uint32(math.Floor(tile.X)),
		Y: uint32(math.Floor(tile.Y)),
		Z: uint32(tile.Z),
	}
}

// LonLatToFractionalBin converts a geographic coordinate into a floating point bin coordinate.
func LonLatToFractionalBin(lonlat *LonLat, level uint32, numBins uint32) *FractionalBinCoord {
	tile := LonLatToFractionalTile(lonlat, level)
	fbins := float64(numBins)
	return &FractionalBinCoord{
		X: util.Fract(tile.X) * fbins,
		Y: util.Fract(tile.Y) * fbins,
	}
}

// LonLatToBin converts a geographic coordinate into a bin coordinate.
func LonLatToBin(lonlat *LonLat, level uint32, numBins uint32) *BinCoord {
	bin := LonLatToFractionalBin(lonlat, level, numBins)
	return &BinCoord{
		X: uint32(math.Floor(bin.X)),
		Y: uint32(math.Floor(bin.Y)),
	}
}

// GetTileGeoBounds returns the geographic bounds of the tile coordinate.
func GetTileGeoBounds(tile *TileCoord) *GeoBounds {
	top := tileToLat(float64(tile.Y), tile.Z)
	bottom := tileToLat(float64(tile.Y+1), tile.Z)
	right := tileToLon(float64(tile.X+1), tile.Z)
	left := tileToLon(float64(tile.X), tile.Z)
	return &GeoBounds{
		TopLeft: &LonLat{
			Lon: left,
			Lat: top,
		},
		BottomRight: &LonLat{
			Lon: right,
			Lat: bottom,
		},
	}
}
