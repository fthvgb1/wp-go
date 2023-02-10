package db

import (
	"context"
	"github.com/fthvgb1/wp-go/internal/pkg/config"
	"github.com/fthvgb1/wp-go/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

var db *sqlx.DB

func InitDb() (*sqlx.DB, error) {
	c := config.GetConfig()
	dsn := c.Mysql.Dsn.GetDsn()
	var err error
	db, err = sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if c.Mysql.Pool.ConnMaxIdleTime != 0 {
		db.SetConnMaxIdleTime(c.Mysql.Pool.ConnMaxLifetime)
	}
	if c.Mysql.Pool.MaxIdleConn != 0 {
		db.SetMaxIdleConns(c.Mysql.Pool.MaxIdleConn)
	}
	if c.Mysql.Pool.MaxOpenConn != 0 {
		db.SetMaxOpenConns(c.Mysql.Pool.MaxOpenConn)
	}
	if c.Mysql.Pool.ConnMaxLifetime != 0 {
		db.SetConnMaxLifetime(c.Mysql.Pool.ConnMaxLifetime)
	}
	return db, err
}

func QueryDb(db *sqlx.DB) model.UniversalDb {
	query := model.NewUniversalDb(

		func(ctx context.Context, a any, s string, args ...any) error {
			if config.GetConfig().ShowQuerySql {
				go log.Println(model.FormatSql(s, args...))
			}
			return db.Select(a, s, args...)
		},

		func(ctx context.Context, a any, s string, args ...any) error {
			if config.GetConfig().ShowQuerySql {
				go log.Println(model.FormatSql(s, args...))
			}
			return db.Get(a, s, args...)
		})

	return query
}
