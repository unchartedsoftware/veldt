package citus

import (
	"fmt"
	"math"

	"github.com/jackc/pgx"

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

func (b *Bivariate) AddQuery(coord *binning.TileCoord, query *Query) *Query {

	extents := &binning.Bounds{
		BottomLeft: &binning.Coord{
			X: b.Left,
			Y: b.Top,
		},
		TopRight: &binning.Coord{
			X: b.Right,
			Y: b.Bottom,
		},
	}
	bounds := binning.GetTileBounds(coord, extents)
	minX := int64(math.Min(bounds.BottomLeft.X, bounds.TopRight.X))
	maxX := int64(math.Max(bounds.BottomLeft.X, bounds.TopRight.X))
	minY := int64(math.Min(bounds.BottomLeft.Y, bounds.TopRight.Y))
	maxY := int64(math.Max(bounds.BottomLeft.Y, bounds.TopRight.Y))
	b.Bounds = bounds

	minXArg := query.AddParameter(minX)
	maxXArg := query.AddParameter(maxX)
	rangeQueryX := fmt.Sprintf("%s >= %s and %s < %s", b.XField, minXArg, b.XField, maxXArg)
	query.Where(rangeQueryX)

	minYArg := query.AddParameter(minY)
	maxYArg := query.AddParameter(maxY)
	rangeQueryY := fmt.Sprintf("%s >= %s and %s < %s", b.YField, minYArg, b.YField, maxYArg)
	query.Where(rangeQueryY)

	return query
}

func (b *Bivariate) AddAgg(coord *binning.TileCoord, query *Query) *Query {

	extents := &binning.Bounds{
		BottomLeft: &binning.Coord{
			X: b.Left,
			Y: b.Top,
		},
		TopRight: &binning.Coord{
			X: b.Right,
			Y: b.Bottom,
		},
	}
	bounds := binning.GetTileBounds(coord, extents)
	minX := int64(math.Min(bounds.BottomLeft.X, bounds.TopRight.X))
	minY := int64(math.Min(bounds.BottomLeft.Y, bounds.TopRight.Y))
	xRange := math.Abs(bounds.TopRight.X - bounds.BottomLeft.X)
	yRange := math.Abs(bounds.TopRight.Y - bounds.BottomLeft.Y)
	intervalX := int64(math.Max(1, xRange/float64(b.Resolution)))
	intervalY := int64(math.Max(1, yRange/float64(b.Resolution)))
	b.Bounds = bounds
	b.BinSizeX = xRange / float64(b.Resolution)
	b.BinSizeY = yRange / float64(b.Resolution)

	minXArg := query.AddParameter(minX)
	intervalXArg := query.AddParameter(intervalX)
	queryString := fmt.Sprintf("((%s - %s) / %s * %s)", b.XField, minXArg, intervalXArg, intervalXArg)
	query.GroupBy(queryString)
	query.Select(fmt.Sprintf("%s + %s as x", minXArg, queryString))

	minYArg := query.AddParameter(minY)
	intervalYArg := query.AddParameter(intervalY)
	queryString = fmt.Sprintf("((%s - %s) / %s * %s)", b.YField, minYArg, intervalYArg, intervalYArg)
	query.GroupBy(queryString)
	query.Select(fmt.Sprintf("%s + %s as y", minYArg, queryString))

	return query
}

func (b *Bivariate) clampBin(bin int64) int64 {
	if bin > int64(b.Resolution)-1 {
		return int64(b.Resolution) - 1
	}
	if bin < 0 {
		return 0
	}
	return bin
}

func (b *Bivariate) GetXBin(x int64) int64 {
	bounds := b.Bounds
	fx := float64(x)
	var bin int64
	if bounds.BottomLeft.X > bounds.TopRight.X {
		bin = int64(float64(b.Resolution-1) - ((fx - bounds.TopRight.X) / b.BinSizeX))
	} else {
		bin = int64((fx - bounds.BottomLeft.X) / b.BinSizeX)
	}
	return b.clampBin(bin)
}

// GetYBin given an y value, returns the corresponding bin.
func (b *Bivariate) GetYBin(y int64) int64 {
	bounds := b.Bounds
	fy := float64(y)
	var bin int64
	if bounds.BottomLeft.Y > bounds.TopRight.Y {
		bin = int64(float64(b.Resolution-1) - ((fy - bounds.TopRight.Y) / b.BinSizeY))
	} else {
		bin = int64((fy - bounds.BottomLeft.Y) / b.BinSizeY)
	}
	return b.clampBin(bin)
}

// GetBins parses the resulting histograms into bins.
func (b *Bivariate) GetBins(rows *pgx.Rows) ([]float64, error) {
	// allocate bins buffer
	bins := make([]float64, b.Resolution*b.Resolution)
	// fill bins buffer
	for rows.Next() {
		var x, y int64
		var value float64

		err := rows.Scan(&x, &y, &value)
		if err != nil {
			return nil, fmt.Errorf("Error parsing histogram aggregation: %v",
				err)
		}

		xBin := b.GetXBin(x)
		yBin := b.GetYBin(y)

		index := xBin + int64(b.Resolution)*yBin
		bins[index] += value
	}

	return bins, nil
}
