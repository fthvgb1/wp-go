package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// Finds can use offset
//
// Conditions 中可用 Where Fields Group Having Join Order Offset Limit In 函数
func Finds[T Model](ctx context.Context, q *QueryCondition) (r []T, err error) {
	var rr T
	w := ""
	var args []any
	if q.where != nil {
		w, args, err = q.where.ParseWhere(&q.in)
		if err != nil {
			return r, err
		}
	}
	h := ""
	if q.having != nil {
		hh, arg, err := q.having.ParseWhere(&q.in)
		if err != nil {
			return r, err
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
	sq := fmt.Sprintf(tp, q.fields, rr.Table(), j, w, groupBy, h, q.order.parseOrderBy(), l)
	err = globalBb.Select(ctx, &r, sq, args...)
	return
}

// ChunkFind 分片查询并直接返回所有结果
//
// Conditions 中可用 Where Fields Group Having Join Order Limit In 函数
func ChunkFind[T Model](ctx context.Context, perLimit int, q *QueryCondition) (r []T, err error) {
	i := 1
	var rr []T
	var total int
	var offset int
	for {
		if 1 == i {
			rr, total, err = SimplePagination[T](ctx, q.where, q.fields, q.group, i, perLimit, q.order, q.join, q.having, q.in...)
		} else {
			rr, err = Finds[T](ctx, Conditions(
				Where(q.where),
				Fields(q.fields),
				Group(q.group),
				Having(q.having),
				Join(q.join),
				Order(q.order),
				Offset(offset),
				Limit(perLimit),
				In(q.in...),
			))
		}
		offset += perLimit
		if (err != nil && err != sql.ErrNoRows) || len(rr) < 1 {
			return
		}
		r = append(r, rr...)
		if len(r) >= total {
			break
		}
		i++
	}
	return
}

// Chunk 分片查询并函数过虑返回新类型的切片
//
// Conditions 中可用 Where Fields Group Having Join Order Limit In 函数
func Chunk[T Model, R any](ctx context.Context, perLimit int, fn func(rows T) (R, bool), q *QueryCondition) (r []R, err error) {
	i := 1
	var rr []T
	var count int
	var total int
	var offset int
	for {
		if 1 == i {
			rr, total, err = SimplePagination[T](ctx, q.where, q.fields, q.group, i, perLimit, q.order, q.join, q.having, q.in...)
		} else {
			rr, err = Finds[T](ctx, Conditions(
				Where(q.where),
				Fields(q.fields),
				Group(q.group),
				Having(q.having),
				Join(q.join),
				Order(q.order),
				Offset(offset),
				Limit(perLimit),
				In(q.in...),
			))
		}
		offset += perLimit
		if (err != nil && err != sql.ErrNoRows) || len(rr) < 1 {
			return
		}
		for _, t := range rr {
			v, ok := fn(t)
			if ok {
				r = append(r, v)
			}
		}
		count += len(rr)
		if count >= total {
			break
		}
		i++
	}
	return
}
