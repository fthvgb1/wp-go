package model

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

func setTable[T Model](q *QueryCondition) {
	if q.From == "" {
		q.From = Table[T]()
	}
}

func Relation[T Model](db dbQuery, ctx context.Context, r *T, q *QueryCondition) (err error) {
	var rr T
	t := reflect.TypeOf(rr)
	v := reflect.ValueOf(r).Elem()
	for tableTag, relation := range q.Relation {
		if tableTag == "" {
			continue
		}
		for i := 0; i < t.NumField(); i++ {
			tag := t.Field(i).Tag
			table, ok := tag.Lookup("table")
			if !ok || table == "" {
				continue
			}
			tables := strings.Split(table, " ")
			if tables[len(tables)-1] != tableTag {
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
					id = fmt.Sprintf("%v", v.Field(j).Interface())
					break
				}
			}
			{
				var w any = relation.Where
				if w == nil {
					w = SqlBuilder{}
				}
				ww, ok := w.(SqlBuilder)
				if ok {
					ww = append(ww, SqlBuilder{{
						tag.Get("foreignKey"), "=", id, "int",
					}}...)
					relation.Where = ww
				}
			}
			sq, args, er := BuildQuerySql(*relation)
			if er != nil {
				err = er
				return
			}
			vv := reflect.New(v.Field(i).Type().Elem()).Interface()
			switch tag.Get("relation") {
			case "hasOne":
				err = db.Get(ctx, vv, sq, args...)
			case "hasMany":
				err = db.Select(ctx, vv, sq, args...)
			}
			if err != nil {
				return
			}
			v.Field(i).Set(reflect.ValueOf(vv))
		}
	}
	return
}
