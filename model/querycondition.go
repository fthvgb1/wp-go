package model

import (
	"fmt"
	"strings"
)

func BuildQuerySql[T Model](q *QueryCondition) (r string, args []any, err error) {
	var rr T
	w := ""
	if q.where != nil {
		w, args, err = q.where.ParseWhere(&q.in)
		if err != nil {
			return
		}
	}
	h := ""
	if q.having != nil {
		hh, arg, er := q.having.ParseWhere(&q.in)
		if er != nil {
			err = er
			return
		}
		args = append(args, arg...)
		h = strings.Replace(hh, " where", " having", 1)
	}

	j := q.join.parseJoin()
	groupBy := ""
	if q.group != "" {
		g := strings.Builder{}
		g.WriteString(" group by ")
		g.WriteString(q.group)
		groupBy = g.String()
	}
	tp := "select %s from %s %s %s %s %s %s %s"
	l := ""

	if q.limit > 0 {
		l = fmt.Sprintf(" limit %d", q.limit)
	}
	if q.offset > 0 {
		l = fmt.Sprintf(" %s offset %d", l, q.offset)
	}
	table := rr.Table()
	if q.from != "" {
		table = q.from
	}
	r = fmt.Sprintf(tp, q.fields, table, j, w, groupBy, h, q.order.parseOrderBy(), l)
	return
}
