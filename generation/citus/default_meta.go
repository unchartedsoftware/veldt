package citus

import (
	"encoding/json"
	"strings"
	"errors"

	"github.com/jackc/pgx"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/generation/meta"
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

// DefaultMeta represents a meta data generator that produces default
// metadata with property types and extrema.
type DefaultMeta struct {
	MetaGenerator
}

// NewDefaultMeta instantiates and returns a pointer to a new generator.
func NewDefaultMeta(host string, port string) meta.GeneratorConstructor {
	return func(metaReq *meta.Request) (meta.Generator, error) {
		client, err := NewClient(host, port)
		if err != nil {
			return nil, err
		}
		m := &DefaultMeta{}
		m.host = host
		m.port = port
		m.req = metaReq
		m.client = client
		return m, nil
	}
}

// GetMeta returns the meta data for a given table.
func (g *DefaultMeta) GetMeta() ([]byte, error) {
	client := g.client
	metaReq := g.req

	split := strings.Split(metaReq.URI, ".")
	if len(split) != 2 {
		return nil, errors.New("Incorrect format for table. Expect 'schema.table'.")
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
