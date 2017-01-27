package citus

import (
	"fmt"
	"strconv"
	"strings"
)

// QueryString represents a citus implementation of the veldt.Query interface.
type QueryString interface {
	Get(*Query) (string, error)
}

// Query represents a citus query object.
type Query struct {
	QueryArgs      []interface{}
	WhereClauses   []string
	GroupByClauses []string
	Fields         []string
	Tables         []string
	OrderByClauses []string
	RowLimit       uint32
}

// NewQuery instantiates and returns a new query object.
func NewQuery() (*Query, error) {
	return &Query{
		WhereClauses:   []string{},
		GroupByClauses: []string{},
		Fields:         []string{},
		Tables:         []string{},
		OrderByClauses: []string{},
		RowLimit:       0,
		QueryArgs:      make([]interface{}, 0),
	}, nil
}

// GetHash returns the string hash of the query
func (q *Query) GetHash() string {
	numWheres := len(q.WhereClauses)
	numGroups := len(q.GroupByClauses)
	numFields := len(q.Fields)
	numTables := len(q.Tables)
	numOrders := len(q.OrderByClauses)
	numArgs := len(q.QueryArgs)
	hashes := make([]string, numWheres+numGroups+numFields+numTables+numOrders+numArgs)
	for i, clause := range q.WhereClauses {
		hashes[i] = clause
	}
	for i, clause := range q.GroupByClauses {
		hashes[i+numWheres] = clause
	}
	for i, clause := range q.Fields {
		hashes[i+numWheres+numGroups] = clause
	}
	for i, clause := range q.Tables {
		hashes[i+numWheres+numGroups+numFields] = clause
	}
	for i, clause := range q.OrderByClauses {
		hashes[i+numWheres+numGroups+numFields+numTables] = clause
	}

	// May want to revisit the query argument hashing.
	for i, arg := range q.QueryArgs {
		hashes[i+numWheres+numGroups+numFields+numOrders+numTables] = fmt.Sprintf("%v", arg)
	}

	hash := strings.Join(hashes, "::")
	if q.RowLimit > 0 {
		hash = hash + "::Limit=" + fmt.Sprint(q.RowLimit)
	}

	return hash
}

// GetQuery returns the query string.
func (q *Query) GetQuery(nested bool) string {
	queryString := fmt.Sprintf("SELECT %s", strings.Join(q.Fields, ", "))

	queryString += fmt.Sprintf(" FROM %s", strings.Join(q.Tables, ", "))

	if len(q.WhereClauses) > 0 {
		queryString += fmt.Sprintf(" WHERE %s", strings.Join(q.WhereClauses, " AND "))
	}

	if len(q.GroupByClauses) > 0 {
		queryString += fmt.Sprintf(" GROUP BY %s", strings.Join(q.GroupByClauses, ", "))
	}

	if len(q.OrderByClauses) > 0 {
		queryString += fmt.Sprintf(" ORDER BY %s", strings.Join(q.OrderByClauses, ", "))
	}

	if q.RowLimit > 0 {
		queryString += fmt.Sprintf(" LIMIT %d", q.RowLimit)
	}

	if !nested {
		queryString = queryString + ";"
	}

	return queryString
}

// AddParameter adds a parameter to the query and returns the parameter number.
func (q *Query) AddParameter(param interface{}) string {
	q.QueryArgs = append(q.QueryArgs, param)
	return "$" + strconv.Itoa(len(q.QueryArgs))
}

// Where adds a where clause to the query.
func (q *Query) Where(clause string) {
	q.WhereClauses = append(q.WhereClauses, clause)
}

// GroupBy adds a groupby to the query.
func (q *Query) GroupBy(clause string) {
	q.GroupByClauses = append(q.GroupByClauses, clause)
}

// Select adds a field to the query.
func (q *Query) Select(field string) {
	q.Fields = append(q.Fields, field)
}

// From adds a table to the query.
func (q *Query) From(table string) {
	q.Tables = append(q.Tables, table)
}

// OrderBy adds an orber by clause to the query.
func (q *Query) OrderBy(clause string) {
	q.OrderByClauses = append(q.OrderByClauses, clause)
}

// Limit sets the limit to the query.
func (q *Query) Limit(limit uint32) {
	q.RowLimit = limit
}
