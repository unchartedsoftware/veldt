package binning

// TileCoord represents a TMS tile's coordinates (0,0) being at the bottom-left.
type TileCoord struct {
	X uint32 `json:"x"`
	Y uint32 `json:"y"`
	Z uint32 `json:"z"`
}

// FractionalTileCoord represents a TMS tile's coordinates, using floating point components (0,0) being at the bottom-left.
type FractionalTileCoord struct {
	X float64
	Y float64
	Z uint32
}
