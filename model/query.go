package model

import (
	"context"
	"fmt"
	"golang.org/x/exp/constraints"
	"math/rand"
	"strings"
)

func pagination[T Model](db dbQuery, ctx context.Context, where ParseWhere, fields, group string, page, pageSize int, order SqlBuilder, join SqlBuilder, having SqlBuilder, in ...[]any) (r []T, total int, err error) {
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
		err = db.Get(ctx, &n, sq, args...)
	} else {
		tpx := "select count(*) n from (select %s from %s %s %s %s %s ) %s"
		sq := fmt.Sprintf(tpx, group, rr.Table(), j, w, groupBy, h, fmt.Sprintf("table%d", rand.Int()))
		err = db.Get(ctx, &n, sq, args...)
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
	sq := fmt.Sprintf(tp, fields, rr.Table(), j, w, groupBy, h, order.parseOrderBy(), offset, pageSize)
	err = db.Select(ctx, &r, sq, args...)
	if err != nil {
		return
	}
	return
}

func SimplePagination[T Model](ctx context.Context, where ParseWhere, fields, group string, page, pageSize int, order SqlBuilder, join SqlBuilder, having SqlBuilder, in ...[]any) (r []T, total int, err error) {
	r, total, err = pagination[T](globalBb, ctx, where, fields, group, page, pageSize, order, join, having, in...)
	return
}

func FindOneById[T Model, I constraints.Integer](ctx context.Context, id I) (T, error) {
	var r T
	sq := fmt.Sprintf("select * from `%s` where `%s`=?", r.Table(), r.PrimaryKey())
	err := globalBb.Get(ctx, &r, sq, id)
	if err != nil {
		return r, err
	}
	return r, nil
}

func FirstOne[T Model](ctx context.Context, where ParseWhere, fields string, order SqlBuilder, in ...[]any) (r T, err error) {
	s, args, err := BuildQuerySql[T](&QueryCondition{
		where:  where,
		fields: fields,
		order:  order,
		in:     in,
		limit:  1,
	})
	if err != nil {
		return
	}
	err = globalBb.Get(ctx, &r, s, args...)
	return
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
	sq := fmt.Sprintf(tp, fields, r.Table(), w, r.PrimaryKey())
	err = globalBb.Get(ctx, &r, sq, args...)
	if err != nil {
		return r, err
	}
	return r, nil
}

func SimpleFind[T Model](ctx context.Context, where ParseWhere, fields string, in ...[]any) (r []T, err error) {
	s, args, err := BuildQuerySql[T](&QueryCondition{
		where:  where,
		fields: fields,
		in:     in,
	})
	if err != nil {
		return
	}
	err = globalBb.Select(ctx, &r, s, args...)
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
	q := QueryCondition{
		where:  where,
		fields: fields,
		group:  group,
		order:  order,
		join:   join,
		having: having,
		limit:  limit,
		in:     in,
	}
	s, args, err := BuildQuerySql[T](&q)
	if err != nil {
		return
	}
	err = globalBb.Select(ctx, &r, s, args...)
	return
}

func Get[T Model](ctx context.Context, sql string, params ...any) (r T, err error) {
	sql = strings.Replace(sql, "{table}", r.Table(), -1)
	err = globalBb.Get(ctx, &r, sql, params...)
	return
}
