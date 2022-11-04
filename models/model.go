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

var _ ParseWhere = SqlBuilder{}

type Model interface {
	PrimaryKey() string
	Table() string
}

type ParseWhere interface {
	ParseWhere(*[][]any) (string, []any, error)
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

func (w SqlBuilder) parseIn(ss []string, s *strings.Builder, c *int, args *[]any, in *[][]any) (t bool) {
	if helper.IsContainInArr(ss[1], []string{"in", "not in"}) && len(*in) > 0 {
		s.WriteString(" (")
		for _, p := range (*in)[*c] {
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

func (w SqlBuilder) parseType(ss []string, args *[]any) error {
	if len(ss) == 4 && ss[3] == "int" {
		i, err := strconv.ParseInt(ss[2], 10, 64)
		if err != nil {
			return err
		}
		*args = append(*args, i)
	} else if len(ss) == 4 && ss[3] == "float" {
		i, err := strconv.ParseFloat(ss[2], 64)
		if err != nil {
			return err
		}
		*args = append(*args, i)
	} else {
		*args = append(*args, ss[2])
	}
	return nil
}

func (w SqlBuilder) ParseWhere(in *[][]any) (string, []any, error) {
	var s strings.Builder
	args := make([]any, 0, len(w))
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
			err := w.parseType(ss, &args)
			if err != nil {
				return "", nil, err
			}
		} else if len(ss) >= 5 && len(ss)%5 == 0 {
			j := len(ss) / 5
			for i := 0; i < j; i++ {
				start := i * 5
				end := start + 5
				st := s.String()
				if strings.Contains(st, "and ") && ss[start] == "or" {
					st = strings.TrimRight(st, "and ")
					s.Reset()
					s.WriteString(st)
					s.WriteString(fmt.Sprintf(" %s ", ss[start]))
				}
				if i == 0 {
					s.WriteString("( ")
				}
				w.parseField(ss[start+1:end], &s)
				s.WriteString(ss[start+2])
				if w.parseIn(ss[start+1:end], &s, &c, &args, in) {
					s.WriteString(" and ")
					continue
				}
				s.WriteString(" ? and ")
				err := w.parseType(ss[start+1:end], &args)
				if err != nil {
					return "", nil, err
				}
			}
			st := s.String()
			st = strings.TrimRight(st, "and ")
			s.Reset()
			s.WriteString(st)
			s.WriteString(") and ")
		}
	}
	ss := strings.TrimRight(s.String(), "and ")
	if ss != "" {
		s.Reset()
		s.WriteString(" where ")
		s.WriteString(ss)
		ss = s.String()
	}
	if len(*in) > c {
		*in = (*in)[c:]
	}
	return ss, args, nil
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

func SimplePagination[T Model](where ParseWhere, fields, group string, page, pageSize int, order SqlBuilder, join SqlBuilder, having SqlBuilder, in ...[]any) (r []T, total int, err error) {
	var rr T
	w, args, err := where.ParseWhere(&in)
	if err != nil {
		return r, total, err
	}
	h := ""
	if having != nil {
		hh, arg, err := having.ParseWhere(&in)
		if err != nil {
			return r, total, err
		}
		args = append(args, arg...)
		h = strings.Replace(hh, " where", " having", 1)
	}

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
	if having != nil {
		tm := map[string]struct{}{}
		for _, s := range strings.Split(group, ",") {
			tm[s] = struct{}{}
		}
		for _, ss := range having {
			if _, ok := tm[ss[0]]; !ok {
				group = fmt.Sprintf("%s,%s", group, ss[0])
			}
		}
		group = strings.Trim(group, ",")
	}
	j := join.parseJoin()
	if group == "" {
		tpx := "select count(*) n from %s %s %s limit 1"
		sq := fmt.Sprintf(tpx, rr.Table(), j, w)
		err = db.Db.Get(&n, sq, args...)
	} else {
		tpx := "select count(*) n from (select %s from %s %s %s %s %s ) %s"
		rand.Seed(int64(time.Now().Nanosecond()))
		sq := fmt.Sprintf(tpx, group, rr.Table(), j, w, groupBy, h, fmt.Sprintf("table%d", rand.Int()))
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
	tp := "select %s from %s %s %s %s %s %s limit %d,%d"
	sql := fmt.Sprintf(tp, fields, rr.Table(), j, w, groupBy, h, order.parseOrderBy(), offset, pageSize)
	err = db.Db.Select(&r, sql, args...)
	if err != nil {
		return
	}
	return
}

func FindOneById[T Model, I helper.IntNumber](id I) (T, error) {
	var r T
	sql := fmt.Sprintf("select * from `%s` where `%s`=?", r.Table(), r.PrimaryKey())
	err := db.Db.Get(&r, sql, id)
	if err != nil {
		return r, err
	}
	return r, nil
}

func FirstOne[T Model](where ParseWhere, fields string, order SqlBuilder, in ...[]any) (T, error) {
	var r T
	w, args, err := where.ParseWhere(&in)
	if err != nil {
		return r, err
	}
	tp := "select %s from %s %s %s"
	sql := fmt.Sprintf(tp, fields, r.Table(), w, order.parseOrderBy())
	err = db.Db.Get(&r, sql, args...)
	if err != nil {
		return r, err
	}
	return r, nil
}

func LastOne[T Model](where ParseWhere, fields string, in ...[]any) (T, error) {
	var r T
	w, args, err := where.ParseWhere(&in)
	if err != nil {
		return r, err
	}
	tp := "select %s from %s %s order by %s desc limit 1"
	sql := fmt.Sprintf(tp, fields, r.Table(), w, r.PrimaryKey())
	err = db.Db.Get(&r, sql, args...)
	if err != nil {
		return r, err
	}
	return r, nil
}

func SimpleFind[T Model](where ParseWhere, fields string, in ...[]any) ([]T, error) {
	var r []T
	var rr T
	w, args, err := where.ParseWhere(&in)
	if err != nil {
		return r, err
	}
	tp := "select %s from %s %s"
	sql := fmt.Sprintf(tp, fields, rr.Table(), w)
	err = db.Db.Select(&r, sql, args...)
	if err != nil {
		return r, err
	}
	return r, nil
}

func Select[T Model](sql string, params ...any) ([]T, error) {
	var r []T
	var rr T
	sql = strings.Replace(sql, "{table}", rr.Table(), -1)
	err := db.Db.Select(&r, sql, params...)
	if err != nil {
		return r, err
	}
	return r, nil
}

func Find[T Model](where ParseWhere, fields, group string, order SqlBuilder, join SqlBuilder, having SqlBuilder, limit int, in ...[]any) (r []T, err error) {
	var rr T
	w := ""
	var args []any
	if where != nil {
		w, args, err = where.ParseWhere(&in)
		if err != nil {
			return r, err
		}
	}
	h := ""
	if having != nil {
		hh, arg, err := having.ParseWhere(&in)
		if err != nil {
			return r, err
		}
		args = append(args, arg...)
		h = strings.Replace(hh, " where", " having", 1)
	}

	j := join.parseJoin()
	groupBy := ""
	if group != "" {
		g := strings.Builder{}
		g.WriteString(" group by ")
		g.WriteString(group)
		groupBy = g.String()
	}
	tp := "select %s from %s %s %s %s %s %s %s"
	l := ""
	if limit > 0 {
		l = fmt.Sprintf(" limit %d", limit)
	}
	sql := fmt.Sprintf(tp, fields, rr.Table(), j, w, groupBy, h, order.parseOrderBy(), l)
	err = db.Db.Select(&r, sql, args...)
	return
}

func Get[T Model](sql string, params ...any) (r T, err error) {
	sql = strings.Replace(sql, "{table}", r.Table(), -1)
	err = db.Db.Get(&r, sql, params...)
	return
}
