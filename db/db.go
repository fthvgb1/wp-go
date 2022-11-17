package db

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github/fthvgb1/wp-go/config"
	"log"
	"os"
	"strconv"
	"strings"
)

var Db *sqlx.DB

type SqlxDb struct {
	sqlx *sqlx.DB
}

func NewSqlxDb(sqlx *sqlx.DB) *SqlxDb {
	return &SqlxDb{sqlx: sqlx}
}

func (r SqlxDb) Select(ctx context.Context, dest any, sql string, params ...any) error {
	if os.Getenv("SHOW_SQL") == "true" {
		go log.Println(formatSql(sql, params))
	}
	return r.sqlx.Select(dest, sql, params...)
}

func (r SqlxDb) Get(ctx context.Context, dest any, sql string, params ...any) error {
	if os.Getenv("SHOW_SQL") == "true" {
		go log.Println(formatSql(sql, params))
	}
	return r.sqlx.Get(dest, sql, params...)
}

func formatSql(sql string, params []any) string {
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

func InitDb() error {
	c := config.Conf.Load()
	dsn := c.Mysql.Dsn.GetDsn()
	var err error
	Db, err = sqlx.Open("mysql", dsn)
	if err != nil {
		return err
	}
	if c.Mysql.Pool.ConnMaxIdleTime != 0 {
		Db.SetConnMaxIdleTime(c.Mysql.Pool.ConnMaxLifetime)
	}
	if c.Mysql.Pool.MaxIdleConn != 0 {
		Db.SetMaxIdleConns(c.Mysql.Pool.MaxIdleConn)
	}
	if c.Mysql.Pool.MaxOpenConn != 0 {
		Db.SetMaxOpenConns(c.Mysql.Pool.MaxOpenConn)
	}
	if c.Mysql.Pool.ConnMaxLifetime != 0 {
		Db.SetConnMaxLifetime(c.Mysql.Pool.ConnMaxLifetime)
	}
	return err
}
