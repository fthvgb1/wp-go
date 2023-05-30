package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"golang.org/x/exp/constraints"
	"strings"
)

func setTable[T Model](q *QueryCondition) {
	if q.From == "" {
		q.From = Table[T]()
	}
}

const (
	HasOne  = "hasOne"
	HasMany = "hasMany"
)

// Relationship join table
//
// # RelationType HasOne| HasMany
//
// eg: hasOne, post has a user. ForeignKey is user's id , Local is post's userId field
//
// eg: hasMany, post has many comments,ForeignKey is comment's postId field, Local is post's id field
//
// On is additional join on conditions
type Relationship struct {
	RelationType string
	Table        string
	ForeignKey   string
	Local        string
	On           string
	Middle       *Relationship
}

func parseBeforeJoin(qq *QueryCondition, ship Relationship) {
	var fromTable, foreignKey, local string
	if ship.Middle != nil {
		parseBeforeJoin(qq, *ship.Middle)
		local = ship.Local
		fromTable = ship.Middle.Table
	} else {
		fromTable = qq.From
	}
	foreignKey = ship.ForeignKey
	local = ship.Local
	tables := strings.Split(ship.Table, " ")
	from := strings.Split(fromTable, " ")
	on := ""
	if ship.On != "" {
		on = fmt.Sprintf("and %s", on)
	}
	qq.Join = append(qq.Join, []string{
		"left join", ship.Table,
		fmt.Sprintf("%s.%s=%s.%s %s",
			tables[len(tables)-1], foreignKey, from[len(from)-1], local, on,
		)})

}

func parseAfterJoin(fromTable string, ids [][]any, qq *QueryCondition, ship Relationship) bool {
	tables := strings.Split(ship.Middle.Table, " ")
	from := strings.Split(fromTable, " ")
	on := ""
	if ship.On != "" {
		on = fmt.Sprintf("and %s", on)
	}
	foreignKey := ship.ForeignKey
	local := ship.Local
	qq.Join = append(qq.Join, []string{
		"left join", ship.Middle.Table,
		fmt.Sprintf("%s.%s=%s.%s %s",
			tables[len(tables)-1], foreignKey, from[len(from)-1], local, on,
		),
	})
	if ship.Middle != nil && ship.Middle.Middle != nil {
		r := parseAfterJoin(tables[len(tables)-1], ids, qq, *ship.Middle)
		return ship.RelationType == HasMany || r
	} else {
		from := strings.Split(qq.From, " ")
		ww, ok := qq.Where.(SqlBuilder)
		if ok {
			ww = append(ww, []string{fmt.Sprintf("%s.%s",
				tables[len(tables)-1], ship.Middle.ForeignKey), "in", ""},
			)
			qq.Where = ww
		} else {
			aw, ok := helper.IsImplements[AndWhere](qq.Where)
			if ok {
				vv := aw.AndWhere(fmt.Sprintf("%s.%s",
					tables[len(tables)-1], ship.Middle.ForeignKey), "in", strings.Join(slice.DecompressBy(ids, func(t any) (string, bool) {
					return fmt.Sprintf("%v", t), true
				}), ","), "int")
				wa, ok := helper.IsImplements[ParseWhere](vv)
				if ok {
					qq.Where = wa
				}
			}
		}
		if qq.Fields == "" || qq.Fields == "*" {
			qq.Fields = str.Join(from[len(from)-1], ".", "*", ",", tables[len(tables)-1], ".", ship.Middle.ForeignKey)
		}
		qq.In = ids
		return ship.RelationType == HasMany || ship.Middle.RelationType == HasMany
	}
}

func Relation(isPlural bool, db dbQuery, ctx context.Context, r any, q *QueryCondition) ([]func(), []func() error) {
	var beforeFn []func()
	var afterFn []func() error
	qx := helper.GetContextVal(ctx, "ancestorsQueryCondition", q)

	for _, f := range q.RelationFn {
		getVal, isJoin, qq, relationship := f()
		idFn, assignmentFn, rr, rrs, ship := relationship()
		if isJoin {
			beforeFn = append(beforeFn, func() {
				parseBeforeJoin(qx, ship)
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
			if qq.From == "" {
				qq.From = ship.Table
			}
			var w any = qq.Where
			if w == nil {
				qq.Where = SqlBuilder{}
			}
			in := [][]any{ids}
			if ship.Middle != nil {
				isPlural = parseAfterJoin(qq.From, in, qq, ship)
			} else {
				ww, ok := qq.Where.(SqlBuilder)
				if ok {
					ww = append(ww, SqlBuilder{{
						ship.ForeignKey, "in", "",
					}}...)
					qq.Where = ww
				} else {
					aw, ok := helper.IsImplements[AndWhere](qq.Where)
					if ok {
						ww := aw.AndWhere(ship.ForeignKey, "in", strings.Join(slice.Map(ids, func(t any) string {
							return fmt.Sprintf("%v", t)
						}), ","), "int")
						qq.Where = ww
					}
				}
				qq.In = in
			}
			err = ParseRelation(isPlural || ship.RelationType == HasMany, db, ctx, helper.Or(isPlural, rrs, rr), qq)
			if err != nil {
				if err == sql.ErrNoRows {
					err = nil
				} else {
					return err
				}
			}
			assignmentFn(r, helper.Or(isPlural, rrs, rr))
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

// RelationHasOne
// eg: post has a user. fId is post's userId, pId is user's id
func RelationHasOne[M, P any, I constraints.Integer | constraints.Unsigned](fId func(*M) I, pId func(*P) I, setVal func(*M, *P), r Relationship) RelationFn {
	idFn := GetWithID(fId)
	setFn := SetHasOne(setVal, fId, pId)
	return func() (func(any) []any, func(any, any), any, any, Relationship) {
		var s P
		var ss []P
		return idFn, setFn, &s, &ss, r
	}
}

// RelationHasMany
// eg: post has many comments,mId is post's id, pId is comment's postId
func RelationHasMany[M, P any, I constraints.Integer | constraints.Unsigned](mId func(*M) I, pId func(*P) I, setVal func(*M, *[]P), r Relationship) RelationFn {
	idFn := GetWithID(mId)
	setFn := SetHasMany(setVal, mId, pId)
	return func() (func(any) []any, func(any, any), any, any, Relationship) {
		var ss []P
		return idFn, setFn, &ss, &ss, r
	}
}

func AddRelationFn(getVal, join bool, q *QueryCondition, r RelationFn) func() (bool, bool, *QueryCondition, RelationFn) {
	return func() (bool, bool, *QueryCondition, RelationFn) {
		return getVal, join, q, r
	}
}
