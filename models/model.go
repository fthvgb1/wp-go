package models

import (
	"fmt"
	"github/fthvgb1/wp-go/db"
	"strings"
)

type mod interface {
	PrimaryKey() string
	Table() string
}

type model[T mod] struct {
}

type SimpleWhere [][]string

func (w SimpleWhere) parseWhere() (string, []interface{}) {
	var s strings.Builder
	args := make([]interface{}, 0, len(w))
	for _, ss := range w {
		if len(ss) == 2 {
			s.WriteString("`")
			s.WriteString(ss[0])
			s.WriteString("`=? and ")
			args = append(args, ss[1])
		}
		if len(ss) == 3 {
			s.WriteString("`")
			s.WriteString(ss[0])
			s.WriteString("`")
			s.WriteString(ss[1])
			s.WriteString("? and ")
			args = append(args, ss[2])
		}
	}
	return strings.TrimRight(s.String(), "and "), args
}

func (m model[T]) FindOneById(id int) (T, error) {
	var r T
	sql := fmt.Sprintf("select * from `%s` where `%s`=?", r.Table(), r.PrimaryKey())
	err := db.Db.Get(&r, sql, id)
	if err != nil {
		return r, err
	}
	return r, nil
}

func (m model[T]) FirstOne(where SimpleWhere, fields string) (T, error) {
	var r T
	w, args := where.parseWhere()
	tp := "select %s from %s where %s"
	sql := fmt.Sprintf(tp, fields, r.Table(), w)
	err := db.Db.Get(&r, sql, args...)
	if err != nil {
		return r, err
	}
	return r, nil
}

func (m model[T]) FindMany(where SimpleWhere, fields string) ([]T, error) {
	var r []T
	var rr T
	w, args := where.parseWhere()
	tp := "select %s from %s where %s"
	sql := fmt.Sprintf(tp, fields, rr.Table(), w)
	err := db.Db.Select(&r, sql, args...)
	if err != nil {
		return r, err
	}
	return r, nil
}

func (m model[T]) Get(sql string, params ...interface{}) (T, error) {
	var r T
	sql = strings.Replace(sql, "%table%", r.Table(), -1)
	err := db.Db.Get(&r, sql, params...)
	if err != nil {
		return r, err
	}
	return r, nil
}

func (m model[T]) Select(sql string, params ...interface{}) ([]T, error) {
	var r []T
	var rr T
	sql = strings.Replace(sql, "%table%", rr.Table(), -1)
	err := db.Db.Select(&r, sql, params...)
	if err != nil {
		return r, err
	}
	return r, nil
}
