package citus

import (
	"fmt"

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
	maxX := int64(bounds.MaxX())
	minY := int64(bounds.MinY())
	maxY := int64(bounds.MaxY())
	// x_bucket
	minXArg := query.AddParameter(minX)
	maxXArg := query.AddParameter(maxX)
	bucketArg := query.AddParameter(b.Resolution)
	queryString := fmt.Sprintf("width_bucket(%s, %s, %s, %s) - 1 AS x_bucket", b.XField, minXArg, maxXArg, bucketArg)
	query.Select(queryString);
	// y_bucket
	minYArg := query.AddParameter(minY)
	maxYArg := query.AddParameter(maxY)
	queryString = fmt.Sprintf("width_bucket(%s, %s, %s, %s) - 1 AS y_bucket", b.YField, minYArg, maxYArg, bucketArg)
	query.Select(queryString);
	query.GroupBy("x_bucket");
	query.GroupBy("y_bucket");
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

		index := x + int64(b.Resolution)*y
		bins[index] += value
	}

	return bins, nil
}

// GetXBin given an x value, returns the corresponding bin.
// This is a passthru for citus since we are binning from 0 to resolution in the query
func (b *Bivariate) GetXBin(coord *binning.TileCoord, x float64) int {
	return int(x)
}

// GetYBin given a y value, returns the corresponding bin.
// This is a passthru for citus since we are binning from 0 to resolution in the query
func (b *Bivariate) GetYBin(coord *binning.TileCoord, y float64) int {
	return int(y)
}
