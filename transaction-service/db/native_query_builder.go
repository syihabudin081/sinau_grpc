package db

import (
	"fmt"
	"strings"
)

type QueryBuilderNative struct {
	BaseQuery     string
	Scopes        []string
	Args          []interface{}
	OrderByClause string
	Limit         int
	Offset        int
}

func NewQueryBuilderNative(baseQuery string) *QueryBuilderNative {
	return &QueryBuilderNative{
		BaseQuery: baseQuery,
		Scopes:    []string{},
		Args:      []interface{}{},
	}
}

// OrderBy menambahkan pengurutan
func (qb *QueryBuilderNative) OrderBy(column, direction string) *QueryBuilderNative {
	if column != "" && (direction == "ASC" || direction == "DESC") {
		qb.OrderByClause = fmt.Sprintf("ORDER BY %s %s", column, direction)
	}
	return qb
}

// Pagination menambahkan limit dan offset
func (qb *QueryBuilderNative) Pagination(limit, offset int) *QueryBuilderNative {
	qb.Limit = limit
	qb.Offset = offset
	return qb
}

// Scope menambahkan filter ke dalam query
func (qb *QueryBuilderNative) Scope(condition string, args ...interface{}) *QueryBuilderNative {
	// Only add condition if it's not empty and args are not nil
	if condition != "" && len(args) > 0 && args[0] != nil && args[0] != "0" {
		qb.Scopes = append(qb.Scopes, condition)
		qb.Args = append(qb.Args, args...)
	}
	return qb
}

// Build menyusun query akhir
func (qb *QueryBuilderNative) Build() (string, []interface{}) {
	query := qb.BaseQuery

	// Filter out empty conditions
	var validScopes []string
	var validArgs []interface{}
	for i, scope := range qb.Scopes {
		if scope != "" {
			validScopes = append(validScopes, scope)
			validArgs = append(validArgs, qb.Args[i])
		}
	}

	// Only add WHERE clause if there are valid conditions
	if len(validScopes) > 0 {
		query += " WHERE " + strings.Join(validScopes, " AND ")
		qb.Args = validArgs
	}

	if qb.OrderByClause != "" {
		query += " " + qb.OrderByClause
	}

	if qb.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", qb.Limit)
		if qb.Offset > 0 {
			query += fmt.Sprintf(" OFFSET %d", qb.Offset)
		}
	}

	return query, qb.Args
}
