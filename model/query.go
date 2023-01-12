package model

import (
	"context"
	"fmt"
	"github/fthvgb1/wp-go/helper"
	"math/rand"
	"strings"
	"time"
)

func SimplePagination[T Model](ctx context.Context, where ParseWhere, fields, group string, page, pageSize int, order SqlBuilder, join SqlBuilder, having SqlBuilder, in ...[]any) (r []T, total int, err error) {
	var rr T
	var w string
	var args []any
	if where != nil {
		w, args, err = where.ParseWhere(&in)
		if err != nil {
			return r, total, err
		}
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
		err = globalBb.Get(ctx, &n, sq, args...)
	} else {
		tpx := "select count(*) n from (select %s from %s %s %s %s %s ) %s"
		rand.Seed(int64(time.Now().Nanosecond()))
		sq := fmt.Sprintf(tpx, group, rr.Table(), j, w, groupBy, h, fmt.Sprintf("table%d", rand.Int()))
		err = globalBb.Get(ctx, &n, sq, args...)
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
	err = globalBb.Select(ctx, &r, sql, args...)
	if err != nil {
		return
	}
	return
}

func FindOneById[T Model, I helper.IntNumber](ctx context.Context, id I) (T, error) {
	var r T
	sql := fmt.Sprintf("select * from `%s` where `%s`=?", r.Table(), r.PrimaryKey())
	err := globalBb.Get(ctx, &r, sql, id)
	if err != nil {
		return r, err
	}
	return r, nil
}

func FirstOne[T Model](ctx context.Context, where ParseWhere, fields string, order SqlBuilder, in ...[]any) (T, error) {
	var r T
	var w string
	var args []any
	var err error
	if where != nil {
		w, args, err = where.ParseWhere(&in)
		if err != nil {
			return r, err
		}
	}
	tp := "select %s from %s %s %s"
	sql := fmt.Sprintf(tp, fields, r.Table(), w, order.parseOrderBy())
	err = globalBb.Get(ctx, &r, sql, args...)
	if err != nil {
		return r, err
	}
	return r, nil
}

func LastOne[T Model](ctx context.Context, where ParseWhere, fields string, in ...[]any) (T, error) {
	var r T
	var w string
	var args []any
	var err error
	if where != nil {
		w, args, err = where.ParseWhere(&in)
		if err != nil {
			return r, err
		}
	}
	tp := "select %s from %s %s order by %s desc limit 1"
	sql := fmt.Sprintf(tp, fields, r.Table(), w, r.PrimaryKey())
	err = globalBb.Get(ctx, &r, sql, args...)
	if err != nil {
		return r, err
	}
	return r, nil
}

func SimpleFind[T Model](ctx context.Context, where ParseWhere, fields string, in ...[]any) ([]T, error) {
	var r []T
	var rr T
	var err error
	var args []any
	var w string
	if where != nil {
		w, args, err = where.ParseWhere(&in)
		if err != nil {
			return r, err
		}
	}
	tp := "select %s from %s %s"
	sql := fmt.Sprintf(tp, fields, rr.Table(), w)
	err = globalBb.Select(ctx, &r, sql, args...)
	if err != nil {
		return r, err
	}
	return r, nil
}

func Select[T Model](ctx context.Context, sql string, params ...any) ([]T, error) {
	var r []T
	var rr T
	sql = strings.Replace(sql, "{table}", rr.Table(), -1)
	err := globalBb.Select(ctx, &r, sql, params...)
	if err != nil {
		return r, err
	}
	return r, nil
}

func Find[T Model](ctx context.Context, where ParseWhere, fields, group string, order SqlBuilder, join SqlBuilder, having SqlBuilder, limit int, in ...[]any) (r []T, err error) {
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
	err = globalBb.Select(ctx, &r, sql, args...)
	return
}

func Get[T Model](ctx context.Context, sql string, params ...any) (r T, err error) {
	sql = strings.Replace(sql, "{table}", r.Table(), -1)
	err = globalBb.Get(ctx, &r, sql, params...)
	return
}
