package binning

// TileCoord represents a TMS tile's coordinates (0,0) being at the bottom-left.
type TileCoord struct {
	X uint32
	Y uint32
	Z uint32
}

// BinCoord represents a bin position within a Tile.
type BinCoord struct {
	X uint32
	Y uint32
}

// FractionalTileCoord represents a TMS tile's coordinates, using floating point components (0,0) being at the bottom-left.
type FractionalTileCoord struct {
	X float64
	Y float64
	Z uint32
}

// FractionalBinCoord represents a bin position within a Tile, using floating point components.
type FractionalBinCoord struct {
	X float64
	Y float64
}
