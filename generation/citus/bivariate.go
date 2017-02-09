package citus

import (
	"fmt"
	"math"

	"github.com/jackc/pgx"

	"github.com/unchartedsoftware/veldt/binning"
	"github.com/unchartedsoftware/veldt/tile"
)

// Bivariate represents a bivariate tile generator.
type Bivariate struct {
	tile.Bivariate
}

// AddQuery adds the tiling query to the provided query object.
func (b *Bivariate) AddQuery(coord *binning.TileCoord, query *Query) *Query {
	// get tile bounds
	bounds := b.TileBounds(coord)
	// x
	minXArg := query.AddParameter(int64(bounds.MinX()))
	maxXArg := query.AddParameter(int64(bounds.MaxX()))
	rangeQueryX := fmt.Sprintf("%s >= %s and %s < %s", b.XField, minXArg, b.XField, maxXArg)
	query.Where(rangeQueryX)
	// y
	minYArg := query.AddParameter(int64(bounds.MinY()))
	maxYArg := query.AddParameter(int64(bounds.MaxY()))
	rangeQueryY := fmt.Sprintf("%s >= %s and %s < %s", b.YField, minYArg, b.YField, maxYArg)
	query.Where(rangeQueryY)
	// result
	return query
}

// AddAggs adds the tiling aggregations to the provided query object.
func (b *Bivariate) AddAggs(coord *binning.TileCoord, query *Query) *Query {
	bounds := b.TileBounds(coord)
	// bin
	minX := int64(bounds.MinX())
	minY := int64(bounds.MinY())
	intervalX := int64(math.Max(1, b.BinSizeX(coord)))
	intervalY := int64(math.Max(1, b.BinSizeY(coord)))
	// x
	minXArg := query.AddParameter(minX)
	intervalXArg := query.AddParameter(intervalX)
	queryString := fmt.Sprintf("((%s - %s) / %s * %s)", b.XField, minXArg, intervalXArg, intervalXArg)
	query.GroupBy(queryString)
	query.Select(fmt.Sprintf("%s + %s as x", minXArg, queryString))
	// y
	minYArg := query.AddParameter(minY)
	intervalYArg := query.AddParameter(intervalY)
	queryString = fmt.Sprintf("((%s - %s) / %s * %s)", b.YField, minYArg, intervalYArg, intervalYArg)
	query.GroupBy(queryString)
	query.Select(fmt.Sprintf("%s + %s as y", minYArg, queryString))
	// result
	return query
}

// GetBins parses the resulting histograms into bins.
func (b *Bivariate) GetBins(coord *binning.TileCoord, rows *pgx.Rows) ([]float64, error) {
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

		xBin := b.GetXBin(coord, float64(x))
		yBin := b.GetYBin(coord, float64(y))

		index := xBin + b.Resolution*yBin
		bins[index] += value
	}

	return bins, nil
}
