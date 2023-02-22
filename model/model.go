package model

import (
	"context"
)

var _ ParseWhere = SqlBuilder{}
var globalBb dbQuery[Model]

func InitDB(db dbQuery[Model]) {
	globalBb = db
}

type QuerySelect[T any] func(context.Context, string, ...any) ([]T, error)
type QueryGet[T any] func(context.Context, string, ...any) (T, error)

type Model interface {
	PrimaryKey() string
	Table() string
}

type ParseWhere interface {
	ParseWhere(*[][]any) (string, []any, error)
}

type dbQuery[T any] interface {
	Select(context.Context, string, ...any) ([]T, error)
	Get(context.Context, string, ...any) (T, error)
}

type SqlBuilder [][]string
