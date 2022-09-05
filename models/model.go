package models

import (
	"fmt"
	"github/fthvgb1/wp-go/db"
	"github/fthvgb1/wp-go/helper"
	"strconv"
	"strings"
)

type model interface {
	PrimaryKey() string
	Table() string
}

type SqlBuilder [][]string

func (w SqlBuilder) parseWhere(in ...[]interface{}) (string, []interface{}) {
	var s strings.Builder
	args := make([]interface{}, 0, len(w))
	c := 0
	for _, ss := range w {
		if len(ss) == 2 {
			s.WriteString("`")
			if strings.Contains(ss[0], ".") {
				sx := strings.Split(ss[0], ".")
				s.WriteString(sx[0])
				s.WriteString("`.`")
				s.WriteString(sx[1])
			} else {
				s.WriteString(ss[0])
			}
			s.WriteString("`=? and ")
			args = append(args, ss[1])
		}
		if len(ss) >= 3 {
			s.WriteString("`")
			if strings.Contains(ss[0], ".") {
				sx := strings.Split(ss[0], ".")
				s.WriteString(sx[0])
				s.WriteString("`.`")
				s.WriteString(sx[1])
			} else {
				s.WriteString(ss[0])
			}
			s.WriteString("`")
			s.WriteString(ss[1])
			if ss[1] == "in" && len(in) > 0 {
				s.WriteString(" (")
				for _, p := range in[c] {
					s.WriteString("?,")
					args = append(args, p)
				}
				sx := s.String()
				s.Reset()
				s.WriteString(strings.TrimRight(sx, ","))
				s.WriteString(")")
				c++
				s.WriteString(" and ")
				continue
			}
			s.WriteString(" ? and ")
			if len(ss) == 4 && ss[3] == "int" {
				i, _ := strconv.Atoi(ss[2])
				args = append(args, i)
			} else if len(ss) == 4 && ss[3] == "float" {
				i, _ := strconv.ParseFloat(ss[2], 64)
				args = append(args, i)
			} else {
				args = append(args, ss[2])
			}
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
func (w SqlBuilder) parseJoin() string {
	s := strings.Builder{}
	for _, ss := range w {
		l := len(ss)
		for j := 0; j < l; j++ {
			s.WriteString(" ")
			if (l == 4 && j == 3) || (l == 3 && j == 2) {
				s.WriteString("on ")
			}
			s.WriteString(ss[j])
			s.WriteString(" ")
		}

	}
	return s.String()
}

func SimplePagination[T model](where SqlBuilder, fields string, page, pageSize int, order SqlBuilder, join SqlBuilder, in ...[]interface{}) (r []T, total int, err error) {
	var rr T
	w, args := where.parseWhere(in...)
	n := struct {
		N int `db:"n" json:"n"`
	}{}
	j := join.parseJoin()
	tpx := "select count(*) n from %s %s %s limit 1"
	sq := fmt.Sprintf(tpx, rr.Table(), j, w)
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
	tp := "select %s from %s %s %s %s limit %d,%d"
	sql := fmt.Sprintf(tp, fields, rr.Table(), j, w, order.parseOrderBy(), offset, pageSize)
	err = db.Db.Select(&r, sql, args...)
	if err != nil {
		return
	}
	return
}

func FindOneById[T model](id int) (T, error) {
	var r T
	sql := fmt.Sprintf("select * from `%s` where `%s`=?", r.Table(), r.PrimaryKey())
	err := db.Db.Get(&r, sql, id)
	if err != nil {
		return r, err
	}
	return r, nil
}

func FirstOne[T model](where SqlBuilder, fields string, in ...[]interface{}) (T, error) {
	var r T
	w, args := where.parseWhere(in...)
	tp := "select %s from %s %s"
	sql := fmt.Sprintf(tp, fields, r.Table(), w)
	err := db.Db.Get(&r, sql, args...)
	if err != nil {
		return r, err
	}
	return r, nil
}

func LastOne[T model](where SqlBuilder, fields string, in ...[]interface{}) (T, error) {
	var r T
	w, args := where.parseWhere(in...)
	tp := "select %s from %s %s order by %s desc limit 1"
	sql := fmt.Sprintf(tp, fields, r.Table(), w, r.PrimaryKey())
	err := db.Db.Get(&r, sql, args...)
	if err != nil {
		return r, err
	}
	return r, nil
}

func SimpleFind[T model](where SqlBuilder, fields string, in ...[]interface{}) ([]T, error) {
	var r []T
	var rr T
	w, args := where.parseWhere(in...)
	tp := "select %s from %s %s"
	sql := fmt.Sprintf(tp, fields, rr.Table(), w)
	err := db.Db.Select(&r, sql, args...)
	if err != nil {
		return r, err
	}
	return r, nil
}

func Select[T model](sql string, params ...interface{}) ([]T, error) {
	var r []T
	var rr T
	sql = strings.Replace(sql, "{table}", rr.Table(), -1)
	err := db.Db.Select(&r, sql, params...)
	if err != nil {
		return r, err
	}
	return r, nil
}

func Find[T model](where SqlBuilder, fields string, order SqlBuilder, join SqlBuilder, limit int, in ...[]interface{}) (r []T, err error) {
	var rr T
	w, args := where.parseWhere(in...)
	j := join.parseJoin()
	tp := "select %s from %s %s %s %s %s"
	l := ""
	if limit > 0 {
		l = fmt.Sprintf(" limit %d", limit)
	}
	sql := fmt.Sprintf(tp, fields, rr.Table(), j, w, order.parseOrderBy(), l)
	err = db.Db.Select(&r, sql, args...)
	return
}

func Get[T model](sql string, params ...interface{}) (r T, err error) {
	sql = strings.Replace(sql, "{table}", r.Table(), -1)
	err = db.Db.Get(&r, sql, params...)
	return
}
