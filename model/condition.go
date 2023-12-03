package model

import "context"

type QueryCondition struct {
	Where      ParseWhere
	From       string
	Fields     string
	Group      string
	Order      SqlBuilder
	Join       SqlBuilder
	Having     SqlBuilder
	Limit      int
	Offset     int
	TotalRow   int
	In         [][]any
	RelationFn []func() (bool, bool, *QueryCondition, RelationFn)
}

type RelationFn func() (func(any) []any, func(any, any), func(bool) any, Relationship)

func Conditions(fns ...Condition) *QueryCondition {
	r := &QueryCondition{}
	for _, fn := range fns {
		fn(r)
	}
	if r.Fields == "" {
		r.Fields = "*"
	}
	return r
}

type Condition func(c *QueryCondition)

func Where(where ParseWhere) Condition {
	return func(c *QueryCondition) {
		c.Where = where
	}
}
func Fields(fields string) Condition {
	return func(c *QueryCondition) {
		c.Fields = fields
	}
}

func From(from string) Condition {
	return func(c *QueryCondition) {
		c.From = from
	}
}

func Group(group string) Condition {
	return func(c *QueryCondition) {
		c.Group = group
	}
}

func Order[T ~[][]string](order T) Condition {
	return func(c *QueryCondition) {
		c.Order = SqlBuilder(order)
	}
}

func Join[T ~[][]string](join T) Condition {
	return func(c *QueryCondition) {
		c.Join = SqlBuilder(join)
	}
}

func Having[T ~[][]string](having T) Condition {
	return func(c *QueryCondition) {
		c.Having = SqlBuilder(having)
	}
}

func Limit(limit int) Condition {
	return func(c *QueryCondition) {
		c.Limit = limit
	}
}

// TotalRaw only effect on Pagination,when TotalRaw>0 ,will not query count
func TotalRaw(total int) Condition {
	return func(c *QueryCondition) {
		c.TotalRow = total
	}
}

func Offset(offset int) Condition {
	return func(c *QueryCondition) {
		c.Offset = offset
	}
}

func In(in ...[]any) Condition {
	return func(c *QueryCondition) {
		c.In = append(c.In, in...)
	}
}

func WithCtx(ctx *context.Context) Condition {
	return func(c *QueryCondition) {
		*ctx = context.WithValue(*ctx, "ancestorsQueryCondition", c)
	}
}

func WithFn(getVal, isJoin bool, q *QueryCondition, fn func() (func(any) []any, func(any, any), func(bool) any, Relationship)) Condition {
	return func(c *QueryCondition) {
		c.RelationFn = append(c.RelationFn, func() (bool, bool, *QueryCondition, RelationFn) {
			return getVal, isJoin, q, fn
		})
	}
}
