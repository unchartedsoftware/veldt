package citus

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
)

// PropertyMeta represents the meta data for a single property.
type PropertyMeta struct {
	Type    string           `json:"type"`
	Extrema *binning.Extrema `json:"extrema,omitempty"`
}

func isNumeric(typ string) bool {
	return typ == "smallint" ||
		typ == "integer" ||
		typ == "bigint" ||
		typ == "decimal" ||
		typ == "numeric" ||
		typ == "real" ||
		typ == "double precision" ||
		typ == "smallserial" ||
		typ == "serial" ||
		typ == "bigserial"

}

func isTimestamp(typ string) bool {
	return typ == "timestamp" ||
		typ == "timestamp with time zone" ||
		typ == "date" ||
		typ == "time" ||
		typ == "time with time zone" ||
		typ == "interval"
}

func getPropertyMeta(connPool *pgx.ConnPool, schema string, table string, column string, typ string) (*PropertyMeta, error) {
	p := PropertyMeta{
		Type: typ,
	}
	// if field is 'ordinal', get the extrema
	if isNumeric(typ) {
		extrema, err := GetNumericExtrema(connPool, schema, table, column)
		if err != nil {
			return nil, err
		}
		p.Extrema = extrema
	} else if isTimestamp(typ) {
		extrema, err := GetTimestampExtrema(connPool, schema, table, column)
		if err != nil {
			return nil, err
		}
		p.Extrema = extrema
	}
	return &p, nil
}

// GetNumericExtrema returns the extrema of a numeric field for the provided table.
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

// GetTimestampExtrema returns the extrema of a timestamp field for the provided table.
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

// DefaultMeta represents a meta data generator that produces default
// metadata with property types and extrema.
type DefaultMeta struct {
	Host string
	Port string
}

// NewDefaultMeta instantiates and returns a pointer to a new generator.
func NewDefaultMeta(host string, port string) veldt.MetaCtor {
	return func() (veldt.Meta, error) {
		return &DefaultMeta{
			Host: host,
			Port: port,
		}, nil
	}
}

// Parse parses the provided JSON object and populates the structs attributes.
func (g *DefaultMeta) Parse(params map[string]interface{}) error {
	return nil
}

// Create generates metadata from the provided URI.
func (g *DefaultMeta) Create(uri string) ([]byte, error) {
	client, err := NewClient(g.Host, g.Port)
	if err != nil {
		return nil, err
	}

	split := strings.Split(uri, ".")
	if len(split) != 2 {
		return nil, errors.New("incorrect format for table, expect 'schema.table'")
	}
	schemaInput := split[0]
	tableInput := split[1]

	schemaQuery := "select table_schema as schema, table_name as table, column_name as column, data_type as typ from information_schema.columns where table_schema = $1 and table_name = $2;"
	rows, err := client.Query(schemaQuery, schemaInput, tableInput)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	meta := make(map[string]interface{})
	for rows.Next() {
		var table string
		var schema string
		var column string
		var typ string
		err := rows.Scan(&schema, &table, &column, &typ)
		if err != nil {
			return nil, err
		}

		metaColumn, err := getPropertyMeta(client, schema, table, column, typ)
		if err != nil {
			return nil, err
		}
		meta[column] = metaColumn
	}

	return json.Marshal(meta)
}
