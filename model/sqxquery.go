package model

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"strconv"
	"strings"
)

type SqlxQuery struct {
	sqlx *sqlx.DB
}

func NewSqlxQuery(sqlx *sqlx.DB) SqlxQuery {
	return SqlxQuery{sqlx: sqlx}
}

func (r SqlxQuery) Select(ctx context.Context, dest any, sql string, params ...any) error {
	if os.Getenv("SHOW_SQL") == "true" {
		go log.Println(FormatSql(sql, params...))
	}
	return r.sqlx.Select(dest, sql, params...)
}

func (r SqlxQuery) Get(ctx context.Context, dest any, sql string, params ...any) error {
	if os.Getenv("SHOW_SQL") == "true" {
		go log.Println(FormatSql(sql, params...))
	}
	return r.sqlx.Get(dest, sql, params...)
}

func FormatSql(sql string, params ...any) string {
	for _, param := range params {
		switch param.(type) {
		case string:
			sql = strings.Replace(sql, "?", fmt.Sprintf("'%s'", param.(string)), 1)
		case int64:
			sql = strings.Replace(sql, "?", strconv.FormatInt(param.(int64), 10), 1)
		case int:
			sql = strings.Replace(sql, "?", strconv.Itoa(param.(int)), 1)
		case uint64:
			sql = strings.Replace(sql, "?", strconv.FormatUint(param.(uint64), 10), 1)
		case float64:
			sql = strings.Replace(sql, "?", fmt.Sprintf("%f", param.(float64)), 1)
		}
	}
	return sql
}
