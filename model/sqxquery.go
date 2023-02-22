package model

import (
	"context"
	"fmt"
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
)

type SqlxQuery[T any] struct {
	sqlx *sqlx.DB
	UniversalDb[T]
}

func NewSqlxQuery[T any](sqlx *sqlx.DB, u UniversalDb[T]) *SqlxQuery[T] {

	s := &SqlxQuery[T]{sqlx: sqlx, UniversalDb: u}
	if u.selects == nil {
		s.UniversalDb.selects = s.Selects
	}
	if u.gets == nil {
		s.UniversalDb.gets = s.Gets
	}
	return s
}

func SetSelect[T any](db *SqlxQuery[T], fn QuerySelect[T]) {
	db.selects = fn
}
func SetGet[T any](db *SqlxQuery[T], fn QueryGet[T]) {
	db.gets = fn
}

func (s *SqlxQuery[T]) Selects(ctx context.Context, sql string, params ...any) (r []T, err error) {
	v := ctx.Value("handle=>")
	if v != nil {
		vv, ok := v.(string)
		if ok && vv != "" {
			switch vv {
			case "string":
				//return ToMapSlice(r.sqlx, dest.(*[]map[string]string), sql, params...)
			case "scanner":
				fn := ctx.Value("fn")
				return nil, Scanner[T](s.sqlx, sql, params...)(fn.(func(T)))
			}
		}
	}
	//var a T
	err = s.sqlx.Select(&r, sql, params...)
	return
}

func (s *SqlxQuery[T]) Gets(ctx context.Context, sql string, params ...any) (r T, err error) {
	v := ctx.Value("handle=>")
	if v != nil {
		vv, ok := v.(string)
		if ok && vv != "" {
			switch vv {
			case "string":
				//return GetToMap(r.sqlx, dest.(*map[string]string), sql, params...)
			}
		}
	}
	err = s.sqlx.Get(&r, sql, params...)
	return
}

func Scanner[T any](db *sqlx.DB, s string, params ...any) func(func(T)) error {
	var v T
	return func(fn func(T)) error {
		rows, err := db.Queryx(s, params...)
		if err != nil {
			return err
		}
		for rows.Next() {
			err = rows.StructScan(&v)
			if err != nil {
				return err
			}
			fn(v)
		}
		return nil
	}
}

func ToMapSlice[V any](db *sqlx.DB, dest *[]map[string]V, sql string, params ...any) (err error) {
	rows, err := db.Query(sql, params...)
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	defer rows.Close()
	columnLen := len(columns)
	c := make([]*V, columnLen)
	for i, _ := range c {
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

func GetToMap[V any](db *sqlx.DB, dest *map[string]V, sql string, params ...any) (err error) {
	rows := db.QueryRowx(sql, params...)
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	columnLen := len(columns)
	c := make([]*V, columnLen)
	for i, _ := range c {
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
	return sql
}
