package query

import (
	"fmt"
	"strconv"
	"strings"
)

type Query struct {
	QueryArgs      []interface{}
	WhereClauses   []string
	GroupByClauses []string
	Fields         []string
	Tables         []string
	OrderByClauses []string
	Limit		   uint32
}

func NewQuery() (*Query, error) {
	return &Query{
		WhereClauses:   []string{},
		GroupByClauses: []string{},
		Fields:         []string{},
		Tables:         []string{},
		OrderByClauses: []string{},
		Limit:			0,
		QueryArgs:      make([]interface{}, 0),
	}, nil
}

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
	if q.Limit > 0 {
		hash = hash +"::Limit=" + fmt.Sprint(q.Limit)
	}

	return hash
}

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

	if q.Limit > 0 {
		queryString += fmt.Sprintf(" LIMIT %s", q.Limit)
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

func (q *Query) AddWhereClause(clause string) {
	q.WhereClauses = append(q.WhereClauses, clause)
}

func (q *Query) AddGroupByClause(clause string) {
	q.GroupByClauses = append(q.GroupByClauses, clause)
}

func (q *Query) AddField(field string) {
	q.Fields = append(q.Fields, field)
}

func (q *Query) AddTable(table string) {
	q.Tables = append(q.Tables, table)
}

func (q *Query) AddOrderByClause(clause string) {
	q.OrderByClauses = append(q.OrderByClauses, clause)
}

func (q *Query) SetLimit(limit uint32) {
	q.Limit = limit
}
