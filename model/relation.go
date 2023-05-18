package model

import (
	"context"
	"fmt"
	"github.com/fthvgb1/wp-go/helper"
	"reflect"
	"strings"
)

func setTable[T Model](q *QueryCondition) {
	if q.From == "" {
		q.From = Table[T]()
	}
}

func Relation(db dbQuery, ctx context.Context, r any, q *QueryCondition) ([]func(), []func() error) {
	var fn []func()
	var fns []func() error
	t := reflect.TypeOf(r).Elem()
	v := reflect.ValueOf(r).Elem()
	for tableTag, relation := range q.Relation {
		if tableTag == "" {
			continue
		}
		tableTag := tableTag
		relation := relation
		for i := 0; i < t.NumField(); i++ {
			i := i
			tag := t.Field(i).Tag
			table, ok := tag.Lookup("table")
			if !ok || table == "" {
				continue
			}
			tables := strings.Split(table, " ")
			if tables[len(tables)-1] != tableTag {
				continue
			}
			foreignKey := tag.Get("foreignKey")
			if foreignKey == "" {
				continue
			}
			localKey := tag.Get("local")
			if localKey == "" {
				continue
			}
			if relation == nil {
				relation = &QueryCondition{
					Fields: "*",
				}
			}
			relation.From = table
			id := ""
			j := 0
			for ; j < t.NumField(); j++ {
				vvv, ok := t.Field(j).Tag.Lookup("db")
				if ok && vvv == tag.Get("local") {
					break
				}
			}
			if relation.WithJoin {
				from := strings.Split(q.From, " ")
				fn = append(fn, func() {
					qq := helper.GetContextVal(ctx, "ancestorsQueryCondition", q)
					qq.Join = append(q.Join, SqlBuilder{
						{"left join", table, fmt.Sprintf("%s.%s=%s.%s", tables[len(tables)-1], foreignKey, from[len(from)-1], localKey)},
					}...)
				})
			}
			fns = append(fns, func() error {
				{
					var w any = relation.Where
					if w == nil {
						w = SqlBuilder{}
					}
					ww, ok := w.(SqlBuilder)
					if ok {
						id = fmt.Sprintf("%v", v.Field(j).Interface())
						ww = append(ww, SqlBuilder{{
							foreignKey, "=", id, "int",
						}}...)
						relation.Where = ww
					}
				}
				var err error
				vv := reflect.New(v.Field(i).Type().Elem()).Interface()
				switch tag.Get("relation") {
				case "hasOne":
					err = parseRelation(false, db, ctx, vv, relation)
				case "hasMany":
					err = parseRelation(true, db, ctx, vv, relation)
				}
				if err != nil {
					return err
				}
				v.Field(i).Set(reflect.ValueOf(vv))
				return nil
			})
		}
	}
	return fn, fns
}
