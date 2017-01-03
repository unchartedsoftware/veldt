package common

import (
	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/tile"
)

// Bivariate represents a bivariate tile generator.
type Bivariate struct {
	tile.Bivariate
	Bounds   *binning.Bounds
	BinSizeX float64
	BinSizeY float64
}

func (b *Bivariate) ClampBin(bin int64) int {
	if bin > int64(b.Resolution)-1 {
		return b.Resolution - 1
	}
	if bin < 0 {
		return 0
	}
	return int(bin)
}

func (b *Bivariate) GetXBin(x int64) int {
	bounds := b.Bounds
	fx := float64(x)
	var bin int64
	if bounds.BottomLeft.X > bounds.TopRight.X {
		bin = int64(float64(b.Resolution-1) - ((fx - bounds.TopRight.X) / b.BinSizeX))
	} else {
		bin = int64((fx - bounds.BottomLeft.X) / b.BinSizeX)
	}
	return b.ClampBin(bin)
}

// GetX given an x value, returns the corresponding coord within the range of
// [0 : 256) for the tile.
func (b *Bivariate) GetX(x float64) float64 {
	bounds := b.Bounds
	if bounds.BottomLeft.X > bounds.TopRight.X {
		rang := bounds.BottomLeft.X - bounds.TopRight.X
		return binning.MaxTileResolution - (((x - bounds.TopRight.X) / rang) * binning.MaxTileResolution)
	}
	rang := bounds.TopRight.X - bounds.BottomLeft.X
	return ((x - bounds.BottomLeft.X) / rang) * binning.MaxTileResolution
}

// GetYBin given an y value, returns the corresponding bin.
func (b *Bivariate) GetYBin(y int64) int {
	bounds := b.Bounds
	fy := float64(y)
	var bin int64
	if bounds.BottomLeft.Y > bounds.TopRight.Y {
		bin = int64(float64(b.Resolution-1) - ((fy - bounds.TopRight.Y) / b.BinSizeY))
	} else {
		bin = int64((fy - bounds.BottomLeft.Y) / b.BinSizeY)
	}
	return b.ClampBin(bin)
}

// GetY given an y value, returns the corresponding coord within the range of
// [0 : 256) for the tile.
func (b *Bivariate) GetY(y float64) float64 {
	bounds := b.Bounds
	if bounds.BottomLeft.Y > bounds.TopRight.Y {
		rang := bounds.BottomLeft.Y - bounds.TopRight.Y
		return binning.MaxTileResolution - (((y - bounds.TopRight.Y) / rang) * binning.MaxTileResolution)
	}
	rang := bounds.TopRight.Y - bounds.BottomLeft.Y
	return ((y - bounds.BottomLeft.Y) / rang) * binning.MaxTileResolution
}
