package model

import (
	"context"
	"fmt"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/safety"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
)

type SqlxQuery struct {
	sqlx *safety.Var[*sqlx.DB]
	UniversalDb
}

func NewSqlxQuery(sqlx *safety.Var[*sqlx.DB], u UniversalDb) *SqlxQuery {
	s := &SqlxQuery{sqlx: sqlx, UniversalDb: u}
	if u.selects == nil {
		s.UniversalDb.selects = s.Selects
	}
	if u.gets == nil {
		s.UniversalDb.gets = s.Gets
	}
	return s
}

func SetSelect(db *SqlxQuery, fn func(context.Context, any, string, ...any) error) {
	db.selects = fn
}
func SetGet(db *SqlxQuery, fn func(context.Context, any, string, ...any) error) {
	db.gets = fn
}

func (r *SqlxQuery) Selects(ctx context.Context, dest any, sql string, params ...any) error {
	v := helper.GetContextVal(ctx, "handle=>", "")
	db := r.sqlx.Load()
	if v != "" {
		switch v {
		case "string":
			return ToMapSlice(ctx, db, dest.(*[]map[string]string), sql, params...)
		case "scanner":
			fn := ctx.Value("fn")
			return Scanner[any](ctx, db, dest, sql, params...)(fn.(func(any)))
		}
	}

	return db.SelectContext(ctx, dest, sql, params...)
}

func (r *SqlxQuery) Gets(ctx context.Context, dest any, sql string, params ...any) error {
	db := r.sqlx.Load()
	v := helper.GetContextVal(ctx, "handle=>", "")
	if v != "" {
		switch v {
		case "string":
			return GetToMap(ctx, db, dest.(*map[string]string), sql, params...)
		}
	}
	return db.GetContext(ctx, dest, sql, params...)
}

func Scanner[T any](ctx context.Context, db *sqlx.DB, v T, s string, params ...any) func(func(T)) error {
	return func(fn func(T)) error {
		rows, err := db.QueryxContext(ctx, s, params...)
		if err != nil {
			return err
		}
		for rows.Next() {
			err = rows.StructScan(v)
			if err != nil {
				return err
			}
			fn(v)
		}
		return nil
	}
}

func ToMapSlice[V any](ctx context.Context, db *sqlx.DB, dest *[]map[string]V, sql string, params ...any) (err error) {
	rows, err := db.QueryContext(ctx, sql, params...)
	if err != nil {
		return err
	}
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	defer rows.Close()
	columnLen := len(columns)
	c := make([]*V, columnLen)
	for i := range c {
		var a V
		c[i] = &a
	}
	args := slice.ToAnySlice(c)
	var m []map[string]V
	for rows.Next() {
		err = rows.Scan(args...)
		if err != nil {
			return
		}
		v := make(map[string]V)
		for i, data := range c {
			v[columns[i]] = *data
		}
		m = append(m, v)
	}
	*dest = m
	return
}

func GetToMap[V any](ctx context.Context, db *sqlx.DB, dest *map[string]V, sql string, params ...any) (err error) {
	rows := db.QueryRowxContext(ctx, sql, params...)
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	columnLen := len(columns)
	c := make([]*V, columnLen)
	for i := range c {
		var a V
		c[i] = &a
	}
	err = rows.Scan(slice.ToAnySlice(c)...)
	if err != nil {
		return
	}
	v := make(map[string]V)
	for i, data := range c {
		v[columns[i]] = *data
	}
	*dest = v
	return
}

func FormatSql(sql string, params ...any) string {
	for _, param := range params {
		switch param.(type) {
		case string:
			sql = strings.Replace(sql, "?", fmt.Sprintf("'%s'", param.(string)), 1)
		case int64:
			sql = strings.Replace(sql, "?", strconv.FormatInt(param.(int64), 10), 1)
		case int:
			sql = strings.Replace(sql, "?", strconv.Itoa(param.(int)), 1)
		case uint64:
			sql = strings.Replace(sql, "?", strconv.FormatUint(param.(uint64), 10), 1)
		case float64:
			sql = strings.Replace(sql, "?", fmt.Sprintf("%f", param.(float64)), 1)
		}
	}
	return str.CutSpecialDuplicate(sql, " ")
}
