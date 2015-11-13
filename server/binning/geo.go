package binning

import (
	"math"
)

// GeoBounds represents a geographical bounding box
type GeoBounds struct {
	BottomLeft *LonLat
	TopRight *LonLat
}

// LonLat represents a geographic point
type LonLat struct {
	Lon float64
	Lat float64
}

const degreesToRadians = math.Pi / 180.0 // Factor for changing degrees to radians
const radiansToDegrees = 180.0 / math.Pi // Factor for changing radians to degrees

func tileToLon( x uint32, level uint32 ) float64 {
	pow2 := math.Pow( 2, float64( level ) )
	return float64( x ) / pow2 * 360.0 - 180.0
}

func tileToLat( y uint32, level uint32 ) float64 {
	pow2 := math.Pow( 2, float64( level ) )
	n := math.Pi - ( 2.0 * math.Pi * float64( y ) ) / pow2
	return math.Atan( math.Sinh( n ) ) * radiansToDegrees
	/*
	pow2 := math.Pow( 2, float64( level ) )
	n := -math.Pi + ( 2.0 * math.Pi * float64( y ) ) / pow2
	return math.Atan( math.Sinh( n ) ) * radiansToDegrees
	*/
}

func fract( value float64 ) float64 {
	_, fract := math.Modf( value )
	return fract
}

// LonLatToFractionalTile converts a geograhic coordinate into a floating point tile coordinate
func LonLatToFractionalTile( lonLat *LonLat, level uint32 ) *FractionalTileCoord {
	latR := lonLat.Lat * degreesToRadians
	pow2 := math.Pow( 2, float64( level ) )
	x := ( lonLat.Lon + 180.0 ) / 360.0 * pow2
	y := ( pow2 * ( 1 - math.Log( math.Tan( latR ) + 1 / math.Cos( latR ) ) / math.Pi ) / 2 )
	return &FractionalTileCoord{
		X: x,
		Y: y, //pow2 - y,
		Z: level,
    }
}

// LonLatToTile converts a geograhic coordinate into tile coordinate
func LonLatToTile( lonlat *LonLat, level uint32 ) *TileCoord {
	tile := LonLatToFractionalTile( lonlat, level )
	return &TileCoord{
		X: uint32( math.Floor( tile.X ) ),
		Y: uint32( math.Floor( tile.Y ) ),
		Z: uint32( tile.Z ),
	}
}

// LonLatToFractionalBin converts a geograhic coordinate into a floating point bin coordinate
func LonLatToFractionalBin( lonlat *LonLat, level uint32, numBins uint32 ) *FractionalBinCoord {
	tile := LonLatToFractionalTile( lonlat, level )
	fbins := float64( numBins )
    return &FractionalBinCoord{
        X: fract( tile.X ) * fbins,
        Y: fract( tile.Y ) * fbins, //( fbins - 1 ) - fract( tile.Y ) * ( fbins - 1 ),
    }
}

// LonLatToBin converts a geograhic coordinate into a bin coordinate
func LonLatToBin( lonlat *LonLat, level uint32, numBins uint32 ) *BinCoord {
	bin := LonLatToFractionalBin( lonlat, level, numBins )
	return &BinCoord{
		X: uint32( math.Floor( bin.X ) ),
		Y: uint32( math.Floor( bin.Y ) ),
	}
}

func LonLatToFlatBin( lonlat *LonLat, level uint32, numBins uint32 ) uint32 {
	bin := LonLatToFractionalBin( lonlat, level, numBins )
	return uint32( math.Floor( bin.X ) ) +
		uint32( math.Floor( bin.Y ) ) * numBins
}

// GetTileGeoBounds returns the geographic bounds of the tile coordinate
func GetTileGeoBounds( tile *TileCoord ) *GeoBounds {
	top := tileToLat( tile.Y, tile.Z )
    bottom := tileToLat( tile.Y + 1, tile.Z )
    // top := tileToLat( tile.Y + 1, tile.Z )
    // bottom := tileToLat( tile.Y, tile.Z )
    right := tileToLon( tile.X + 1, tile.Z )
    left := tileToLon( tile.X, tile.Z )
    return &GeoBounds{
		BottomLeft: &LonLat{
			Lon: left,
			Lat: bottom,
		},
		TopRight: &LonLat{
			Lon: right,
			Lat: top,
		},
    }
}
