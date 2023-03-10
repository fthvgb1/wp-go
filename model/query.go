package model

import (
	"context"
	"github.com/fthvgb1/wp-go/helper/number"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"golang.org/x/exp/constraints"
	"math/rand"
	"strings"
)

type count[T Model] struct {
	t T
	N int `json:"n,omitempty" db:"n" gorm:"n"`
}

func (c count[T]) PrimaryKey() string {
	return c.t.PrimaryKey()
}

func (c count[T]) Table() string {
	return c.t.Table()
}

func pagination[T Model](db dbQuery, ctx context.Context, q QueryCondition, page, pageSize int) (r []T, total int, err error) {
	if page < 1 || pageSize < 1 {
		return
	}
	q.Limit = pageSize
	qx := QueryCondition{
		Where:  q.Where,
		Having: q.Having,
		Join:   q.Join,
		In:     q.In,
		Group:  q.Group,
		From:   q.From,
		Fields: "count(*) n",
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
			From:   qx.From,
			In:     qx.In,
			Fields: "count(*) n",
		}
	}
	n, err := gets[count[T]](db, ctx, qx)
	total = n.N
	if err != nil || total < 1 {
		return
	}
	offset := 0
	if page > 1 {
		offset = (page - 1) * q.Limit
	}
	if offset >= total {
		return
	}
	q.Offset = offset
	m := ctx.Value("handle=>toMap")
	if m == nil {
		r, err = finds[T](db, ctx, q)
		return
	}
	mm, ok := m.(*[]map[string]string)
	if ok {
		mx, er := findToStringMap[T](db, ctx, q)
		if er != nil {
			err = er
			return
		}
		*mm = mx
	}
	return
}

func paginationToMap[T Model](db dbQuery, ctx context.Context, q QueryCondition, page, pageSize int) (r []map[string]string, total int, err error) {
	ctx = context.WithValue(ctx, "handle=>toMap", &r)
	_, total, err = pagination[T](db, ctx, q, page, pageSize)
	return
}

func PaginationToMap[T Model](ctx context.Context, q QueryCondition, page, pageSize int) (r []map[string]string, total int, err error) {
	return paginationToMap[T](globalBb, ctx, q, page, pageSize)
}
func PaginationToMapFromDB[T Model](db dbQuery, ctx context.Context, q QueryCondition, page, pageSize int) (r []map[string]string, total int, err error) {
	return paginationToMap[T](db, ctx, q, page, pageSize)
}

func FindOneById[T Model, I constraints.Integer](ctx context.Context, id I) (T, error) {
	return gets[T](globalBb, ctx, QueryCondition{
		Fields: "*",
		Where: SqlBuilder{
			{PrimaryKey[T](), "=", number.ToString(id), "int"},
		},
	})
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
	return gets[T](globalBb, ctx, Conditions(
		Where(where),
		Fields(fields),
		In(in...),
		Order(SqlBuilder{{PrimaryKey[T](), "desc"}}),
		Limit(1),
	))
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

// Select ??????????????????T???????????????????????? {table}?????????
func Select[T Model](ctx context.Context, sql string, params ...any) ([]T, error) {
	var r []T
	sql = strings.Replace(sql, "{table}", Table[T](), -1)
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

// Get ???????????? {table}????????? T?????????
func Get[T Model](ctx context.Context, sql string, params ...any) (r T, err error) {
	sql = strings.Replace(sql, "{table}", r.Table(), -1)
	err = globalBb.Get(ctx, &r, sql, params...)
	return
}
