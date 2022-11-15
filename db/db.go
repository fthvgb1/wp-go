package db

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github/fthvgb1/wp-go/config"
	"log"
	"os"
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
		log.Printf(strings.Replace(sql, "?", "'%v'", -1), params...)
	}
	return r.sqlx.Select(dest, sql, params...)
}

func (r SqlxDb) Get(ctx context.Context, dest any, sql string, params ...any) error {
	if os.Getenv("SHOW_SQL") == "true" {
		log.Printf(strings.Replace(sql, "?", "'%v'", -1), params...)
	}
	return r.sqlx.Get(dest, sql, params...)
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
