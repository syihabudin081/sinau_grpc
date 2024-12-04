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

// Scope menambahkan filter ke dalam query
func (qb *QueryBuilderNative) Scope(condition string, args ...interface{}) *QueryBuilderNative {
	if condition != "" {
		qb.Scopes = append(qb.Scopes, condition)
		qb.Args = append(qb.Args, args...)
	}
	return qb
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

// Build menyusun query akhir
func (qb *QueryBuilderNative) Build() (string, []interface{}) {
	query := qb.BaseQuery
	if len(qb.Scopes) > 0 {
		query += " WHERE " + strings.Join(qb.Scopes, " AND ")
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
