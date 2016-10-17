package query

import (
	"fmt"
	"strings"
	"strconv"

	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/plog"
)

type Query struct {
	QueryArgs      []interface{}
	WhereClauses   []string
	GroupByClauses []string
	Fields         []string
	Tables		   []string
}

func NewQuery(tileReq *tile.Request) (*Query, error) {
	return &Query{
		WhereClauses:   []string{},
		GroupByClauses: []string{},
		Fields:         []string{},
		Tables:		    []string{},
		QueryArgs:	make([]interface{}, 0),
	}, nil
}

func (q *Query) GetHash() string {
	numWheres := len(q.WhereClauses)
	numGroups := len(q.GroupByClauses)
	numFields := len(q.Fields)
	numTables := len(q.Tables)
	numArgs   := len(q.QueryArgs)
	hashes := make([]string, numWheres+numGroups+numFields+numTables+numArgs)
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
		hashes[i+numWheres+numGroups+numTables] = clause
	}

	// May want to revisit the query argument hashing.
	for i, arg := range q.QueryArgs {
		hashes[i+numWheres+numGroups+numFields+numTables] = fmt.Sprintf("%v", arg)
	}

	return strings.Join(hashes, "::")
}

func (q *Query) GetQuery() string {
	queryString := fmt.Sprintf("SELECT %s", strings.Join(q.Fields, ", "))

	queryString += fmt.Sprintf(" FROM %s", strings.Join(q.Tables, ", "))

	if len(q.WhereClauses) > 0 {
		queryString += fmt.Sprintf(" WHERE %s", strings.Join(q.WhereClauses, " AND "))
	}

	if len(q.GroupByClauses) > 0 {
		queryString += fmt.Sprintf(" GROUP BY %s", strings.Join(q.GroupByClauses, ", "))
	}

	log.Infof("Query: %v", queryString)
	for i:= 0 ; i < len(q.QueryArgs) ; i++ {
		log.Infof("Arg %v: %v", i, q.QueryArgs[i])
	}

	return queryString + ";"
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
