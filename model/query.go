package model

import (
	"context"
	"fmt"
	"github.com/fthvgb1/wp-go/helper/number"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"golang.org/x/exp/constraints"
	"math/rand"
	"strings"
)

func pagination[T Model](db dbQuery, ctx context.Context, q QueryCondition) (r []T, total int, err error) {
	qx := QueryCondition{
		Where:  q.Where,
		Having: q.Having,
		Join:   q.Join,
		In:     q.In,
		Group:  q.Group,
		From:   q.From,
	}
	if q.Group != "" {
		qx.Fields = q.Fields
		sq, in, er := BuildQuerySql[T](qx)
		qx.In = [][]any{in}
		if er != nil {
			err = er
			return
		}
		qx.From = str.Join("( ", sq, " ) ", "table", number.ToString(rand.Int()))
		qx = QueryCondition{
			From: qx.From,
			In:   qx.In,
		}
	}

	n, err := GetField[T](ctx, "count(*)", qx)
	total = str.ToInt[int](n)
	if err != nil || total < 1 {
		return
	}
	offset := 0
	if q.Page > 1 {
		offset = (q.Page - 1) * q.Limit
	}
	if offset >= total {
		return
	}
	q.Offset = offset
	sq, args, err := BuildQuerySql[T](q)
	if err != nil {
		return
	}
	err = db.Select(ctx, &r, sq, args...)
	if err != nil {
		return
	}
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
	s, args, err := BuildQuerySql[T](QueryCondition{
		Where:  where,
		Fields: fields,
		Order:  order,
		In:     in,
		Limit:  1,
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
	s, args, err := BuildQuerySql[T](QueryCondition{
		Where:  where,
		Fields: fields,
		In:     in,
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
		Where:  where,
		Fields: fields,
		Group:  group,
		Order:  order,
		Join:   join,
		Having: having,
		Limit:  limit,
		In:     in,
	}
	s, args, err := BuildQuerySql[T](q)
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
