package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	setTable[T](q)
	sq, args, err := BuildQuerySql(q)
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
			rr, total, err = pagination[T](db, ctx, q, 1, perLimit)
		} else {
			q.Offset = offset
			q.Limit = perLimit
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
			rr, total, err = pagination[T](db, ctx, q, 1, perLimit)
		} else {
			q.Offset = offset
			q.Limit = perLimit
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

// Pagination 同
//
// Condition 中可使用 Where Fields From Group Having Join Order Limit In 函数
func Pagination[T Model](ctx context.Context, q *QueryCondition, page, pageSize int) ([]T, int, error) {
	return pagination[T](globalBb, ctx, q, page, pageSize)
}

// PaginationFromDB 同 Pagination 方便多个db使用
//
// Condition 中可使用 Where Fields Group Having Join Order Limit In 函数
func PaginationFromDB[T Model](db dbQuery, ctx context.Context, q *QueryCondition, page, pageSize int) ([]T, int, error) {
	return pagination[T](db, ctx, q, page, pageSize)
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

func GetField[T Model](ctx context.Context, field string, q *QueryCondition) (r string, err error) {
	r, err = getField[T](globalBb, ctx, field, q)
	return
}
func getField[T Model](db dbQuery, ctx context.Context, field string, q *QueryCondition) (r string, err error) {
	if q.Fields == "" || q.Fields == "*" {
		q.Fields = field
	}
	res, err := getToStringMap[T](db, ctx, q)
	if err != nil {
		return
	}
	f := strings.Split(field, " ")
	r, ok := res[f[len(f)-1]]
	if !ok {
		err = errors.New("not exists")
	}
	return
}
func GetFieldFromDB[T Model](db dbQuery, ctx context.Context, field string, q *QueryCondition) (r string, err error) {
	return getField[T](db, ctx, field, q)
}

func getToStringMap[T Model](db dbQuery, ctx context.Context, q *QueryCondition) (r map[string]string, err error) {
	setTable[T](q)
	rawSql, in, err := BuildQuerySql(q)
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, "handle=>", "string")
	err = db.Get(ctx, &r, rawSql, in...)
	return
}
func GetToStringMap[T Model](ctx context.Context, q *QueryCondition) (r map[string]string, err error) {
	r, err = getToStringMap[T](globalBb, ctx, q)
	return
}

func findToStringMap[T Model](db dbQuery, ctx context.Context, q *QueryCondition) (r []map[string]string, err error) {
	setTable[T](q)
	rawSql, in, err := BuildQuerySql(q)
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, "handle=>", "string")
	err = db.Select(ctx, &r, rawSql, in...)
	return
}

func FindToStringMap[T Model](ctx context.Context, q *QueryCondition) (r []map[string]string, err error) {
	r, err = findToStringMap[T](globalBb, ctx, q)
	return
}

func FindToStringMapFromDB[T Model](db dbQuery, ctx context.Context, q *QueryCondition) (r []map[string]string, err error) {
	r, err = findToStringMap[T](db, ctx, q)
	return
}

func GetToStringMapFromDB[T Model](db dbQuery, ctx context.Context, q *QueryCondition) (r map[string]string, err error) {
	r, err = getToStringMap[T](db, ctx, q)
	return
}

func BuildQuerySql(q *QueryCondition) (r string, args []any, err error) {
	where := ""
	if q.Where != nil {
		where, args, err = q.Where.ParseWhere(&q.In)
		if err != nil {
			return
		}
	}
	having := ""
	if q.Having != nil {
		hh, arg, er := q.Having.ParseWhere(&q.In)
		if er != nil {
			err = er
			return
		}
		args = append(args, arg...)
		having = strings.Replace(hh, " where", " having", 1)
	}
	if len(args) == 0 && len(q.In) > 0 {
		for _, antes := range q.In {
			args = append(args, antes...)
		}
	}
	join := ""
	if q.Join != nil {
		join = q.Join.parseJoin()
	}
	groupBy := ""
	if q.Group != "" {
		g := strings.Builder{}
		g.WriteString(" group by ")
		g.WriteString(q.Group)
		groupBy = g.String()
	}
	tp := "select %s from %s %s %s %s %s %s %s"
	l := ""
	table := q.From
	if q.Limit > 0 {
		l = fmt.Sprintf(" limit %d", q.Limit)
	}
	if q.Offset > 0 {
		l = fmt.Sprintf(" %s offset %d", l, q.Offset)
	}
	order := ""
	if q.Order != nil {
		order = q.Order.parseOrderBy()
	}
	r = fmt.Sprintf(tp, q.Fields, table, join, where, groupBy, having, order, l)
	return
}

func findScanner[T Model](db dbQuery, ctx context.Context, fn func(T), q *QueryCondition) (err error) {
	setTable[T](q)
	s, args, err := BuildQuerySql(q)
	if err != nil {
		return
	}
	ctx = context.WithValue(ctx, "handle=>", "scanner")
	var v T
	ctx = context.WithValue(ctx, "fn", func(v any) {
		fn(*(v.(*T)))
	})
	err = db.Select(ctx, &v, s, args...)
	return
}

func FindScannerFromDB[T Model](db dbQuery, ctx context.Context, fn func(T), q *QueryCondition) error {
	return findScanner[T](db, ctx, fn, q)
}

func FindScanner[T Model](ctx context.Context, fn func(T), q *QueryCondition) error {
	return findScanner[T](globalBb, ctx, fn, q)
}

func Gets[T Model](ctx context.Context, q *QueryCondition) (T, error) {
	return gets[T](globalBb, ctx, q)
}
func GetsFromDB[T Model](db dbQuery, ctx context.Context, q *QueryCondition) (T, error) {
	return gets[T](db, ctx, q)
}

func gets[T Model](db dbQuery, ctx context.Context, q *QueryCondition) (r T, err error) {
	setTable[T](q)
	if len(q.Relation) < 1 {
		s, args, er := BuildQuerySql(q)
		if er != nil {
			err = er
			return
		}
		err = db.Get(ctx, &r, s, args...)
		return
	}
	err = parseRelation(false, db, ctx, &r, q)
	return
}

func parseRelation(isMultiple bool, db dbQuery, ctx context.Context, r any, q *QueryCondition) (err error) {
	fn, fns := Relation(db, ctx, r, q)
	for _, f := range fn {
		f()
	}
	s, args, err := BuildQuerySql(q)
	if err != nil {
		return
	}
	if isMultiple {
		err = db.Select(ctx, r, s, args...)
	} else {
		err = db.Get(ctx, r, s, args...)
	}

	if err != nil {
		return
	}
	for _, f := range fns {
		err = f()
		if err != nil {
			return
		}
	}
	return
}
