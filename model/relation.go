package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/fthvgb1/wp-go/helper"
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

func Relation(db dbQuery, ctx context.Context, r any, q *QueryCondition) ([]func(), []func() error) {
	var fn []func()
	var fns []func() error
	for _, f := range q.RelationFn {
		getVal, isJoin, qq, ff := f()
		idFn, assignment, rr, ship := ff()
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
					qq.In = [][]any{idFn(r)}
					qq.Where = ww
				}
				if qq.From == "" {
					qq.From = ship.Table
				}
			}
			// todo finds的情况
			switch ship.RelationType {
			case "hasOne":
				err = parseRelation(false, db, ctx, rr, qq)
			case "hasMany":
				err = parseRelation(true, db, ctx, rr, qq)
			}
			if err != nil && err != sql.ErrNoRows {
				return err
			}
			err = assignment(r, rr)
			return err
		})
	}
	return fn, fns
}
