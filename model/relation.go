package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/slice"
	"golang.org/x/exp/constraints"
	"strings"
)

func setTable[T Model](q *QueryCondition) {
	if q.From == "" {
		q.From = Table[T]()
	}
}

type Relationship struct {
	RelationType string
	Table        string
	ForeignKey   string
	Local        string
	On           string
}

func Relation(isMultiple bool, db dbQuery, ctx context.Context, r any, q *QueryCondition) ([]func(), []func() error) {
	var fn []func()
	var fns []func() error
	for _, f := range q.RelationFn {
		getVal, isJoin, qq, ff := f()
		idFn, assignment, rr, rrs, ship := ff()
		if isJoin {
			fn = append(fn, func() {
				tables := strings.Split(ship.Table, " ")
				from := strings.Split(q.From, " ")
				on := ""
				if ship.On != "" {
					on = fmt.Sprintf("and %s", on)
				}
				qq := helper.GetContextVal(ctx, "ancestorsQueryCondition", q)
				qq.Join = append(qq.Join, []string{
					"left join", ship.Table, fmt.Sprintf("%s.%s=%s.%s %s", tables[len(tables)-1], ship.ForeignKey, from[len(from)-1], ship.Local, on)})
			})
		}
		if !getVal {
			continue
		}
		fns = append(fns, func() error {
			ids := idFn(r)
			if len(ids) < 1 {
				return nil
			}
			var err error
			{
				if qq == nil {
					qq = &QueryCondition{
						Fields: "*",
					}
				}
				var w any = qq.Where
				if w == nil {
					w = SqlBuilder{}
				}
				ww, ok := w.(SqlBuilder)
				if ok {
					ww = append(ww, SqlBuilder{{
						ship.ForeignKey, "in", "",
					}}...)
					qq.In = [][]any{ids}
					qq.Where = ww
				}
				if qq.From == "" {
					qq.From = ship.Table
				}
			}
			err = parseRelation(isMultiple || ship.RelationType == "hasMany", db, ctx, helper.Or(isMultiple, rrs, rr), qq)
			if err != nil && err != sql.ErrNoRows {
				return err
			}
			assignment(r, helper.Or(isMultiple, rrs, rr))

			return err
		})
	}
	return fn, fns
}

func GetWithID[T, V any](fn func(*T) V) func(any) []any {
	return func(a any) []any {
		one, ok := a.(*T)
		if ok {
			return []any{fn(one)}
		}
		return slice.ToAnySlice(slice.Unique(slice.Map(*a.(*[]T), func(t T) any {
			return fn(&t)
		})))
	}
}

func SetHasOne[T, V any, K comparable](fn func(*T, *V), idFn func(*T) K, iddFn func(*V) K) func(any, any) {
	return func(m, v any) {
		one, ok := m.(*T)
		if ok {
			fn(one, v.(*V))
			return
		}
		r := m.(*[]T)
		vv := v.(*[]V)
		mm := slice.SimpleToMap(*vv, func(v V) K {
			return iddFn(&v)
		})
		for i := 0; i < len(*r); i++ {
			val := &(*r)[i]
			id := idFn(val)
			v, ok := mm[id]
			if ok {
				fn(val, &v)
			}
		}
	}
}

func SetHasMany[T, V any, K comparable](fn func(*T, *[]V), idFn func(*T) K, iddFn func(*V) K) func(any, any) {
	return func(m, v any) {
		one, ok := m.(*T)
		if ok {
			fn(one, v.(*[]V))
			return
		}
		r := m.(*[]T)
		vv := v.(*[]V)
		mm := slice.GroupBy(*vv, func(t V) (K, V) {
			return iddFn(&t), t
		})
		for i := 0; i < len(*r); i++ {
			val := &(*r)[i]
			id := idFn(val)
			v, ok := mm[id]
			if ok {
				fn(val, &v)
			}
		}
	}
}

func RelationHasOne[M, P any, I constraints.Integer | uint64](fId func(*M) I, pId func(*P) I, setVal func(*M, *P), r Relationship) RelationFn {
	return func() (func(any) []any, func(any, any), any, any, Relationship) {
		var s P
		var ss []P
		return GetWithID(fId), SetHasOne(setVal, fId, pId), &s, &ss, r
	}
}
func RelationHasMany[M, P any, I constraints.Integer | uint64](mId func(*M) I, pId func(*P) I, setVal func(*M, *[]P), r Relationship) RelationFn {
	return func() (func(any) []any, func(any, any), any, any, Relationship) {
		var ss []P
		return GetWithID(mId), SetHasMany(setVal, mId, pId), &ss, &ss, r
	}
}
