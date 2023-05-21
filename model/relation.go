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
	var beforeFn []func()
	var afterFn []func() error
	for _, f := range q.RelationFn {
		getVal, isJoin, qq, relationship := f()
		idFn, assignmentFn, rr, rrs, ship := relationship()
		if isJoin {
			beforeFn = append(beforeFn, func() {
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
		afterFn = append(afterFn, func() error {
			ids := idFn(r)
			if len(ids) < 1 {
				return nil
			}
			var err error
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
			err = ParseRelation(isMultiple || ship.RelationType == "hasMany", db, ctx, helper.Or(isMultiple, rrs, rr), qq)
			if err != nil {
				if err == sql.ErrNoRows {
					err = nil
				} else {
					return err
				}
			}
			assignmentFn(r, helper.Or(isMultiple, rrs, rr))
			return err
		})
	}
	return beforeFn, afterFn
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

// SetHasOne mIdFn is main , pIdFn is part
//
// eg: post has a user. mIdFn is post's userId, iddFn is user's id
func SetHasOne[T, V any, K comparable](assignmentFn func(*T, *V), mIdFn func(*T) K, pIdFn func(*V) K) func(any, any) {
	return func(m, p any) {
		one, ok := m.(*T)
		if ok {
			assignmentFn(one, p.(*V))
			return
		}
		mSlice := m.(*[]T)
		pSLice := p.(*[]V)
		mm := slice.SimpleToMap(*pSLice, func(v V) K {
			return pIdFn(&v)
		})
		for i := 0; i < len(*mSlice); i++ {
			m := &(*mSlice)[i]
			id := mIdFn(m)
			p, ok := mm[id]
			if ok {
				assignmentFn(m, &p)
			}
		}
	}
}

// SetHasMany
// eg: post has many comments,pIdFn is comment's postId, mIdFn is post's id
func SetHasMany[T, V any, K comparable](assignmentFn func(*T, *[]V), pIdFn func(*T) K, mIdFn func(*V) K) func(any, any) {
	return func(m, p any) {
		one, ok := m.(*T)
		if ok {
			assignmentFn(one, p.(*[]V))
			return
		}
		r := m.(*[]T)
		vv := p.(*[]V)
		mm := slice.GroupBy(*vv, func(t V) (K, V) {
			return mIdFn(&t), t
		})
		for i := 0; i < len(*r); i++ {
			m := &(*r)[i]
			id := pIdFn(m)
			p, ok := mm[id]
			if ok {
				assignmentFn(m, &p)
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
