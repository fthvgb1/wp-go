package model

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
)

type SqlxQuery struct {
	sqlx *sqlx.DB
	UniversalDb
}

func NewSqlxQuery(sqlx *sqlx.DB, u UniversalDb) *SqlxQuery {

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
	v := ctx.Value("toMap")
	if v != nil {
		vv, ok := v.(bool)
		if ok && vv {
			d, ok := dest.(*[]map[string]any)
			if ok {
				return r.toMapSlice(d, sql, params...)
			}
		}
	}
	return r.sqlx.Select(dest, sql, params...)
}

func (r *SqlxQuery) Gets(ctx context.Context, dest any, sql string, params ...any) error {
	v := ctx.Value("toMap")
	if v != nil {
		vv, ok := v.(bool)
		if ok && vv {
			d, ok := dest.(*map[string]any)
			if ok {
				return r.toMap(d, sql, params...)
			}
		}
	}
	return r.sqlx.Get(dest, sql, params...)
}

func (r *SqlxQuery) toMap(dest *map[string]any, sql string, params ...any) (err error) {
	rows := r.sqlx.QueryRowx(sql, params...)
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	columnLen := len(columns)
	c := make([]any, columnLen)
	for i, _ := range c {
		var a any
		c[i] = &a
	}
	err = rows.Scan(c...)
	if err != nil {
		return
	}
	v := make(map[string]any)
	for i, data := range c {
		s, ok := data.(*any)
		if ok {
			ss, ok := (*s).([]uint8)
			if ok {
				data = string(ss)
			}
		}
		v[columns[i]] = data
	}
	*dest = v
	return
}

func (r *SqlxQuery) toMapSlice(dest *[]map[string]any, sql string, params ...any) (err error) {
	rows, err := r.sqlx.Query(sql, params...)
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	defer rows.Close()
	columnLen := len(columns)
	c := make([]any, columnLen)
	for i, _ := range c {
		var a any
		c[i] = &a
	}
	var m []map[string]any
	for rows.Next() {
		err = rows.Scan(c...)
		if err != nil {
			return
		}
		v := make(map[string]any)
		for i, data := range c {
			s, ok := data.(*any)
			if ok {
				ss, ok := (*s).([]uint8)
				if ok {
					data = string(ss)
				}
			}
			v[columns[i]] = data
		}
		m = append(m, v)
	}
	*dest = m
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
