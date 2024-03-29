package model

import (
	"context"
)

var _ ParseWhere = SqlBuilder{}
var globalBb dbQuery

func InitDB(db dbQuery) {
	globalBb = db
}

type QueryFn func(context.Context, any, string, ...any) error

type Model interface {
	PrimaryKey() string
	Table() string
}

type ParseWhere interface {
	ParseWhere(*[][]any) (string, []any, error)
}

type AndWhere interface {
	AndWhere(field, operator, val, fieldType string) ParseWhere
}

type dbQuery interface {
	Select(context.Context, any, string, ...any) error
	Get(context.Context, any, string, ...any) error
}

type SqlBuilder [][]string

func Table[T Model]() string {
	var r T
	return r.Table()
}

func PrimaryKey[T Model]() string {
	var r T
	return r.PrimaryKey()
}
