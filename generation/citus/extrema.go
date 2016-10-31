package citus

import (
	"fmt"
	"time"

	"github.com/jackc/pgx"

	"github.com/unchartedsoftware/prism/binning"
)

// GetNumericExtrema returns the extrema of a numeric field for the provided index.
func GetNumericExtrema(connPool *pgx.ConnPool, schema string, table string, column string) (*binning.Extrema, error) {
	// query
	queryString := fmt.Sprintf("SELECT CAST(MIN(%s) AS FLOAT) as min, CAST(MAX(%s) AS FLOAT) as max FROM %s.%s;", column, column, schema, table)
	row := connPool.QueryRow(queryString)

	// Parse min & max values.
	var min *float64
	var max *float64
	err := row.Scan(&min, &max)
	if err != nil {
		return nil, err
	}

	// it seems if the mapping exists, but no documents have the attribute, the min / max are null
	// TODO: TEST THIS FOR CITUS!!!
	if min == nil || max == nil {
		return nil, nil
	}

	return &binning.Extrema{
		Min: *min,
		Max: *max,
	}, nil
}

// GetTimestampExtrema returns the extrema of a timestamp field for the provided index.
func GetTimestampExtrema(connPool *pgx.ConnPool, schema string, table string, column string) (*binning.Extrema, error) {
	// query
	queryString := fmt.Sprintf("SELECT MIN(%s) as min, MAX(%s) as max FROM %s.%s;", column, column, schema, table)
	row := connPool.QueryRow(queryString)

	// Parse min & max values.
	var min *time.Time
	var max *time.Time
	err := row.Scan(&min, &max)
	if err != nil {
		return nil, err
	}

	// it seems if the mapping exists, but no documents have the attribute, the min / max are null
	// TODO: TEST THIS FOR CITUS!!!
	if min == nil || max == nil {
		return nil, nil
	}

	return &binning.Extrema{
		Min: float64(min.Unix()),
		Max: float64(max.Unix()),
	}, nil
}
