package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/slice"
	"strings"
)

// Finds  比 Find 多一个offset
//
// Conditions 中可用 Where Fields Group Having Join Order Offset Limit In 函数
func Finds[T Model](ctx context.Context, q *QueryCondition) (r []T, err error) {
	r, err = finds[T](globalBb, ctx, q)
	return
}

// FindFromDB  同 Finds 使用指定 db 查询
//
// Conditions 中可用 Where Fields Group Having Join Order Offset Limit In 函数
func FindFromDB[T Model](db dbQuery, ctx context.Context, q *QueryCondition) (r []T, err error) {
	r, err = finds[T](db, ctx, q)
	return
}

func finds[T Model](db dbQuery, ctx context.Context, q *QueryCondition) (r []T, err error) {
	sq, args, err := FindRawSql[T](q)
	if err != nil {
		return
	}
	err = db.Select(ctx, &r, sq, args...)
	return
}

func chunkFind[T Model](db dbQuery, ctx context.Context, perLimit int, q *QueryCondition) (r []T, err error) {
	i := 1
	var rr []T
	var total int
	var offset int
	for {
		if 1 == i {
			rr, total, err = pagination[T](db, ctx, q.where, q.fields, q.group, i, perLimit, q.order, q.join, q.having, q.in...)
		} else {
			q.offset = offset
			q.limit = perLimit
			rr, err = finds[T](db, ctx, q)
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

// ChunkFind 分片查询并直接返回所有结果
//
// Conditions 中可用 Where Fields Group Having Join Order Limit In 函数
func ChunkFind[T Model](ctx context.Context, perLimit int, q *QueryCondition) (r []T, err error) {
	r, err = chunkFind[T](globalBb, ctx, perLimit, q)
	return
}

// ChunkFindFromDB 同 ChunkFind
//
// Conditions 中可用 Where Fields Group Having Join Order Limit In 函数
func ChunkFindFromDB[T Model](db dbQuery, ctx context.Context, perLimit int, q *QueryCondition) (r []T, err error) {
	r, err = chunkFind[T](db, ctx, perLimit, q)
	return
}

// Chunk 分片查询并函数过虑返回新类型的切片
//
// Conditions 中可用 Where Fields Group Having Join Order Limit In 函数
func Chunk[T Model, R any](ctx context.Context, perLimit int, fn func(rows T) (R, bool), q *QueryCondition) (r []R, err error) {
	r, err = chunk(globalBb, ctx, perLimit, fn, q)
	return
}

// ChunkFromDB 同 Chunk
//
// Conditions 中可用 Where Fields Group Having Join Order Limit In 函数
func ChunkFromDB[T Model, R any](db dbQuery, ctx context.Context, perLimit int, fn func(rows T) (R, bool), q *QueryCondition) (r []R, err error) {
	r, err = chunk(db, ctx, perLimit, fn, q)
	return
}

func chunk[T Model, R any](db dbQuery, ctx context.Context, perLimit int, fn func(rows T) (R, bool), q *QueryCondition) (r []R, err error) {
	i := 1
	var rr []T
	var count int
	var total int
	var offset int
	for {
		if 1 == i {
			rr, total, err = pagination[T](db, ctx, q.where, q.fields, q.group, i, perLimit, q.order, q.join, q.having, q.in...)
		} else {
			q.offset = offset
			q.limit = perLimit
			rr, err = finds[T](db, ctx, q)
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

// Pagination 同 SimplePagination
//
// Condition 中可使用 Where Fields Group Having Join Order Page Limit In 函数
func Pagination[T Model](ctx context.Context, q *QueryCondition) ([]T, int, error) {
	return SimplePagination[T](ctx, q.where, q.fields, q.group, q.page, q.limit, q.order, q.join, q.having, q.in...)
}

// PaginationFromDB 同 Pagination 方便多个db使用
//
// Condition 中可使用 Where Fields Group Having Join Order Page Limit In 函数
func PaginationFromDB[T Model](db dbQuery, ctx context.Context, q *QueryCondition) ([]T, int, error) {
	return pagination[T](db, ctx, q.where, q.fields, q.group, q.page, q.limit, q.order, q.join, q.having, q.in...)
}

func Column[V Model, T any](ctx context.Context, fn func(V) (T, bool), q *QueryCondition) ([]T, error) {
	return column[V, T](globalBb, ctx, fn, q)
}
func ColumnFromDB[V Model, T any](db dbQuery, ctx context.Context, fn func(V) (T, bool), q *QueryCondition) (r []T, err error) {
	return column[V, T](db, ctx, fn, q)
}

func column[V Model, T any](db dbQuery, ctx context.Context, fn func(V) (T, bool), q *QueryCondition) (r []T, err error) {
	res, err := finds[V](db, ctx, q)
	if err != nil {
		return nil, err
	}
	r = slice.FilterAndMap(res, fn)
	return
}

func GetField[T Model, V any](ctx context.Context, field string, q *QueryCondition) (r V, err error) {
	r, err = getField[T, V](globalBb, ctx, field, q)
	return
}
func getField[T Model, V any](db dbQuery, ctx context.Context, field string, q *QueryCondition) (r V, err error) {
	res, err := getToAnyMap[T](globalBb, ctx, q)
	if err != nil {
		return
	}
	r, ok := maps.GetStrAnyVal[V](res, field)
	if !ok {
		err = errors.New("not exists")
	}
	return
}
func GetFieldFromDB[T Model, V any](db dbQuery, ctx context.Context, field string, q *QueryCondition) (r V, err error) {
	return getField[T, V](db, ctx, field, q)
}

func findToAnyMap[T Model](db dbQuery, ctx context.Context, q *QueryCondition) (r []map[string]any, err error) {
	rawSql, in, err := FindRawSql[T](q)
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, "toMap", true)
	err = db.Select(ctx, &r, rawSql, in...)
	return
}

func FindToAnyMap[T Model](ctx context.Context, q *QueryCondition) (r []map[string]any, err error) {
	r, err = findToAnyMap[T](globalBb, ctx, q)
	return
}

func FindToAnyMapFromDB[T Model](db dbQuery, ctx context.Context, q *QueryCondition) (r []map[string]any, err error) {
	r, err = findToAnyMap[T](db, ctx, q)
	return
}
func getToAnyMap[T Model](db dbQuery, ctx context.Context, q *QueryCondition) (r map[string]any, err error) {
	rawSql, in, err := FindRawSql[T](q)
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, "toMap", true)
	err = db.Get(ctx, &r, rawSql, in...)
	return
}

func FindRawSql[T Model](q *QueryCondition) (r string, args []any, err error) {
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
	r = fmt.Sprintf(tp, q.fields, rr.Table(), j, w, groupBy, h, q.order.parseOrderBy(), l)
	return
}
