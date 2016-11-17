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

func (b *Bivariate) GetQuery(coord *binning.TileCoord, query *Query) *Query {

	extents := &binning.Bounds{
		TopLeft: &binning.Coord{
			X: b.Left,
			Y: b.Top,
		},
		BottomRight: &binning.Coord{
			X: b.Right,
			Y: b.Bottom,
		},
	}
	bounds := binning.GetTileBounds(coord, extents)
	minX := math.Min(bounds.TopLeft.X, bounds.BottomRight.X)
	maxX := math.Max(bounds.TopLeft.X, bounds.BottomRight.X)
	minY := math.Min(bounds.TopLeft.Y, bounds.BottomRight.Y)
	maxY := math.Max(bounds.TopLeft.Y, bounds.BottomRight.Y)
	b.Bounds = bounds

	minXArg := query.AddParameter(minX)
	maxXArg := query.AddParameter(maxX)
	rangeQueryX := fmt.Sprintf("%s >= %s and %s < %s", b.XField, minXArg, b.XField, maxXArg)
	query.AddWhereClause(rangeQueryX)

	minYArg := query.AddParameter(minY)
	maxYArg := query.AddParameter(maxY)
	rangeQueryY := fmt.Sprintf("%s >= %s and %s < %s", b.YField, minYArg, b.YField, maxYArg)
	query.AddWhereClause(rangeQueryY)

	return query
}

func (b *Bivariate) GetAgg(coord *binning.TileCoord, query *Query) *Query {

	extents := &binning.Bounds{
		TopLeft: &binning.Coord{
			X: b.Left,
			Y: b.Top,
		},
		BottomRight: &binning.Coord{
			X: b.Right,
			Y: b.Bottom,
		},
	}
	bounds := binning.GetTileBounds(coord, extents)
	minX := int64(math.Min(bounds.TopLeft.X, bounds.BottomRight.X))
	minY := int64(math.Min(bounds.TopLeft.Y, bounds.BottomRight.Y))
	xRange := math.Abs(bounds.BottomRight.X - bounds.TopLeft.X)
	yRange := math.Abs(bounds.BottomRight.Y - bounds.TopLeft.Y)
	intervalX := int64(math.Max(1, xRange/float64(b.Resolution)))
	intervalY := int64(math.Max(1, yRange/float64(b.Resolution)))
	b.Bounds = bounds
	b.BinSizeX = xRange / float64(b.Resolution)
	b.BinSizeY = yRange / float64(b.Resolution)

	minXArg := query.AddParameter(minX)
	intervalXArg := query.AddParameter(intervalX)
	queryString := fmt.Sprintf("((%s - %s) / %s * %s)", b.XField, minXArg, intervalXArg, intervalXArg)
	query.AddGroupByClause(queryString)
	query.AddField(fmt.Sprintf("%s + %s as x", minXArg, queryString))

	minYArg := query.AddParameter(minY)
	intervalYArg := query.AddParameter(intervalY)
	queryString = fmt.Sprintf("((%s - %s) / %s * %s)", b.YField, minYArg, intervalYArg, intervalYArg)
	query.AddGroupByClause(queryString)
	query.AddField(fmt.Sprintf("%s + %s as y", minYArg, queryString))

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
	if bounds.TopLeft.X > bounds.BottomRight.X {
		bin = int64(float64(b.Resolution) - ((fx - bounds.BottomRight.X) / b.BinSizeX))
	} else {
		bin = int64((fx - bounds.TopLeft.X) / b.BinSizeX)
	}
	return b.clampBin(bin)
}

// GetYBin given an y value, returns the corresponding bin.
func (b *Bivariate) GetYBin(y int64) int64 {
	bounds := b.Bounds
	fy := float64(y)
	var bin int64
	if bounds.TopLeft.Y > bounds.BottomRight.Y {
		bin = int64(float64(b.Resolution) - ((fy - bounds.BottomRight.Y) / b.BinSizeY))
	} else {
		bin = int64((fy - bounds.TopLeft.Y) / b.BinSizeY)
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
