package models

import (
	"fmt"
	"github/fthvgb1/wp-go/db"
	"github/fthvgb1/wp-go/helper"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Model interface {
	PrimaryKey() string
	Table() string
}

type ParseWhere interface {
	ParseWhere(in ...[]interface{}) (string, []interface{})
}

type SqlBuilder [][]string

func (w SqlBuilder) parseField(ss []string, s *strings.Builder) {
	if strings.Contains(ss[0], ".") && !strings.Contains(ss[0], "(") {
		s.WriteString("`")
		sx := strings.Split(ss[0], ".")
		s.WriteString(sx[0])
		s.WriteString("`.`")
		s.WriteString(sx[1])
		s.WriteString("`")
	} else if !strings.Contains(ss[0], ".") && !strings.Contains(ss[0], "(") {
		s.WriteString("`")
		s.WriteString(ss[0])
		s.WriteString("`")
	} else {
		s.WriteString(ss[0])
	}
}

func (w SqlBuilder) parseIn(ss []string, s *strings.Builder, c *int, args *[]interface{}, in [][]interface{}) (t bool) {
	if ss[1] == "in" && len(in) > 0 {
		s.WriteString(" (")
		for _, p := range in[*c] {
			s.WriteString("?,")
			*args = append(*args, p)
		}
		sx := s.String()
		s.Reset()
		s.WriteString(strings.TrimRight(sx, ","))
		s.WriteString(")")
		*c++
		t = true
	}
	return t
}

func (w SqlBuilder) parseType(ss []string, s *strings.Builder, args *[]interface{}) {
	if len(ss) == 4 && ss[3] == "int" {
		i, _ := strconv.Atoi(ss[2])
		*args = append(*args, i)
	} else if len(ss) == 4 && ss[3] == "float" {
		i, _ := strconv.ParseFloat(ss[2], 64)
		*args = append(*args, i)
	} else {
		*args = append(*args, ss[2])
	}
}

func (w SqlBuilder) ParseWhere(in ...[]interface{}) (string, []interface{}) {
	var s strings.Builder
	args := make([]interface{}, 0, len(w))
	c := 0
	for _, ss := range w {
		if len(ss) == 2 {
			w.parseField(ss, &s)
			s.WriteString("=? and ")
			args = append(args, ss[1])
		} else if len(ss) >= 3 && len(ss) < 5 {
			w.parseField(ss, &s)
			s.WriteString(ss[1])
			if w.parseIn(ss, &s, &c, &args, in) {
				s.WriteString(" and ")
				continue
			}
			s.WriteString(" ? and ")
			w.parseType(ss, &s, &args)
		} else if len(ss) >= 5 && len(ss)%5 == 0 {
			j := len(ss) / 5
			fl := false
			for i := 0; i < j; i++ {
				start := i * 5
				end := start + 5
				if ss[start] == "or" {
					st := s.String()
					if strings.Contains(st, "and ") {
						st = strings.TrimRight(st, "and ")
						s.Reset()
						s.WriteString(st)
						s.WriteString(" or ")
					}
					if i == 0 {
						s.WriteString("( ")
						fl = true
					}

					w.parseField(ss[start+1:end], &s)
					if w.parseIn(ss[start+1:end], &s, &c, &args, in) {
						s.WriteString(" and ")
						continue
					}
					s.WriteString(ss[start+2])
					s.WriteString(" ? and ")
					w.parseType(ss[start+1:end], &s, &args)
				} else {
					w.parseField(ss[start+1:end], &s)
					if w.parseIn(ss[start+1:end], &s, &c, &args, in) {
						s.WriteString(" and ")
						continue
					}
					s.WriteString(ss[start+2])
					s.WriteString(" ? and ")
					w.parseType(ss[start+1:start+4], &s, &args)

				}
				if i == j-1 && fl {
					st := s.String()
					st = strings.TrimRight(st, "and ")
					s.Reset()
					s.WriteString(st)
					s.WriteString(") and ")
				}
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
			s.WriteString(" ")
			s.WriteString(ss[0])
			s.WriteString(" ")
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

func SimplePagination[T Model](where ParseWhere, fields, group string, page, pageSize int, order SqlBuilder, join SqlBuilder, in ...[]interface{}) (r []T, total int, err error) {
	var rr T
	w, args := where.ParseWhere(in...)
	n := struct {
		N int `db:"n" json:"n"`
	}{}
	groupBy := ""
	if group != "" {
		g := strings.Builder{}
		g.WriteString(" group by ")
		g.WriteString(group)
		groupBy = g.String()
	}
	j := join.parseJoin()
	if group == "" {
		tpx := "select count(*) n from %s %s %s limit 1"
		sq := fmt.Sprintf(tpx, rr.Table(), j, w)
		err = db.Db.Get(&n, sq, args...)
	} else {
		tpx := "select count(*) n from (select %s from %s %s %s %s ) %s"
		rand.Seed(int64(time.Now().Nanosecond()))
		sq := fmt.Sprintf(tpx, group, rr.Table(), j, w, groupBy, fmt.Sprintf("table%d", rand.Int()))
		err = db.Db.Get(&n, sq, args...)
	}

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
	tp := "select %s from %s %s %s %s %s limit %d,%d"
	sql := fmt.Sprintf(tp, fields, rr.Table(), j, w, groupBy, order.parseOrderBy(), offset, pageSize)
	err = db.Db.Select(&r, sql, args...)
	if err != nil {
		return
	}
	return
}

func FindOneById[T Model](id int) (T, error) {
	var r T
	sql := fmt.Sprintf("select * from `%s` where `%s`=?", r.Table(), r.PrimaryKey())
	err := db.Db.Get(&r, sql, id)
	if err != nil {
		return r, err
	}
	return r, nil
}

func FirstOne[T Model](where ParseWhere, fields string, in ...[]interface{}) (T, error) {
	var r T
	w, args := where.ParseWhere(in...)
	tp := "select %s from %s %s"
	sql := fmt.Sprintf(tp, fields, r.Table(), w)
	err := db.Db.Get(&r, sql, args...)
	if err != nil {
		return r, err
	}
	return r, nil
}

func LastOne[T Model](where ParseWhere, fields string, in ...[]interface{}) (T, error) {
	var r T
	w, args := where.ParseWhere(in...)
	tp := "select %s from %s %s order by %s desc limit 1"
	sql := fmt.Sprintf(tp, fields, r.Table(), w, r.PrimaryKey())
	err := db.Db.Get(&r, sql, args...)
	if err != nil {
		return r, err
	}
	return r, nil
}

func SimpleFind[T Model](where ParseWhere, fields string, in ...[]interface{}) ([]T, error) {
	var r []T
	var rr T
	w, args := where.ParseWhere(in...)
	tp := "select %s from %s %s"
	sql := fmt.Sprintf(tp, fields, rr.Table(), w)
	err := db.Db.Select(&r, sql, args...)
	if err != nil {
		return r, err
	}
	return r, nil
}

func Select[T Model](sql string, params ...interface{}) ([]T, error) {
	var r []T
	var rr T
	sql = strings.Replace(sql, "{table}", rr.Table(), -1)
	err := db.Db.Select(&r, sql, params...)
	if err != nil {
		return r, err
	}
	return r, nil
}

func Find[T Model](where ParseWhere, fields, group string, order SqlBuilder, join SqlBuilder, limit int, in ...[]interface{}) (r []T, err error) {
	var rr T
	w := ""
	var args []interface{}
	if where != nil {
		w, args = where.ParseWhere(in...)
	}

	j := join.parseJoin()
	groupBy := ""
	if group != "" {
		g := strings.Builder{}
		g.WriteString(" group by ")
		g.WriteString(group)
		groupBy = g.String()
	}
	tp := "select %s from %s %s %s %s %s %s"
	l := ""
	if limit > 0 {
		l = fmt.Sprintf(" limit %d", limit)
	}
	sql := fmt.Sprintf(tp, fields, rr.Table(), j, w, groupBy, order.parseOrderBy(), l)
	err = db.Db.Select(&r, sql, args...)
	return
}

func Get[T Model](sql string, params ...interface{}) (r T, err error) {
	sql = strings.Replace(sql, "{table}", r.Table(), -1)
	err = db.Db.Get(&r, sql, params...)
	return
}
