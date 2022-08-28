package models

import (
	"fmt"
	"github/fthvgb1/wp-go/db"
	"github/fthvgb1/wp-go/helper"
	"strings"
)

type mod interface {
	PrimaryKey() string
	Table() string
}

type model[T mod] struct {
}

func (m model[T]) PrimaryKey() string {
	return "id"
}

func (m model[T]) Table() string {
	return ""
}

type SqlBuilder [][]string

func (w SqlBuilder) parseWhere() (string, []interface{}) {
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
			s.WriteString(" ? and ")
			args = append(args, ss[2])
		}
	}
	ss := strings.TrimRight(s.String(), "and ")
	if ss != "" {
		s.Reset()
		s.WriteString(" where ")
		s.WriteString(ss)
		ss = s.String()
	}
	return ss, args
}

func (w SqlBuilder) parseOrderBy() string {
	s := strings.Builder{}
	for _, ss := range w {
		if len(ss) == 2 && ss[0] != "" && helper.IsContainInArr(ss[1], []string{"asc", "desc"}) {
			s.WriteString(" `")
			s.WriteString(ss[0])
			s.WriteString("` ")
			s.WriteString(ss[1])
			s.WriteString(",")
		}
	}
	ss := strings.TrimRight(s.String(), ",")
	if ss != "" {
		s.Reset()
		s.WriteString(" order by ")
		s.WriteString(ss)
		ss = s.String()
	}
	return ss
}

func (m model[T]) SimplePagination(where SqlBuilder, fields string, page, pageSize int, order SqlBuilder) (r []T, total int, err error) {
	var rr T
	w, args := where.parseWhere()
	n := struct {
		N int `db:"n" json:"n"`
	}{}
	tpx := "select count(*) n from %s %s limit 1"
	sq := fmt.Sprintf(tpx, rr.Table(), w)
	err = db.Db.Get(&n, sq, args...)
	if err != nil {
		return
	}
	if n.N == 0 {
		return
	}
	total = n.N
	offset := 0
	if page > 1 {
		offset = (page - 1) * pageSize
	}
	if offset >= total {
		return
	}
	tp := "select %s from %s %s %s limit %d,%d"
	sql := fmt.Sprintf(tp, fields, rr.Table(), w, order.parseOrderBy(), offset, pageSize)
	err = db.Db.Select(&r, sql, args...)
	if err != nil {
		return
	}
	return
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

func (m model[T]) FirstOne(where SqlBuilder, fields string) (T, error) {
	var r T
	w, args := where.parseWhere()
	tp := "select %s from %s %s"
	sql := fmt.Sprintf(tp, fields, r.Table(), w)
	err := db.Db.Get(&r, sql, args...)
	if err != nil {
		return r, err
	}
	return r, nil
}

func (m model[T]) LastOne(where SqlBuilder, fields string) (T, error) {
	var r T
	w, args := where.parseWhere()
	tp := "select %s from %s %s order by %s desc limit 1"
	sql := fmt.Sprintf(tp, fields, r.Table(), w, r.PrimaryKey())
	err := db.Db.Get(&r, sql, args...)
	if err != nil {
		return r, err
	}
	return r, nil
}

func (m model[T]) FindMany(where SqlBuilder, fields string) ([]T, error) {
	var r []T
	var rr T
	w, args := where.parseWhere()
	tp := "select %s from %s %s"
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
